package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mamenzul/go-rest/configs"
	"github.com/mamenzul/go-rest/services/user"
	"github.com/tursodatabase/go-libsql"
)

type APIServer struct {
	db *sql.DB
}

func main() {
	dbName := "local.db"
	authToken := configs.Envs.AUTH_TOKEN
	url := configs.Envs.DATABASE_URL

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, url,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()
	db := sql.OpenDB(connector)
	defer db.Close()

	// The HTTP Server
	server := &http.Server{Addr: fmt.Sprintf(":%s", configs.Envs.Port), Handler: service(&APIServer{db: db})}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		fmt.Println("Shutting down server..")
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Server gracefully stopped")
		serverStopCtx()
	}()

	// Run the server
	fmt.Printf("Server running on %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func service(s *APIServer) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(r)

	return r
}

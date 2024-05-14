package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost      string
	Port            string
	DATABASE_URL    string
	AUTH_TOKEN      string
	MAILGUN_API_KEY string
	MAILGUN_SENDER  string
	MAILGUN_DOMAIN  string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:      getEnv("PUBLIC_HOST", "http://localhost"),
		Port:            getEnv("PORT", "8080"),
		DATABASE_URL:    getEnv("DATABASE_URL", "panic"),
		AUTH_TOKEN:      getEnv("AUTH_TOKEN", "panic"),
		MAILGUN_API_KEY: getEnv("MAILGUN_API_KEY", "panic"),
		MAILGUN_SENDER:  getEnv("MAILGUN_SENDER", "panic"),
		MAILGUN_DOMAIN:  getEnv("MAILGUN_DOMAIN", "panic"),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "panic" {
		log.Panic("Env not found: " + key)
	}
	return fallback
}

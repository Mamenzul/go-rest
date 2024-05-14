package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mamenzul/go-rest/services/auth"
	mailgun "github.com/mamenzul/go-rest/services/mail"
	"github.com/mamenzul/go-rest/types"
	"github.com/mamenzul/go-rest/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Post("/register", h.handleRegister)
	router.Post("/reset-password", h.handleResetPassword)
	router.Post("/reset-password-token", h.handleResetPasswordToken)
	router.Get("/users", h.handleGetUsers)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		Email:    user.Email,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var user types.ResetPasswordPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		token, err := h.store.StoreResetToken(user.Email)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		body := "Click here to reset your password: http://localhost:3000/reset-password-token?token=" + token
		_, err = mailgun.SendSimpleMessage("Reset password", body, user.Email)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	message := "If the email exists, a password reset link will be sent to it."

	utils.WriteJSON(w, http.StatusCreated, message)
}

func (h *Handler) handleResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("token is required"))
		return
	}
	var user types.ResetPasswordTokenPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	//check if token is valid
	valid, err := h.store.CheckResetToken(token)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !valid {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid token"))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdatePassword(user.Email, hashedPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.DeleteResetToken(token)
	if err != nil {
		log.Default().Println(err)
		return
	}

	message := "Password updated successfully"

	utils.WriteJSON(w, http.StatusCreated, message)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

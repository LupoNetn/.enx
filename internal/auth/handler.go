package auth

import (
	"encoding/json"
	"net/http"

	"github.com/luponetn/.enx/internal/utils"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Name == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "email, name and password are required")
		return
	}

	res, err := h.service.Register(r.Context(), RegisterInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, "account created successfully", res)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	res, err := h.service.Login(r.Context(), LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "login successful", res)
}

func (h *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		utils.WriteError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	tokens, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "tokens refreshed successfully", tokens)
}
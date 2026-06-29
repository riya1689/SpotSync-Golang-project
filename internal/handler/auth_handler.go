package handler 

import (
	"net/http"
	"spotSync-golang-Project/internal/dto"
	"spotSync-golang-Project/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService 
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, validate *validator.Validate) *AuthHandler { 
	return &AuthHandler{authService, validate} 
}

func (h *AuthHandler) Register(c echo.Context) error { 
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil { 
		return JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if err := h.validate.Struct(req); err != nil { 
		return JSONError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.authService.Register(req) 
	if err != nil { 
		if err.Error() == "email already in use" { 
			return JSONError(c, http.StatusBadRequest, "Registration failed", err.Error()) 
		}
		return JSONError(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(c echo.Context) error { 
	var req dto.LoginRequest 
	if err := c.Bind(&req); err != nil { 
		return JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error()) 
	}

	if err := h.validate.Struct(req); err != nil { 
		return JSONError(c, http.StatusBadRequest, "Validation failed", err.Error()) 
	}

	res, err := h.authService.Login(req)
	if err != nil { 
		if err.Error() == "invalid email or password" {
			return JSONError(c, http.StatusUnauthorized, "Login failed", err.Error()) 
		}
		return JSONError(c, http.StatusInternalServerError, "Internal server error", err.Error()) 
	}

	return JSONSuccess(c, http.StatusOK, "Login successful", res)
}

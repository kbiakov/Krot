package main

import (
	"github.com/labstack/echo"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const jwtSecret = "secret"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
	User        User   `json:"user"`
}

func CheckUserAuth(h echo.HandlerFunc) echo.HandlerFunc {
	isAuthenticated := func(ctx echo.Context) bool {
		u := ctx.Get("user").(*jwt.Token)
		claims := u.Claims.(jwt.MapClaims)
		authID := claims["name"].(string)
		return authID == ctx.Param("uid")
	}

	return func(ctx echo.Context) error {
		if !isAuthenticated(ctx) {
			return echo.ErrUnauthorized
		}
		return h(ctx)
	}
}

func NewAuthResponse(user *User) (*AuthResponse, error) {
	// Create token
	t := jwt.New(jwt.SigningMethodHS256)
	expiresAt := time.Now().Add(time.Hour * 72).Unix()

	// Set claims
	claims := t.Claims.(jwt.MapClaims)
	claims["name"] = user.ID
	claims["admin"] = false
	claims["exp"] = expiresAt

	// Generate encoded token
	accessToken, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	// Form response
	return &AuthResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		User:        *user,
	}, nil
}

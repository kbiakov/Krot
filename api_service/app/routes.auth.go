package main

import (
	"github.com/labstack/echo"
	"net/http"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type AuthRequest struct {
	Email	 string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
	User	    User   `json:"user"`
}

func isAuthenticated(ctx echo.Context) bool {
	u := ctx.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	authID := claims["name"].(string)
	return authID == ctx.Param("uid")
}

func newAuthResponse(user User) (*AuthResponse, error) {
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
		ExpiresAt: expiresAt,
		User: user,
	}, nil
}

func SignUp(ctx echo.Context) error {
	u := new(User)
	if err := ctx.Bind(u); err != nil {
		return err
	}

	existingUser, err := GetUser(u.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("User with this email is already exists")
	}

	if err := u.CreateUser(); err != nil {
		return err
	}

	res, err := newAuthResponse(u)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func Login(ctx echo.Context) error {
	req := new(AuthRequest)
	if err := ctx.Bind(req); err != nil {
		return err
	}

	u, err := GetUser(req.Email)
	if err != nil {
		return err
	}
	if !u.IsValidPassword(req.Password) {
		return echo.ErrUnauthorized
	}

	res, err := newAuthResponse(u)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func Logout(ctx echo.Context) error {
	// TODO: stop all subscriptions
	return ctx.JSON(http.StatusNotImplemented, nil)
}

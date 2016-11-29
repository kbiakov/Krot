package krot

import (
	jwt "github.com/dgrijalva/jwt-go"
	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	"sync"
)

type JwtMiddleware struct {
	mw *jwtMiddleware.Middleware
}

var instance *JwtMiddleware

func GetJwtMiddlewareInstance() *JwtMiddleware {
	sync.Once.Do(func() {
		instance = &JwtMiddleware{mw: newMiddleware()}
	})
	return instance
}

func newMiddleware() jwtMiddleware.Middleware {
	return jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// TODO: put JWT secret here
			return []byte("..."), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

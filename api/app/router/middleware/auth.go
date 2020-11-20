package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

var (
	// CtxKeyJWTClaims is a context key for saving data from JWT's payload
	CtxKeyJWTClaims = contextKey("ContextKeyJWTClaims")
)

// TokenValidation is a middleware that checks the incoming requests for token validity
func TokenValidation(AtJwtSecretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

			if len(authHeader) != 2 {
				// w.WriteHeader(http.StatusUnauthorized)
				// w.Write([]byte("Malformed Token"))
				next.ServeHTTP(w, r)
			} else {
				jwtToken := authHeader[1]
				token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return AtJwtSecretKey, nil
				})

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					ctx := context.WithValue(r.Context(), CtxKeyJWTClaims, claims)
					// Access context values in handlers like this
					// props, _ := r.Context().Value("props").(jwt.MapClaims)
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					fmt.Println(err)
					// w.WriteHeader(http.StatusUnauthorized)
					// w.Write([]byte("Unauthorized"))
					next.ServeHTTP(w, r)
				}
			}
		})
	}
}

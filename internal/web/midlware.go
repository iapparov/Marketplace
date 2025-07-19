package web

import (
	"context"
	"marketplace/internal/app"
	"net/http"
	"strings"
)

type ContextKey string
const UserIDKey ContextKey = "user_id"

func AuthMiddleware(jwtProvider *app.JwtProvider) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing auth header", http.StatusUnauthorized)
                return
            }
            parts := strings.Split(authHeader, "Bearer ")
            if len(parts) != 2 {
                http.Error(w, "invalid auth header", http.StatusUnauthorized)
                return
            }
            tokenStr := parts[1]

            claims, err := jwtProvider.ValidateAccessToken(tokenStr)
            if err != nil {
                http.Error(w, "invalid or expired token", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), UserIDKey, claims["uuid"])
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func OptionalAuthMiddleware(jwtProvider *app.JwtProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}
			parts := strings.Split(authHeader, "Bearer ")
			if len(parts) != 2 {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := parts[1]

			claims, err := jwtProvider.ValidateAccessToken(tokenStr)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims["uuid"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
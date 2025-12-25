package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"social-network/pkg/repository"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow OPTIONS (CORS)
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		ses, err := r.Cookie("session")
		if err != nil {
			http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Validate session against DB (Shared DB pattern for migration)
		userID, err := repository.GetUserIDBySessionToken(ses.Value)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			// Log internal error but return 401/500
			log.Printf("Database error in middleware: %v\n", err)
			http.Error(w, `{"message":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Store user ID in context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS is handled by API Gateway. Services should not add these headers to avoid duplicates.
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		// Preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

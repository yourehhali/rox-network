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

		
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		ses, err := r.Cookie("session")
		if err != nil {
			http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		
		userID, err := repository.GetUserIDBySessionToken(ses.Value)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			
			log.Printf("Database error in middleware: %v\n", err)
			http.Error(w, `{"message":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		
		
		
		

		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

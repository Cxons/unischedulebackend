package authmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Cxons/unischedulebackend/internal/shared/constants"
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
)


func JwtMiddleware()func(http.Handler)http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// tokenString := strings.TrimPrefix(authHeader,"Bearer ")
			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			claims,err := jwt.VerifyToken(tokenString)
			 if err != nil{
				fmt.Print(err)
				http.Error(w,"Invalid or expired token",http.StatusUnauthorized)
				return
			 }
			//  ctxKey := &userContext{Value: "UserInfo"}
			 ctx := context.WithValue(r.Context(),constants.UserInfoKey,claims)
			 next.ServeHTTP(w,r.WithContext(ctx))
		})
	}
} 

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow your frontend origin
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        // Allow credentials like cookies
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        // Allowed headers
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        // Allowed methods
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

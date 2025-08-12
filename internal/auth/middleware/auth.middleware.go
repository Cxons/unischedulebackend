package authmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Cxons/unischedulebackend/internal/auth/dto"
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
)
type userContext dto.UserContext;


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
			 ctxKey := &userContext{Value: "UserInfo"}
			 ctx := context.WithValue(r.Context(),ctxKey,claims)
			 next.ServeHTTP(w,r.WithContext(ctx))
		})
	}
} 
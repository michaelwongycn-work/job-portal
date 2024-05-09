package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/michaelwongycn/job-portal/lib/auth"
	"github.com/michaelwongycn/job-portal/lib/cache"
)

func Authorize(allowedRoles []int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 {
				http.Error(w, "Malformed token", http.StatusUnauthorized)
				return
			}
			accessToken := tokenParts[1]

			cachedAccessToken := cache.GetCache(accessToken)

			if cachedAccessToken == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ParseToken(accessToken)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			role := int(claims["rle"].(float64))
			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Package middleware provides HTTP middleware functions for the API Gateway.
// This package includes authentication handling to verify JWT tokens and attach user-specific data to requests.

package middleware

import (
	"api-gateway/internal/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is an HTTP middleware function that verifies the JWT token from incoming requests.
// If the token is valid, it extracts user information and forwards the request to the next handler.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.CFG.PublicRoutes != nil {
			if _, ok := config.CFG.PublicRoutes[r.URL.Path]; ok {
				next(w, r)
				return
			}
		}

		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
			return
		}

		r.Header.Add("User-ID", userID)

		next(w, r)
	}
}

// extractToken retrieves the Bearer token from the Authorization header of the request.
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}

// validateToken parses and validates the JWT token string.
// It checks for the correct signing method and returns the claims if the token is valid.
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(config.CFG.SecretKey)

		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return nil, errors.New("token expired")
		}
	}

	if iat, ok := claims["iat"].(float64); ok {
		issuedAt := time.Unix(int64(iat), 0)
		if time.Now().Before(issuedAt) {
			return nil, errors.New("token issued in the future")
		}
	}

	if nbf, ok := claims["nbf"].(float64); ok {
		notBefore := time.Unix(int64(nbf), 0)
		if time.Now().Before(notBefore) {
			return nil, errors.New("token not valid yet")
		}
	}

	return claims, nil
}

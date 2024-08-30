package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware(jwtKey []byte) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		SigningKey:    jwtKey,
		SigningMethod: "HS256",
		ContextKey:    "user",
		TokenLookup:   "header:Authorization",
		AuthScheme:    "Bearer",
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			var status int
			var message string

			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					status = http.StatusUnauthorized
					message = "Token has expired"
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					status = http.StatusUnauthorized
					message = "Invalid token signature"
				}
			}

			return c.JSON(status, map[string]string{"error": message})
		},
	}
	return middleware.JWTWithConfig(config)
}

// RoleMiddleware function to check user role
func RoleMiddleware(requiredRole string, jwtKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token not found"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			role, ok := claims["role"].(string)
			if !ok || role != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have permission to access this route"})
			}

			return next(c)
		}
	}
}

package middlewares

import (
	"app/global"
	"app/pkg/jwt"
	"app/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.DataDetailResponse(c, 401, response.ErrCodeUnauthorized, nil)
			c.Abort()
			return
		}

		// Check if the header has the "Bearer" prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.DataDetailResponse(c, 401, response.ErrInvalidToken, nil)
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ValidateToken(tokenString, global.Config.JWT.SecretKey)
		if err != nil {
			response.DataDetailResponse(c, 401, response.ErrInvalidToken, nil)
			c.Abort()
			return
		}

		// Store user info in context for later use
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("system_role", claims.SystemRole)

		c.Next()
	}
}

// RoleMiddleware checks if the user has a required role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("system_role")
		if !exists {
			response.DataDetailResponse(c, 401, response.ErrInvalidToken, nil)
			c.Abort()
			return
		}

		// Check if a user role is in the allowed roles
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.DataDetailResponse(c, 403, response.ErrCodeAccessDenied, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

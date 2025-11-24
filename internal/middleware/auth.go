package middleware

import (
	"fmt"
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your_secret_key") // Should match the one in handlers

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Parse user_id from claims
			// We expect user_id to be a string (UUID)
			if idStr, ok := claims["user_id"].(string); ok {
				if uid, err := uuid.Parse(idStr); err == nil {
					c.Set("userID", uid)

					// Load User with Role and Permissions
					var user models.User
					if err := database.DB.Preload("Role.Permissions").First(&user, "id = ?", uid).Error; err == nil {
						c.Set("user", user)
						c.Set("role", user.Role.Name)

						// Create a map of permissions for easy lookup
						perms := make(map[string]bool)
						for _, p := range user.Role.Permissions {
							perms[p.Resource+":"+p.Action] = true
						}
						c.Set("permissions", perms)
					}
				}
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}

func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			c.Abort()
			return
		}

		permissions := perms.(map[string]bool)
		if !permissions[resource+":"+action] && !permissions["*:*"] { // Check specific or superadmin
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

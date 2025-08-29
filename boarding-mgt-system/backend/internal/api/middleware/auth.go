package middleware

import (
	"net/http"
	"strings"

	"github.com/ferryflow/boarding-mgt-system/internal/auth"
	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate token
		jwtCfg := &auth.JWTConfig{
			Secret: []byte(jwtConfig.Secret),
			Issuer: "ferryflow",
		}

		claims, err := auth.ValidateToken(jwtCfg, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check token type
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Set("operator_id", claims.OperatorID)
		c.Set("session_id", claims.SessionID)

		c.Next()
	}
}

// RequireRole checks if user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User type not found"})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user type"})
			c.Abort()
			return
		}

		// Check if user has required role
		for _, role := range roles {
			if userTypeStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// RequireOperator ensures user belongs to specific operator
func RequireOperator() gin.HandlerFunc {
	return func(c *gin.Context) {
		operatorID, exists := c.Get("operator_id")
		if !exists || operatorID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Operator access required"})
			c.Abort()
			return
		}

		// Get operator ID from path parameter if present
		pathOperatorID := c.Param("operator_id")
		if pathOperatorID != "" {
			// Check if user's operator matches the requested operator
			if operatorID != pathOperatorID {
				userType, _ := c.Get("user_type")
				// System admins can access any operator
				if userType != "system_admin" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this operator"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// OptionalAuth allows both authenticated and unauthenticated requests
func OptionalAuth(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Try to validate token but don't fail if invalid
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err == nil {
			jwtCfg := &auth.JWTConfig{
				Secret: []byte(jwtConfig.Secret),
				Issuer: "ferryflow",
			}

			claims, err := auth.ValidateToken(jwtCfg, token)
			if err == nil && claims.TokenType == "access" {
				// Store user info if valid
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_type", claims.UserType)
				c.Set("operator_id", claims.OperatorID)
				c.Set("session_id", claims.SessionID)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}

// GetUserID gets the authenticated user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	
	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetOperatorID gets the user's operator ID from context
func GetOperatorID(c *gin.Context) (string, bool) {
	operatorID, exists := c.Get("operator_id")
	if !exists {
		return "", false
	}
	
	operatorIDStr, ok := operatorID.(string)
	return operatorIDStr, ok
}

// IsAuthenticated checks if the request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	auth, exists := c.Get("authenticated")
	if !exists {
		return false
	}
	
	isAuth, ok := auth.(bool)
	return ok && isAuth
}
package middleware

import (
	"net/http"
	"strings"
	"time"

	"japanese-learning-app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTSecret is the secret key for JWT
var JWTSecret []byte

// SetJWTSecret sets the JWT secret from config
func SetJWTSecret(secret string) {
	JWTSecret = []byte(secret)
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// AuthRequired is a middleware that requires authentication
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Get database from context (we'll set this up in main.go)
		db, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database connection not available",
			})
			c.Abort()
			return
		}

		// Get user from database to ensure they still exist and are active
		var user models.User
		if err := db.(*gorm.DB).Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found or inactive",
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// GetCurrentUser returns the current user from the Gin context
func GetCurrentUser(c *gin.Context) (*models.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, jwt.ErrSignatureInvalid
	}

	currentUser, ok := user.(*models.User)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return currentUser, nil
}

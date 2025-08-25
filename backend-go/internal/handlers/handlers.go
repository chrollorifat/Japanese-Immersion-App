package handlers

import (
	"net/http"
	"strconv"
	"time"

	"japanese-learning-app/internal/middleware"
	"japanese-learning-app/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler holds the database connection and other dependencies
type Handler struct {
	db *gorm.DB
}

// New creates a new handler with the given database connection
func New(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username or email already exists",
		})
		return
	}

	// Create new user
	user := models.User{
		Username:          req.Username,
		Email:             req.Email,
		Password:          req.Password,
		PreferredLanguage: req.PreferredLanguage,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Save user to database
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user.ToResponse(),
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find user by username
	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Account is inactive",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Update last activity
	now := time.Now()
	user.LastActivity = &now
	h.db.Save(&user)

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		User:        user.ToResponse(),
	})
}

// GetCurrentUser returns the current authenticated user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// Logout handles user logout (client-side token removal)
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetUserBooks returns all books for the current user
func (h *Handler) GetUserBooks(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	var books []models.Book
	if err := h.db.Where("user_id = ?", user.ID).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch books",
		})
		return
	}

	// Convert to response format
	var bookResponses []models.BookResponse
	for _, book := range books {
		bookResponses = append(bookResponses, book.ToResponse())
	}

	c.JSON(http.StatusOK, bookResponses)
}

// UploadBook handles book file upload
func (h *Handler) UploadBook(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded",
		})
		return
	}

	// TODO: Implement file validation and processing
	// For now, just return a placeholder response

	c.JSON(http.StatusOK, gin.H{
		"message": "Book upload functionality not yet implemented",
		"file":    file.Filename,
		"user_id": user.ID,
	})
}

// GetBook returns a specific book
func (h *Handler) GetBook(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID",
		})
		return
	}

	var book models.Book
	if err := h.db.Where("id = ? AND user_id = ?", bookID, user.ID).First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Book not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch book",
		})
		return
	}

	c.JSON(http.StatusOK, book.ToResponse())
}

// DeleteBook deletes a book
func (h *Handler) DeleteBook(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID",
		})
		return
	}

	// TODO: Also delete the actual file
	if err := h.db.Where("id = ? AND user_id = ?", bookID, user.ID).Delete(&models.Book{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete book",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}

// Placeholder handlers for features to be implemented
func (h *Handler) LookupWord(c *gin.Context) {
	word := c.Param("word")
	c.JSON(http.StatusOK, gin.H{
		"word":    word,
		"message": "Word lookup not yet implemented",
	})
}

func (h *Handler) MarkWordAsKnown(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Mark word as known not yet implemented",
	})
}

func (h *Handler) GetDueCards(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"due_cards": []interface{}{},
		"message":   "SRS system not yet implemented",
	})
}

func (h *Handler) ReviewCard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "SRS review not yet implemented",
	})
}

func (h *Handler) StartReadingSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Reading session tracking not yet implemented",
	})
}

func (h *Handler) EndReadingSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Reading session tracking not yet implemented",
	})
}

func (h *Handler) GetReadingStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_reading_time": 0,
		"books_read":         0,
		"words_learned":      0,
		"message":            "Reading statistics not yet implemented",
	})
}

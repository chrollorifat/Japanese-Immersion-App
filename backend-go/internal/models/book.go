package models

import (
	"time"

	"gorm.io/gorm"
)

// Book represents an uploaded ebook
type Book struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign key
	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Basic metadata
	Title    string `json:"title" gorm:"size:500;not null"`
	Author   string `json:"author" gorm:"size:200"`
	Language string `json:"language" gorm:"size:10;default:ja"`

	// File information
	FilePath string `json:"file_path" gorm:"size:500;not null"`
	FileName string `json:"file_name" gorm:"size:255;not null"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type" gorm:"size:100"`

	// Processing status
	ProcessingStatus string `json:"processing_status" gorm:"size:20;default:pending"` // pending, processing, completed, failed
	WordCount        int    `json:"word_count" gorm:"default:0"`
	UniqueWordCount  int    `json:"unique_word_count" gorm:"default:0"`
	DifficultyLevel  string `json:"difficulty_level" gorm:"size:10"` // beginner, intermediate, advanced

	// Reading progress
	LastReadAt      *time.Time `json:"last_read_at"`
	ReadingProgress float64    `json:"reading_progress" gorm:"type:decimal(5,2);default:0.00"` // percentage

	// Content extraction
	ExtractedText string                 `json:"extracted_text" gorm:"type:text"`
	ChapterData   []map[string]interface{} `json:"chapter_data" gorm:"serializer:json"` // [{title, start_pos, end_pos}]

	// Timestamps
	UploadedAt time.Time `json:"uploaded_at" gorm:"autoCreateTime"`

	// Relationships
	ReadingSessions []ReadingSession `json:"reading_sessions,omitempty" gorm:"foreignKey:BookID"`
	Annotations     []BookAnnotation `json:"annotations,omitempty" gorm:"foreignKey:BookID"`
}

// TableName specifies the table name for GORM
func (Book) TableName() string {
	return "books"
}

// ToResponse converts the book to a response format
func (b *Book) ToResponse() BookResponse {
	return BookResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Title:            b.Title,
		Author:           b.Author,
		Language:         b.Language,
		FileName:         b.FileName,
		FileSize:         b.FileSize,
		MimeType:         b.MimeType,
		ProcessingStatus: b.ProcessingStatus,
		WordCount:        b.WordCount,
		UniqueWordCount:  b.UniqueWordCount,
		DifficultyLevel:  b.DifficultyLevel,
		UploadedAt:       b.UploadedAt,
		LastReadAt:       b.LastReadAt,
		ReadingProgress:  b.ReadingProgress,
		ChapterData:      b.ChapterData,
	}
}

// BookResponse is the response format for book data
type BookResponse struct {
	ID               uint                     `json:"id"`
	UserID           uint                     `json:"user_id"`
	Title            string                   `json:"title"`
	Author           string                   `json:"author"`
	Language         string                   `json:"language"`
	FileName         string                   `json:"file_name"`
	FileSize         int64                    `json:"file_size"`
	MimeType         string                   `json:"mime_type"`
	ProcessingStatus string                   `json:"processing_status"`
	WordCount        int                      `json:"word_count"`
	UniqueWordCount  int                      `json:"unique_word_count"`
	DifficultyLevel  string                   `json:"difficulty_level"`
	UploadedAt       time.Time                `json:"uploaded_at"`
	LastReadAt       *time.Time               `json:"last_read_at"`
	ReadingProgress  float64                  `json:"reading_progress"`
	ChapterData      []map[string]interface{} `json:"chapter_data"`
}

// ReadingSession represents a reading session
type ReadingSession struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign keys
	UserID uint `json:"user_id" gorm:"not null;index"`
	BookID uint `json:"book_id" gorm:"not null;index"`
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Book   Book `json:"book,omitempty" gorm:"foreignKey:BookID"`

	// Session timing
	StartTime       time.Time  `json:"start_time" gorm:"autoCreateTime"`
	EndTime         *time.Time `json:"end_time"`
	DurationMinutes int        `json:"duration_minutes"`

	// Progress tracking
	StartPosition int `json:"start_position" gorm:"default:0"` // Character position in text
	EndPosition   int `json:"end_position" gorm:"default:0"`
	WordsLearned  int `json:"words_learned" gorm:"default:0"`
	WordsReviewed int `json:"words_reviewed" gorm:"default:0"`

	// Session data
	NewWordsEncountered []uint `json:"new_words_encountered" gorm:"serializer:json"` // Array of word IDs
	WordsLookedUp       []uint `json:"words_looked_up" gorm:"serializer:json"`       // Array of word IDs
}

// TableName specifies the table name for GORM
func (ReadingSession) TableName() string {
	return "reading_sessions"
}

// BookAnnotation represents an annotation on a book
type BookAnnotation struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign keys
	UserID uint `json:"user_id" gorm:"not null;index"`
	BookID uint `json:"book_id" gorm:"not null;index"`
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Book   Book `json:"book,omitempty" gorm:"foreignKey:BookID"`

	// Position in text
	StartPosition int    `json:"start_position" gorm:"not null"`
	EndPosition   int    `json:"end_position" gorm:"not null"`
	SelectedText  string `json:"selected_text" gorm:"type:text;not null"`

	// Annotation data
	AnnotationType string                 `json:"annotation_type" gorm:"size:20;default:word_lookup"` // word_lookup, note, highlight
	AnnotationData map[string]interface{} `json:"annotation_data" gorm:"serializer:json"`             // Flexible data storage

	// Visual styling
	HighlightColor string `json:"highlight_color" gorm:"size:7"` // Hex color code
}

// TableName specifies the table name for GORM
func (BookAnnotation) TableName() string {
	return "book_annotations"
}

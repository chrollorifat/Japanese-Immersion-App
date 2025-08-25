from sqlalchemy import Column, Integer, String, Boolean, DateTime, Text, JSON, BigInteger, ForeignKey, DECIMAL
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from .database import Base

class Book(Base):
    __tablename__ = "books"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    
    # Basic metadata
    title = Column(String(500), nullable=False)
    author = Column(String(200))
    language = Column(String(10), default="ja")
    
    # File information
    file_path = Column(String(500), nullable=False)
    file_name = Column(String(255), nullable=False)
    file_size = Column(BigInteger)
    mime_type = Column(String(100))
    
    # Processing status
    processing_status = Column(String(20), default="pending")  # pending, processing, completed, failed
    word_count = Column(Integer, default=0)
    unique_word_count = Column(Integer, default=0)
    difficulty_level = Column(String(10))  # beginner, intermediate, advanced
    
    # Timestamps
    uploaded_at = Column(DateTime(timezone=True), server_default=func.now())
    last_read_at = Column(DateTime(timezone=True))
    
    # Reading progress
    reading_progress = Column(DECIMAL(5, 2), default=0.00)  # percentage
    
    # Content extraction
    extracted_text = Column(Text)
    chapter_data = Column(JSON, default=list)  # [{title, start_pos, end_pos}]
    
    # Relationships
    user = relationship("User", back_populates="books")
    reading_sessions = relationship("ReadingSession", back_populates="book")
    annotations = relationship("BookAnnotation", back_populates="book")
    
    def __repr__(self):
        return f"<Book(title='{self.title}', author='{self.author}', user_id={self.user_id})>"
    
    def to_dict(self):
        """Convert book object to dictionary for API responses"""
        return {
            "id": self.id,
            "user_id": self.user_id,
            "title": self.title,
            "author": self.author,
            "language": self.language,
            "file_name": self.file_name,
            "file_size": self.file_size,
            "mime_type": self.mime_type,
            "processing_status": self.processing_status,
            "word_count": self.word_count,
            "unique_word_count": self.unique_word_count,
            "difficulty_level": self.difficulty_level,
            "uploaded_at": self.uploaded_at.isoformat() if self.uploaded_at else None,
            "last_read_at": self.last_read_at.isoformat() if self.last_read_at else None,
            "reading_progress": float(self.reading_progress) if self.reading_progress else 0.0,
            "chapter_data": self.chapter_data or []
        }

# Add the relationship to User model
from .user import User
User.books = relationship("Book", back_populates="user")

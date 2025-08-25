from sqlalchemy import Column, Integer, String, DateTime, Text, JSON, ForeignKey
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from .database import Base

class BookAnnotation(Base):
    __tablename__ = "book_annotations"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    book_id = Column(Integer, ForeignKey("books.id", ondelete="CASCADE"), nullable=False)
    
    # Position in text
    start_position = Column(Integer, nullable=False)
    end_position = Column(Integer, nullable=False)
    selected_text = Column(Text, nullable=False)
    
    # Annotation data
    annotation_type = Column(String(20), default="word_lookup")  # word_lookup, note, highlight
    annotation_data = Column(JSON, default=dict)  # Flexible data storage
    
    # Visual styling
    highlight_color = Column(String(7))  # Hex color code
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    
    # Relationships
    user = relationship("User")
    book = relationship("Book")
    
    def __repr__(self):
        return f"<BookAnnotation(book_id={self.book_id}, annotation_type='{self.annotation_type}')>"
    
    def to_dict(self):
        """Convert annotation object to dictionary for API responses"""
        return {
            "id": self.id,
            "user_id": self.user_id,
            "book_id": self.book_id,
            "start_position": self.start_position,
            "end_position": self.end_position,
            "selected_text": self.selected_text,
            "annotation_type": self.annotation_type,
            "annotation_data": self.annotation_data or {},
            "highlight_color": self.highlight_color,
            "created_at": self.created_at.isoformat() if self.created_at else None,
            "updated_at": self.updated_at.isoformat() if self.updated_at else None
        }

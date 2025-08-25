from sqlalchemy import Column, Integer, String, Boolean, DateTime, Text, JSON
from sqlalchemy.sql import func
from .database import Base

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(50), unique=True, nullable=False, index=True)
    email = Column(String(100), unique=True, nullable=False, index=True)
    hashed_password = Column(String(255), nullable=False)
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    last_activity = Column(DateTime(timezone=True))
    
    # Account status
    is_active = Column(Boolean, default=True)
    preferred_language = Column(String(10), default="en")
    learning_preferences = Column(JSON, default=dict)
    
    # Learning statistics
    total_words_learned = Column(Integer, default=0)
    total_reading_time = Column(Integer, default=0)  # in minutes
    streak_days = Column(Integer, default=0)
    
    def __repr__(self):
        return f"<User(username='{self.username}', email='{self.email}')>"
    
    def to_dict(self):
        """Convert user object to dictionary for API responses"""
        return {
            "id": self.id,
            "username": self.username,
            "email": self.email,
            "created_at": self.created_at.isoformat() if self.created_at else None,
            "is_active": self.is_active,
            "preferred_language": self.preferred_language,
            "learning_preferences": self.learning_preferences,
            "total_words_learned": self.total_words_learned,
            "total_reading_time": self.total_reading_time,
            "streak_days": self.streak_days,
            "last_activity": self.last_activity.isoformat() if self.last_activity else None
        }

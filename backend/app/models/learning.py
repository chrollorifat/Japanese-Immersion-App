from sqlalchemy import Column, Integer, String, Boolean, DateTime, Text, JSON, ForeignKey, DECIMAL
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from .database import Base

class UserWordKnowledge(Base):
    __tablename__ = "user_word_knowledge"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    word_id = Column(Integer, ForeignKey("words.id", ondelete="CASCADE"), nullable=False)
    
    # Knowledge status
    knowledge_level = Column(Integer, default=0)  # 0=unknown, 1=recognized, 2=familiar, 3=known, 4=mastered
    first_encountered_at = Column(DateTime(timezone=True), server_default=func.now())
    last_reviewed_at = Column(DateTime(timezone=True))
    
    # SRS data
    srs_level = Column(Integer, default=0)  # SRS interval level
    next_review_at = Column(DateTime(timezone=True))
    review_count = Column(Integer, default=0)
    correct_count = Column(Integer, default=0)
    streak = Column(Integer, default=0)
    
    # Context tracking
    first_seen_book_id = Column(Integer, ForeignKey("books.id"))
    times_encountered = Column(Integer, default=1)
    
    # Learning metadata
    notes = Column(Text)  # User's personal notes
    is_ignored = Column(Boolean, default=False)  # User marked as "ignore this word"
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

class ReadingSession(Base):
    __tablename__ = "reading_sessions"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    book_id = Column(Integer, ForeignKey("books.id", ondelete="CASCADE"), nullable=False)
    
    start_time = Column(DateTime(timezone=True), server_default=func.now())
    end_time = Column(DateTime(timezone=True))
    duration_minutes = Column(Integer)  # calculated from start/end time
    
    # Progress tracking
    start_position = Column(Integer, default=0)  # Character position in text
    end_position = Column(Integer, default=0)
    words_learned = Column(Integer, default=0)
    words_reviewed = Column(Integer, default=0)
    
    # Session data
    new_words_encountered = Column(JSON, default=list)  # Array of word IDs
    words_looked_up = Column(JSON, default=list)  # Array of word IDs
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    # Relationships
    user = relationship("User")
    book = relationship("Book")

class SRSCard(Base):
    __tablename__ = "srs_cards"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    word_id = Column(Integer, ForeignKey("words.id", ondelete="CASCADE"), nullable=False)
    
    # Card type and content
    card_type = Column(String(20), default="recognition")  # recognition, recall, production
    front_content = Column(JSON, nullable=False)  # {text, furigana, audio_url}
    back_content = Column(JSON, nullable=False)  # {definition, example, notes}
    
    # SRS algorithm data
    ease_factor = Column(DECIMAL(4, 2), default=2.50)  # Anki-style ease factor
    interval_days = Column(Integer, default=1)
    repetition_count = Column(Integer, default=0)
    
    # Scheduling
    due_date = Column(DateTime(timezone=True), server_default=func.now())
    last_reviewed_at = Column(DateTime(timezone=True))
    
    # Performance tracking
    total_reviews = Column(Integer, default=0)
    correct_reviews = Column(Integer, default=0)
    current_streak = Column(Integer, default=0)
    longest_streak = Column(Integer, default=0)
    
    # Card state
    is_suspended = Column(Boolean, default=False)
    is_buried = Column(Boolean, default=False)  # Temporarily hidden
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

class ReviewHistory(Base):
    __tablename__ = "review_history"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    card_id = Column(Integer, ForeignKey("srs_cards.id", ondelete="CASCADE"), nullable=False)
    
    reviewed_at = Column(DateTime(timezone=True), server_default=func.now())
    response_quality = Column(Integer, nullable=False)  # 1=again, 2=hard, 3=good, 4=easy
    response_time_ms = Column(Integer)  # Time taken to answer
    
    # Pre-review state
    old_interval = Column(Integer)
    old_ease_factor = Column(DECIMAL(4, 2))
    
    # Post-review state
    new_interval = Column(Integer)
    new_ease_factor = Column(DECIMAL(4, 2))
    new_due_date = Column(DateTime(timezone=True))
    
    # Context
    review_context = Column(String(50))  # 'reader', 'srs_session', 'manual'
    device_type = Column(String(20))  # 'desktop', 'mobile', 'tablet'
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())

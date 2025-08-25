from sqlalchemy import Column, Integer, String, DateTime, Text, JSON, ForeignKey
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from .database import Base

class Word(Base):
    __tablename__ = "words"

    id = Column(Integer, primary_key=True, index=True)
    surface_form = Column(String(100), nullable=False, index=True)  # Original form (e.g., 食べる)
    reading = Column(String(100), index=True)  # Hiragana reading (e.g., たべる)
    pronunciation = Column(String(100))  # For katakana words
    
    # Linguistic information
    part_of_speech = Column(String(50))
    inflection_type = Column(String(50))
    base_form = Column(String(100))  # Dictionary form
    
    # Difficulty indicators
    jlpt_level = Column(Integer)  # 1-5, null if not in JLPT
    wanikani_level = Column(Integer)  # 1-60, null if not in WaniKani
    kanken_level = Column(Integer)  # 1-10, null if not in Kanken
    frequency_rank = Column(Integer)  # Word frequency ranking
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    
    # Relationships
    definitions = relationship("WordDefinition", back_populates="word")
    
    def __repr__(self):
        return f"<Word(surface_form='{self.surface_form}', reading='{self.reading}')>"
    
    def to_dict(self):
        """Convert word object to dictionary for API responses"""
        return {
            "id": self.id,
            "surface_form": self.surface_form,
            "reading": self.reading,
            "pronunciation": self.pronunciation,
            "part_of_speech": self.part_of_speech,
            "inflection_type": self.inflection_type,
            "base_form": self.base_form,
            "jlpt_level": self.jlpt_level,
            "wanikani_level": self.wanikani_level,
            "kanken_level": self.kanken_level,
            "frequency_rank": self.frequency_rank,
            "created_at": self.created_at.isoformat() if self.created_at else None,
            "updated_at": self.updated_at.isoformat() if self.updated_at else None
        }

class WordDefinition(Base):
    __tablename__ = "word_definitions"

    id = Column(Integer, primary_key=True, index=True)
    word_id = Column(Integer, ForeignKey("words.id", ondelete="CASCADE"), nullable=False)
    dictionary_source = Column(String(50), nullable=False, index=True)  # sanseido, daijirin, ookoku, jitendex, jisho
    language = Column(String(10), nullable=False)  # en, ja
    definition = Column(Text, nullable=False)
    example_sentence = Column(Text)
    example_translation = Column(Text)
    
    # Additional metadata
    definition_order = Column(Integer, default=1)  # For multiple definitions
    tags = Column(JSON, default=list)  # Additional tags or categories
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    # Relationships
    word = relationship("Word", back_populates="definitions")
    
    def __repr__(self):
        return f"<WordDefinition(word_id={self.word_id}, dictionary_source='{self.dictionary_source}')>"
    
    def to_dict(self):
        """Convert word definition object to dictionary for API responses"""
        return {
            "id": self.id,
            "word_id": self.word_id,
            "dictionary_source": self.dictionary_source,
            "language": self.language,
            "definition": self.definition,
            "example_sentence": self.example_sentence,
            "example_translation": self.example_translation,
            "definition_order": self.definition_order,
            "tags": self.tags or [],
            "created_at": self.created_at.isoformat() if self.created_at else None
        }

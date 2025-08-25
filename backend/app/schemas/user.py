from pydantic import BaseModel, EmailStr, validator
from typing import Optional, Dict, Any
from datetime import datetime

class UserBase(BaseModel):
    """Base user schema with common fields"""
    username: str
    email: EmailStr
    preferred_language: Optional[str] = "en"

class UserCreate(UserBase):
    """Schema for user registration"""
    password: str
    
    @validator('username')
    def validate_username(cls, v):
        if len(v) < 3:
            raise ValueError('Username must be at least 3 characters long')
        if len(v) > 50:
            raise ValueError('Username must be less than 50 characters long')
        return v
    
    @validator('password')
    def validate_password(cls, v):
        if len(v) < 6:
            raise ValueError('Password must be at least 6 characters long')
        return v

class UserLogin(BaseModel):
    """Schema for user login"""
    username: str
    password: str

class UserResponse(UserBase):
    """Schema for user data in API responses"""
    id: int
    created_at: Optional[datetime]
    is_active: bool
    learning_preferences: Optional[Dict[str, Any]] = {}
    total_words_learned: int = 0
    total_reading_time: int = 0  # in minutes
    streak_days: int = 0
    last_activity: Optional[datetime]
    
    class Config:
        from_attributes = True  # Updated for Pydantic v2
        
class UserUpdate(BaseModel):
    """Schema for updating user information"""
    preferred_language: Optional[str] = None
    learning_preferences: Optional[Dict[str, Any]] = None

class Token(BaseModel):
    """Schema for JWT token response"""
    access_token: str
    token_type: str

class TokenData(BaseModel):
    """Schema for token data"""
    username: Optional[str] = None

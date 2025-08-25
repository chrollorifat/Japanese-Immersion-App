from fastapi import APIRouter, Depends, HTTPException, status, UploadFile, File
from sqlalchemy.orm import Session
from typing import List
import os
import shutil

from ..models.database import get_db
from ..models.user import User
from ..models.book import Book
from ..api.auth import get_current_user

router = APIRouter()

@router.get("/", response_model=List[dict])
async def get_user_books(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get all books for the current user"""
    books = db.query(Book).filter(Book.user_id == current_user.id).all()
    return [book.to_dict() for book in books]

@router.post("/upload")
async def upload_book(
    file: UploadFile = File(...),
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Upload a new ebook file"""
    
    # Validate file type
    allowed_types = ["application/epub+zip", "text/plain", "application/pdf"]
    if file.content_type not in allowed_types:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="File type not supported. Please upload EPUB, TXT, or PDF files."
        )
    
    # Create upload directory if it doesn't exist
    upload_dir = "../data/uploads"
    os.makedirs(upload_dir, exist_ok=True)
    
    # Save file
    file_path = os.path.join(upload_dir, f"{current_user.id}_{file.filename}")
    with open(file_path, "wb") as buffer:
        shutil.copyfileobj(file.file, buffer)
    
    # Create book record
    new_book = Book(
        user_id=current_user.id,
        title=file.filename,  # For now, use filename as title
        file_path=file_path,
        file_name=file.filename,
        file_size=os.path.getsize(file_path),
        mime_type=file.content_type,
        processing_status="pending"
    )
    
    db.add(new_book)
    db.commit()
    db.refresh(new_book)
    
    return {
        "message": "Book uploaded successfully",
        "book": new_book.to_dict()
    }

@router.get("/{book_id}")
async def get_book(
    book_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get a specific book"""
    book = db.query(Book).filter(
        Book.id == book_id,
        Book.user_id == current_user.id
    ).first()
    
    if not book:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Book not found"
        )
    
    return book.to_dict()

@router.delete("/{book_id}")
async def delete_book(
    book_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Delete a book"""
    book = db.query(Book).filter(
        Book.id == book_id,
        Book.user_id == current_user.id
    ).first()
    
    if not book:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Book not found"
        )
    
    # Delete file if it exists
    if os.path.exists(book.file_path):
        os.remove(book.file_path)
    
    db.delete(book)
    db.commit()
    
    return {"message": "Book deleted successfully"}

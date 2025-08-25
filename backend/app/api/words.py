from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from ..models.database import get_db
from ..models.user import User
from ..api.auth import get_current_user

router = APIRouter()

@router.get("/lookup/{word}")
async def lookup_word(
    word: str,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Look up a Japanese word in dictionaries"""
    # TODO: Implement dictionary lookup logic
    return {
        "word": word,
        "definitions": [],
        "message": "Dictionary lookup not yet implemented"
    }

@router.post("/mark-known")
async def mark_word_as_known(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Mark a word as known by the user"""
    # TODO: Implement word knowledge tracking
    return {"message": "Word knowledge tracking not yet implemented"}

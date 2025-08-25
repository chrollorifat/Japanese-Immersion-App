from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from ..models.database import get_db
from ..models.user import User
from ..api.auth import get_current_user

router = APIRouter()

@router.get("/due-cards")
async def get_due_cards(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get cards that are due for review"""
    # TODO: Implement SRS card scheduling
    return {
        "due_cards": [],
        "message": "SRS system not yet implemented"
    }

@router.post("/review")
async def review_card(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Submit a card review"""
    # TODO: Implement SRS review logic
    return {"message": "SRS review not yet implemented"}

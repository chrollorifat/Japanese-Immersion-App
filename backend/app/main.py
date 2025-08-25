from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
import os
from dotenv import load_dotenv

from .models.database import create_tables
from .api import auth, books, words, srs, reading

# Load environment variables
load_dotenv()

# Create FastAPI app instance
app = FastAPI(
    title="Japanese Learning Web App",
    description="A comprehensive web application for learning Japanese, similar to LingQ but focused exclusively on Japanese language learning.",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# Configure CORS
origins = [
    "http://localhost:3000",  # Frontend development server
    "http://127.0.0.1:3000",
    "http://localhost:8080",
    os.getenv("FRONTEND_URL", "http://localhost:3000")
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Mount static files (for serving uploaded files and frontend assets)
if not os.path.exists("../data/uploads"):
    os.makedirs("../data/uploads", exist_ok=True)

app.mount("/uploads", StaticFiles(directory="../data/uploads"), name="uploads")
app.mount("/static", StaticFiles(directory="../frontend/static"), name="static")

# Include API routers
app.include_router(auth.router, prefix="/api/auth", tags=["Authentication"])
app.include_router(books.router, prefix="/api/books", tags=["Books"])
app.include_router(words.router, prefix="/api/words", tags=["Words & Dictionary"])
app.include_router(srs.router, prefix="/api/srs", tags=["SRS System"])
app.include_router(reading.router, prefix="/api/reading", tags=["Reading Sessions"])

# Root endpoint
@app.get("/")
async def root():
    """
    Root endpoint - provides basic API information
    """
    return {
        "message": "Japanese Learning Web App API",
        "version": "1.0.0",
        "status": "running",
        "docs": "/docs",
        "redoc": "/redoc"
    }

# Health check endpoint
@app.get("/health")
async def health_check():
    """
    Health check endpoint for monitoring
    """
    return {
        "status": "healthy",
        "message": "API is running normally"
    }

# Startup event
@app.on_event("startup")
async def startup_event():
    """
    Application startup event - create database tables
    """
    print("Starting Japanese Learning Web App...")
    print("Creating database tables...")
    create_tables()
    print("Database tables created successfully!")
    print("Application startup complete.")

# Shutdown event
@app.on_event("shutdown")
async def shutdown_event():
    """
    Application shutdown event - cleanup resources
    """
    print("Shutting down Japanese Learning Web App...")
    # Add any cleanup code here if needed
    print("Shutdown complete.")

if __name__ == "__main__":
    import uvicorn
    
    host = os.getenv("HOST", "0.0.0.0")
    port = int(os.getenv("PORT", 8000))
    debug = os.getenv("DEBUG", "True").lower() == "true"
    
    print(f"Starting server at http://{host}:{port}")
    print(f"Debug mode: {debug}")
    print(f"API documentation: http://{host}:{port}/docs")
    
    uvicorn.run(
        "app.main:app",
        host=host,
        port=port,
        reload=debug,
        log_level="info" if not debug else "debug"
    )

# Development Setup Guide

This guide will help you set up the development environment for the Japanese Learning Web App.

## Prerequisites Installation

### 1. Install Python (3.9 or higher)
- Download from: https://www.python.org/downloads/
- Make sure to check "Add Python to PATH" during installation
- Verify installation: `python --version`

### 2. Install PostgreSQL
- Download from: https://www.postgresql.org/download/
- Remember your PostgreSQL password (you'll need it later)
- Verify installation: `psql --version`

### 3. Install Git
- Download from: https://git-scm.com/downloads
- Verify installation: `git --version`

### 4. Install a Code Editor
- **Recommended**: VS Code (https://code.visualstudio.com/)
- **Alternatives**: PyCharm Community Edition, Sublime Text

## Project Setup

### 1. Clone or Initialize Repository
```bash
# If cloning from GitHub (future)
git clone https://github.com/yourusername/japanese-learning-app.git
cd japanese-learning-app

# If starting locally (current)
cd japanese-learning-app
git init
git add .
git commit -m "Initial commit"
```

### 2. Create Virtual Environment
```bash
# Create virtual environment
python -m venv venv

# Activate virtual environment
# On Windows:
venv\Scripts\activate
# On macOS/Linux:
source venv/bin/activate

# You should see (venv) in your terminal prompt
```

### 3. Install Python Dependencies
```bash
# Make sure you're in the project root and venv is activated
pip install -r backend/requirements.txt
```

### 4. Set Up Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE japanese_learning_app;

# Create user (optional, for better security)
CREATE USER app_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE japanese_learning_app TO app_user;

# Exit PostgreSQL
\q
```

### 5. Configure Environment Variables
```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your actual values:
# - Update DATABASE_URL with your PostgreSQL credentials
# - Add your OpenAI API key (get from https://openai.com)
# - Generate a SECRET_KEY (you can use: python -c "import secrets; print(secrets.token_hex(32))")
```

### 6. Set Up Database Tables
```bash
# Navigate to backend directory
cd backend

# Run database migrations (we'll create these later)
# For now, we'll create tables manually using our models
```

## Running the Application

### 1. Start Backend Server
```bash
# Make sure you're in the backend directory and venv is activated
cd backend
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# The API will be available at: http://localhost:8000
# API documentation: http://localhost:8000/docs
```

### 2. Open Frontend (for now, simple HTML files)
- Open frontend/templates/index.html in your browser
- Or serve with a simple HTTP server:
```bash
cd frontend
python -m http.server 3000
# Visit: http://localhost:3000
```

## Development Workflow

### 1. Daily Development
```bash
# 1. Activate virtual environment
venv\Scripts\activate

# 2. Make your changes

# 3. Test your changes
pytest backend/tests/

# 4. Commit your changes
git add .
git commit -m "Description of changes"
```

### 2. Installing New Packages
```bash
# Install new package
pip install package_name

# Update requirements.txt
pip freeze > backend/requirements.txt
```

## Troubleshooting

### Common Issues

1. **"pip not found"**: Make sure Python is in your PATH
2. **"psql command not found"**: Make sure PostgreSQL bin directory is in PATH
3. **Connection refused to database**: Check if PostgreSQL service is running
4. **Virtual environment issues**: Deactivate and reactivate: `deactivate` then `venv\Scripts\activate`

### Useful Commands
```bash
# Check if PostgreSQL is running (Windows)
net start postgresql

# Check Python packages
pip list

# Check if port is in use
netstat -an | findstr 8000
```

## Next Steps

After completing this setup:
1. Follow the database schema design in `DATABASE_DESIGN.md`
2. Start with the basic API endpoints
3. Build the frontend reader interface
4. Test with sample ebook files

## Learning Resources

### Python & FastAPI
- FastAPI Tutorial: https://fastapi.tiangolo.com/tutorial/
- Python Crash Course (book)
- Real Python website tutorials

### PostgreSQL
- PostgreSQL Tutorial: https://www.postgresqltutorial.com/
- SQLAlchemy documentation

### Frontend Development
- MDN Web Docs: https://developer.mozilla.org/
- JavaScript.info
- Alpine.js documentation

### Japanese Text Processing
- MeCab documentation
- Unicode in Japanese text processing
- Furigana and ruby text implementation

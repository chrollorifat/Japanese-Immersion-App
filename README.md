# Japanese Learning Web App

A comprehensive web application for learning Japanese, similar to LingQ but focused exclusively on Japanese language learning.

## Features (Planned)

### Core Features
- 📚 **Custom Ebook Reader**: Upload and read Japanese ebooks with integrated dictionary
- 📖 **Word Lookup**: Click any word to see definitions, readings, and pronunciations
- 🎯 **SRS System**: Spaced Repetition System for vocabulary learning
- 🌈 **Progress Tracking**: Visual indicators for known/unknown words
- 📊 **Statistics**: Detailed learning progress and analytics

### Dictionary Integration
- Multiple dictionary sources (Sanseido, Daijirin, Ookoku, Jitendex)
- Monolingual and bilingual definitions
- JLPT, WaniKani, and Kanken level indicators
- Audio pronunciations

### AI Features
- Enhanced definitions and explanations
- Context-aware translations
- Learning recommendations
- Difficulty assessment

## Technology Stack

- **Backend**: Python with FastAPI
- **Database**: PostgreSQL
- **Frontend**: HTML/CSS/JavaScript with Alpine.js
- **File Processing**: Python libraries for ebook parsing
- **AI Integration**: OpenAI API
- **Deployment**: Docker (future)

## Getting Started

### Prerequisites
- Python 3.9+
- PostgreSQL
- Git

### Installation
1. Clone the repository
2. Set up virtual environment
3. Install dependencies
4. Configure database
5. Run the application

## Project Structure
```
japanese-learning-app/
├── backend/               # FastAPI backend
│   ├── app/
│   │   ├── api/          # API routes
│   │   ├── models/       # Database models
│   │   ├── services/     # Business logic
│   │   └── utils/        # Utility functions
│   ├── tests/
│   └── requirements.txt
├── frontend/             # Frontend assets
│   ├── static/
│   │   ├── css/
│   │   └── js/
│   └── templates/
├── data/                 # Dictionary data and uploads
├── docs/                 # Documentation
└── docker/               # Docker configuration
```

## Development Phases

1. **Foundation**: Basic project setup and database design
2. **Core Reader**: Ebook upload and basic reading interface
3. **Dictionary Integration**: Word lookup and definitions
4. **SRS System**: Spaced repetition functionality
5. **Advanced Features**: AI integration and analytics
6. **Polish**: UI/UX improvements and optimization

## License
MIT License

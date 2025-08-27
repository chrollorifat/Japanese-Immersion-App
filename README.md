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


### Installation
1. Clone the repository
2. Set up virtual environment
3. Install dependencies
4. Configure database
5. Run the application

# Japanese Learning App - Go Backend

This is the Go backend for the Japanese Learning Web App, rewritten from Python for better performance and simpler dependency management.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ installed
- PostgreSQL installed and running
- Git (optional)

### Installation

1. **Navigate to the Go backend directory:**
   ```bash
   cd backend-go
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your PostgreSQL credentials
   ```

4. **Create PostgreSQL database:**
   ```sql
   CREATE DATABASE japanese_learning_app;
   ```

5. **Run the application:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8000`

## 📁 Project Structure

```
backend-go/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection & migrations
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware (auth, etc.)
│   ├── models/            # Database models
│   └── utils/             # Utility functions
├── .env.example           # Environment variables template
└── go.mod                 # Go module dependencies
```

## 🔧 Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Database
DATABASE_URL=postgres://username:password@localhost:5432/japanese_learning_app?sslmode=disable

# Server
PORT=8000
DEBUG=true

# JWT Secret (change in production!)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

## 📡 API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user (requires auth)

### Books
- `GET /api/books/` - Get user's books (requires auth)
- `POST /api/books/upload` - Upload new book (requires auth)
- `GET /api/books/:id` - Get specific book (requires auth)
- `DELETE /api/books/:id` - Delete book (requires auth)

### Health Check
- `GET /health` - API health status

## 🏗️ Development

### Running in development mode:
```bash
go run main.go
```

### Building for production:
```bash
go build -o japanese-learning-app main.go
./japanese-learning-app
```

### Running tests:
```bash
go test ./...
```

## 🐳 Docker Support (Coming Soon)

## 🛠️ Technologies Used

- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **JWT** - Authentication tokens
- **bcrypt** - Password hashing

## ✨ Why Go?

Advantages over the Python version:
- **Simpler Dependencies**: No virtual environments or complex package management
- **Better Performance**: Compiled binary, faster execution
- **Easy Deployment**: Single binary with no runtime dependencies
- **Type Safety**: Compile-time error checking
- **Concurrency**: Built-in goroutines for handling multiple requests

## 🔄 Migration from Python

The Go version maintains API compatibility with the Python version, so your frontend code should work without changes.

## 📚 Next Steps

1. Install PostgreSQL and configure the database
2. Run the Go backend: `go run main.go`
3. Test with the existing frontend at `../frontend/templates/index.html`
4. Start implementing additional features (word processing, SRS, etc.)


## Development Phases

1. **Foundation**: Basic project setup and database design
2. **Core Reader**: Ebook upload and basic reading interface
3. **Dictionary Integration**: Word lookup and definitions
4. **SRS System**: Spaced repetition functionality
5. **Advanced Features**: AI integration and analytics
6. **Polish**: UI/UX improvements and optimization

## License
MIT License

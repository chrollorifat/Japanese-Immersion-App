# Database Schema Design

## Overview
This document outlines the database schema for the Japanese Learning Web App, designed to support all core features including user management, ebook processing, dictionary lookups, and the SRS system.

## Core Tables

### 1. Users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    preferred_language VARCHAR(10) DEFAULT 'en',
    learning_preferences JSONB DEFAULT '{}',
    
    -- Learning statistics
    total_words_learned INTEGER DEFAULT 0,
    total_reading_time INTEGER DEFAULT 0, -- in minutes
    streak_days INTEGER DEFAULT 0,
    last_activity TIMESTAMP WITH TIME ZONE
);
```

### 2. Books
```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    author VARCHAR(200),
    language VARCHAR(10) DEFAULT 'ja',
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    
    -- Processing status
    processing_status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, failed
    word_count INTEGER DEFAULT 0,
    unique_word_count INTEGER DEFAULT 0,
    difficulty_level VARCHAR(10), -- beginner, intermediate, advanced
    
    -- Metadata
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_read_at TIMESTAMP WITH TIME ZONE,
    reading_progress DECIMAL(5,2) DEFAULT 0.00, -- percentage
    
    -- Content extraction
    extracted_text TEXT,
    chapter_data JSONB DEFAULT '[]', -- [{title, start_pos, end_pos}]
    
    CONSTRAINT valid_progress CHECK (reading_progress >= 0 AND reading_progress <= 100)
);
```

### 3. Words (Dictionary Entries)
```sql
CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    surface_form VARCHAR(100) NOT NULL, -- Original form (e.g., 食べる)
    reading VARCHAR(100), -- Hiragana reading (e.g., たべる)
    pronunciation VARCHAR(100), -- For katakana words
    
    -- Linguistic information
    part_of_speech VARCHAR(50),
    inflection_type VARCHAR(50),
    base_form VARCHAR(100), -- Dictionary form
    
    -- Difficulty indicators
    jlpt_level INTEGER, -- 1-5, null if not in JLPT
    wanikani_level INTEGER, -- 1-60, null if not in WaniKani
    kanken_level INTEGER, -- 1-10, null if not in Kanken
    frequency_rank INTEGER, -- Word frequency ranking
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(surface_form, reading)
);

-- Index for fast lookups
CREATE INDEX idx_words_surface ON words(surface_form);
CREATE INDEX idx_words_reading ON words(reading);
CREATE INDEX idx_words_jlpt ON words(jlpt_level);
```

### 4. Word Definitions
```sql
CREATE TABLE word_definitions (
    id SERIAL PRIMARY KEY,
    word_id INTEGER REFERENCES words(id) ON DELETE CASCADE,
    dictionary_source VARCHAR(50) NOT NULL, -- sanseido, daijirin, ookoku, jitendex, jisho
    language VARCHAR(10) NOT NULL, -- en, ja
    definition TEXT NOT NULL,
    example_sentence TEXT,
    example_translation TEXT,
    
    -- Additional metadata
    definition_order INTEGER DEFAULT 1, -- For multiple definitions
    tags JSONB DEFAULT '[]', -- Additional tags or categories
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_definitions_word ON word_definitions(word_id);
CREATE INDEX idx_definitions_source ON word_definitions(dictionary_source);
```

### 5. User Word Knowledge (Vocabulary Tracking)
```sql
CREATE TABLE user_word_knowledge (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    word_id INTEGER REFERENCES words(id) ON DELETE CASCADE,
    
    -- Knowledge status
    knowledge_level INTEGER DEFAULT 0, -- 0=unknown, 1=recognized, 2=familiar, 3=known, 4=mastered
    first_encountered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_reviewed_at TIMESTAMP WITH TIME ZONE,
    
    -- SRS data
    srs_level INTEGER DEFAULT 0, -- SRS interval level
    next_review_at TIMESTAMP WITH TIME ZONE,
    review_count INTEGER DEFAULT 0,
    correct_count INTEGER DEFAULT 0,
    streak INTEGER DEFAULT 0,
    
    -- Context tracking
    first_seen_book_id INTEGER REFERENCES books(id),
    times_encountered INTEGER DEFAULT 1,
    
    -- Learning metadata
    notes TEXT, -- User's personal notes
    is_ignored BOOLEAN DEFAULT FALSE, -- User marked as "ignore this word"
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, word_id)
);

CREATE INDEX idx_user_words_user ON user_word_knowledge(user_id);
CREATE INDEX idx_user_words_next_review ON user_word_knowledge(next_review_at);
CREATE INDEX idx_user_words_level ON user_word_knowledge(knowledge_level);
```

### 6. Reading Sessions
```sql
CREATE TABLE reading_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    
    start_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER, -- calculated from start/end time
    
    -- Progress tracking
    start_position INTEGER DEFAULT 0, -- Character position in text
    end_position INTEGER DEFAULT 0,
    words_learned INTEGER DEFAULT 0,
    words_reviewed INTEGER DEFAULT 0,
    
    -- Session data
    new_words_encountered JSONB DEFAULT '[]', -- Array of word IDs
    words_looked_up JSONB DEFAULT '[]', -- Array of word IDs
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user ON reading_sessions(user_id);
CREATE INDEX idx_sessions_book ON reading_sessions(book_id);
```

### 7. SRS Cards
```sql
CREATE TABLE srs_cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    word_id INTEGER REFERENCES words(id) ON DELETE CASCADE,
    
    -- Card type and content
    card_type VARCHAR(20) DEFAULT 'recognition', -- recognition, recall, production
    front_content JSONB NOT NULL, -- {text, furigana, audio_url}
    back_content JSONB NOT NULL, -- {definition, example, notes}
    
    -- SRS algorithm data
    ease_factor DECIMAL(4,2) DEFAULT 2.50, -- Anki-style ease factor
    interval_days INTEGER DEFAULT 1,
    repetition_count INTEGER DEFAULT 0,
    
    -- Scheduling
    due_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_reviewed_at TIMESTAMP WITH TIME ZONE,
    
    -- Performance tracking
    total_reviews INTEGER DEFAULT 0,
    correct_reviews INTEGER DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    
    -- Card state
    is_suspended BOOLEAN DEFAULT FALSE,
    is_buried BOOLEAN DEFAULT FALSE, -- Temporarily hidden
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, word_id, card_type)
);

CREATE INDEX idx_srs_cards_due ON srs_cards(due_date, is_suspended);
CREATE INDEX idx_srs_cards_user ON srs_cards(user_id);
```

### 8. Review History
```sql
CREATE TABLE review_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    card_id INTEGER REFERENCES srs_cards(id) ON DELETE CASCADE,
    
    reviewed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    response_quality INTEGER NOT NULL, -- 1=again, 2=hard, 3=good, 4=easy
    response_time_ms INTEGER, -- Time taken to answer
    
    -- Pre-review state
    old_interval INTEGER,
    old_ease_factor DECIMAL(4,2),
    
    -- Post-review state
    new_interval INTEGER,
    new_ease_factor DECIMAL(4,2),
    new_due_date TIMESTAMP WITH TIME ZONE,
    
    -- Context
    review_context VARCHAR(50), -- 'reader', 'srs_session', 'manual'
    device_type VARCHAR(20), -- 'desktop', 'mobile', 'tablet'
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_review_history_user ON review_history(user_id);
CREATE INDEX idx_review_history_date ON review_history(reviewed_at);
```

### 9. Book Annotations
```sql
CREATE TABLE book_annotations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    
    -- Position in text
    start_position INTEGER NOT NULL,
    end_position INTEGER NOT NULL,
    selected_text TEXT NOT NULL,
    
    -- Annotation data
    annotation_type VARCHAR(20) DEFAULT 'word_lookup', -- word_lookup, note, highlight
    annotation_data JSONB DEFAULT '{}', -- Flexible data storage
    
    -- Visual styling
    highlight_color VARCHAR(7), -- Hex color code
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_annotations_book ON book_annotations(book_id);
CREATE INDEX idx_annotations_position ON book_annotations(start_position, end_position);
```

## Views and Derived Data

### Learning Progress View
```sql
CREATE VIEW user_learning_progress AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(DISTINCT uwk.word_id) as total_words_encountered,
    COUNT(DISTINCT CASE WHEN uwk.knowledge_level >= 3 THEN uwk.word_id END) as words_known,
    COUNT(DISTINCT CASE WHEN uwk.next_review_at <= CURRENT_TIMESTAMP THEN sc.id END) as cards_due,
    COALESCE(AVG(rs.duration_minutes), 0) as avg_session_minutes,
    u.streak_days,
    u.total_reading_time
FROM users u
LEFT JOIN user_word_knowledge uwk ON u.id = uwk.user_id
LEFT JOIN srs_cards sc ON u.id = sc.user_id AND sc.is_suspended = FALSE
LEFT JOIN reading_sessions rs ON u.id = rs.user_id
GROUP BY u.id, u.username, u.streak_days, u.total_reading_time;
```

## Data Relationships Summary

1. **Users** ↔ **Books**: One-to-many (user uploads multiple books)
2. **Books** ↔ **Reading Sessions**: One-to-many (multiple sessions per book)
3. **Words** ↔ **Word Definitions**: One-to-many (multiple definitions per word)
4. **Users** ↔ **User Word Knowledge**: One-to-many (user knows many words)
5. **Users** ↔ **SRS Cards**: One-to-many (user has many cards)
6. **SRS Cards** ↔ **Review History**: One-to-many (card reviewed many times)
7. **Books** ↔ **Book Annotations**: One-to-many (multiple annotations per book)

## Migration Strategy

1. Start with core tables: `users`, `books`, `words`
2. Add learning tracking: `user_word_knowledge`, `reading_sessions`
3. Implement SRS: `srs_cards`, `review_history`
4. Add advanced features: `word_definitions`, `book_annotations`

## Performance Considerations

- Indexes on frequently queried columns
- Partitioning for large tables (review_history by date)
- Regular cleanup of old session data
- Caching for dictionary lookups
- Materialized views for complex statistics

## Security Considerations

- Row-level security for user data
- Encrypted storage for sensitive information
- Audit trails for data modifications
- Backup and recovery procedures

// API Configuration
const API_BASE = 'http://localhost:8000/api';

// Alpine.js data and methods
function appData() {
    return {
        // State
        isLoggedIn: false,
        user: null,
        token: null,
        books: [],
        showLogin: false,
        showRegister: false,
        
        // Forms
        loginForm: {
            username: '',
            password: ''
        },
        registerForm: {
            username: '',
            email: '',
            password: ''
        },
        
        // Notifications
        notification: {
            show: false,
            message: '',
            type: 'info'
        },

        // Initialize app
        async init() {
            // Check for existing token
            this.token = localStorage.getItem('token');
            if (this.token) {
                try {
                    const response = await fetch(`${API_BASE}/auth/me`, {
                        headers: {
                            'Authorization': `Bearer ${this.token}`
                        }
                    });
                    
                    if (response.ok) {
                        this.user = await response.json();
                        this.isLoggedIn = true;
                        await this.loadBooks();
                    } else {
                        // Token is invalid
                        localStorage.removeItem('token');
                        this.token = null;
                    }
                } catch (error) {
                    console.error('Failed to validate token:', error);
                    localStorage.removeItem('token');
                    this.token = null;
                }
            }
        },

        // Authentication
        async login() {
            try {
                const formData = new FormData();
                formData.append('username', this.loginForm.username);
                formData.append('password', this.loginForm.password);

                const response = await fetch(`${API_BASE}/auth/token`, {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    const data = await response.json();
                    this.token = data.access_token;
                    localStorage.setItem('token', this.token);
                    
                    // Get user info
                    const userResponse = await fetch(`${API_BASE}/auth/me`, {
                        headers: {
                            'Authorization': `Bearer ${this.token}`
                        }
                    });
                    
                    if (userResponse.ok) {
                        this.user = await userResponse.json();
                        this.isLoggedIn = true;
                        this.showLogin = false;
                        this.showNotification('Login successful!', 'success');
                        await this.loadBooks();
                        
                        // Reset form
                        this.loginForm = { username: '', password: '' };
                    }
                } else {
                    const error = await response.json();
                    this.showNotification(error.detail || 'Login failed', 'error');
                }
            } catch (error) {
                console.error('Login error:', error);
                this.showNotification('Login failed. Please try again.', 'error');
            }
        },

        async register() {
            try {
                const response = await fetch(`${API_BASE}/auth/register`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.registerForm)
                });

                if (response.ok) {
                    this.showNotification('Registration successful! Please login.', 'success');
                    this.showRegister = false;
                    this.showLogin = true;
                    
                    // Reset form
                    this.registerForm = { username: '', email: '', password: '' };
                } else {
                    const error = await response.json();
                    this.showNotification(error.detail || 'Registration failed', 'error');
                }
            } catch (error) {
                console.error('Registration error:', error);
                this.showNotification('Registration failed. Please try again.', 'error');
            }
        },

        logout() {
            this.isLoggedIn = false;
            this.user = null;
            this.token = null;
            this.books = [];
            localStorage.removeItem('token');
            this.showNotification('Logged out successfully', 'info');
        },

        // Books management
        async loadBooks() {
            if (!this.token) return;
            
            try {
                const response = await fetch(`${API_BASE}/books/`, {
                    headers: {
                        'Authorization': `Bearer ${this.token}`
                    }
                });
                
                if (response.ok) {
                    this.books = await response.json();
                } else {
                    console.error('Failed to load books');
                }
            } catch (error) {
                console.error('Error loading books:', error);
            }
        },

        async uploadBook(event) {
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            try {
                this.showNotification('Uploading book...', 'info');
                
                const response = await fetch(`${API_BASE}/books/upload`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${this.token}`
                    },
                    body: formData
                });

                if (response.ok) {
                    const result = await response.json();
                    this.showNotification(result.message, 'success');
                    await this.loadBooks();
                    
                    // Reset file input
                    event.target.value = '';
                } else {
                    const error = await response.json();
                    this.showNotification(error.detail || 'Upload failed', 'error');
                }
            } catch (error) {
                console.error('Upload error:', error);
                this.showNotification('Upload failed. Please try again.', 'error');
            }
        },

        async deleteBook(bookId) {
            if (!confirm('Are you sure you want to delete this book?')) return;

            try {
                const response = await fetch(`${API_BASE}/books/${bookId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${this.token}`
                    }
                });

                if (response.ok) {
                    this.showNotification('Book deleted successfully', 'success');
                    await this.loadBooks();
                } else {
                    const error = await response.json();
                    this.showNotification(error.detail || 'Delete failed', 'error');
                }
            } catch (error) {
                console.error('Delete error:', error);
                this.showNotification('Delete failed. Please try again.', 'error');
            }
        },

        openReader(book) {
            // TODO: Implement book reader
            this.showNotification('Reader functionality coming soon!', 'info');
        },

        // Notifications
        showNotification(message, type = 'info') {
            this.notification = {
                show: true,
                message,
                type
            };
            
            // Auto-hide after 5 seconds
            setTimeout(() => {
                this.hideNotification();
            }, 5000);
        },

        hideNotification() {
            this.notification.show = false;
        }
    }
}

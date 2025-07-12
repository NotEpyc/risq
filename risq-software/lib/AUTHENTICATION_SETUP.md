# Authentication with Railways Server

This project has been updated to use a Railways server for authentication instead of Firebase.

## Changes Made

1. **Removed Firebase Dependencies**
   - Removed Firebase initialization from `main.dart`
   - Removed Firebase Auth code from `signin_page.dart`
   - Removed Firebase Auth code from `signup_page.dart`

2. **Added HTTP Package**
   - Added `http: ^1.1.0` to `pubspec.yaml` for API calls

3. **Created AuthService**
   - New service class at `lib/services/auth_service.dart`
   - Handles login and registration API calls
   - Centralized error handling

4. **Updated UI**
   - Replaced FontAwesome icons with Material Design icons
   - Updated sign-in page to use AuthService
   - Updated sign-up page to use AuthService

## Railways Server Setup

You'll need to update the base URL in `lib/services/auth_service.dart`:

```dart
static const String _baseUrl = 'https://your-railways-app.railway.app/api';
```

### Expected API Endpoints

Your Railways server should provide these endpoints:

#### POST `/api/auth/login`
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "token": "jwt_token_here",
  "user": {
    "id": "user_id",
    "name": "User Name",
    "email": "user@example.com"
  }
}
```

**Error Response (401/400):**
```json
{
  "success": false,
  "message": "Invalid credentials"
}
```

#### POST `/api/auth/register`
**Request Body:**
```json
{
  "name": "User Name",
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (201/200):**
```json
{
  "success": true,
  "token": "jwt_token_here",
  "user": {
    "id": "user_id",
    "name": "User Name",
    "email": "user@example.com"
  }
}
```

## Server-Side Implementation Requirements

### 1. Database Setup
You'll need to set up a database (PostgreSQL, MySQL, or MongoDB) with a users table/collection:

**Users Table Schema (SQL):**
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Users Collection Schema (MongoDB):**
```javascript
{
  _id: ObjectId,
  name: String,
  email: String (unique),
  password_hash: String,
  created_at: Date,
  updated_at: Date
}
```

### 2. Environment Variables
Set up these environment variables in your Railways deployment:

```bash
# Database
DATABASE_URL=your_database_connection_string

# JWT
JWT_SECRET=your_super_secure_jwt_secret_key
JWT_EXPIRES_IN=7d

# Server
PORT=3000
NODE_ENV=production

# CORS (if needed)
ALLOWED_ORIGINS=*
```

### 3. Required Dependencies (Node.js Example)
```json
{
  "dependencies": {
    "express": "^4.18.2",
    "bcryptjs": "^2.4.3",
    "jsonwebtoken": "^9.0.2",
    "cors": "^2.8.5",
    "helmet": "^7.0.0",
    "express-rate-limit": "^6.10.0",
    "joi": "^17.9.2",
    "pg": "^8.11.3",
    "sequelize": "^6.32.1"
  }
}
```

### 4. Login Endpoint Implementation
```javascript
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const { User } = require('../models');

app.post('/api/auth/login', async (req, res) => {
  try {
    const { email, password } = req.body;

    // Validation
    if (!email || !password) {
      return res.status(400).json({
        success: false,
        message: 'Email and password are required'
      });
    }

    // Find user by email
    const user = await User.findOne({ where: { email } });
    if (!user) {
      return res.status(401).json({
        success: false,
        message: 'Invalid credentials'
      });
    }

    // Verify password
    const isValidPassword = await bcrypt.compare(password, user.password_hash);
    if (!isValidPassword) {
      return res.status(401).json({
        success: false,
        message: 'Invalid credentials'
      });
    }

    // Generate JWT token
    const token = jwt.sign(
      { userId: user.id, email: user.email },
      process.env.JWT_SECRET,
      { expiresIn: process.env.JWT_EXPIRES_IN || '7d' }
    );

    // Return success response
    res.status(200).json({
      success: true,
      token,
      user: {
        id: user.id,
        name: user.name,
        email: user.email
      }
    });

  } catch (error) {
    console.error('Login error:', error);
    res.status(500).json({
      success: false,
      message: 'Internal server error'
    });
  }
});
```

### 5. Register Endpoint Implementation
```javascript
app.post('/api/auth/register', async (req, res) => {
  try {
    const { name, email, password } = req.body;

    // Validation
    if (!name || !email || !password) {
      return res.status(400).json({
        success: false,
        message: 'Name, email, and password are required'
      });
    }

    // Email format validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      return res.status(400).json({
        success: false,
        message: 'Please enter a valid email address'
      });
    }

    // Password strength validation
    if (password.length < 6) {
      return res.status(400).json({
        success: false,
        message: 'Password must be at least 6 characters long'
      });
    }

    // Check if user already exists
    const existingUser = await User.findOne({ where: { email } });
    if (existingUser) {
      return res.status(400).json({
        success: false,
        message: 'An account already exists for this email'
      });
    }

    // Hash password
    const saltRounds = 12;
    const password_hash = await bcrypt.hash(password, saltRounds);

    // Create user
    const user = await User.create({
      name,
      email,
      password_hash
    });

    // Generate JWT token
    const token = jwt.sign(
      { userId: user.id, email: user.email },
      process.env.JWT_SECRET,
      { expiresIn: process.env.JWT_EXPIRES_IN || '7d' }
    );

    // Return success response
    res.status(201).json({
      success: true,
      token,
      user: {
        id: user.id,
        name: user.name,
        email: user.email
      }
    });

  } catch (error) {
    console.error('Registration error:', error);
    
    // Handle database constraint errors
    if (error.name === 'SequelizeUniqueConstraintError') {
      return res.status(400).json({
        success: false,
        message: 'An account already exists for this email'
      });
    }

    res.status(500).json({
      success: false,
      message: 'Internal server error'
    });
  }
});
```

### 6. Security Middleware
```javascript
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
const cors = require('cors');

// Security headers
app.use(helmet());

// CORS configuration
app.use(cors({
  origin: process.env.ALLOWED_ORIGINS?.split(',') || '*',
  credentials: true
}));

// Rate limiting
const authLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // Limit each IP to 5 requests per windowMs
  message: {
    success: false,
    message: 'Too many authentication attempts, please try again later'
  },
  standardHeaders: true,
  legacyHeaders: false,
});

// Apply rate limiting to auth routes
app.use('/api/auth', authLimiter);
```

### 7. Input Validation Middleware
```javascript
const Joi = require('joi');

const validateLogin = (req, res, next) => {
  const schema = Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(1).required()
  });

  const { error } = schema.validate(req.body);
  if (error) {
    return res.status(400).json({
      success: false,
      message: error.details[0].message
    });
  }
  next();
};

const validateRegister = (req, res, next) => {
  const schema = Joi.object({
    name: Joi.string().min(1).max(255).required(),
    email: Joi.string().email().required(),
    password: Joi.string().min(6).required()
  });

  const { error } = schema.validate(req.body);
  if (error) {
    return res.status(400).json({
      success: false,
      message: error.details[0].message
    });
  }
  next();
};

// Use middleware
app.post('/api/auth/login', validateLogin, loginHandler);
app.post('/api/auth/register', validateRegister, registerHandler);
```

### 8. JWT Verification Middleware (for protected routes)
```javascript
const verifyToken = (req, res, next) => {
  const token = req.header('Authorization')?.replace('Bearer ', '');

  if (!token) {
    return res.status(401).json({
      success: false,
      message: 'Access denied. No token provided'
    });
  }

  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    req.user = decoded;
    next();
  } catch (error) {
    res.status(401).json({
      success: false,
      message: 'Invalid token'
    });
  }
};
```

### 9. Basic Server Setup
```javascript
const express = require('express');
const app = express();

// Body parsing middleware
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'OK', timestamp: new Date().toISOString() });
});

// Error handling middleware
app.use((error, req, res, next) => {
  console.error('Unhandled error:', error);
  res.status(500).json({
    success: false,
    message: 'Internal server error'
  });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    success: false,
    message: 'Route not found'
  });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
```

### 10. Testing Your Endpoints

You can test your endpoints using curl or Postman:

**Test Registration:**
```bash
curl -X POST https://your-app.railway.app/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Test Login:**
```bash
curl -X POST https://your-app.railway.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Next Steps

1. Deploy your backend API to Railways
2. Update the `_baseUrl` in `AuthService` with your Railways app URL
3. Implement token storage (optional - see commented code in AuthService)

## Token Storage (Optional)

To implement persistent authentication, you can:

1. Add `shared_preferences` package to `pubspec.yaml`
2. Uncomment and implement the token storage methods in `AuthService`
3. Check for stored tokens on app startup

```yaml
dependencies:
  shared_preferences: ^2.2.2
```

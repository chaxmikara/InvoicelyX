# InvoicelyX 📋

**A modern, secure, and feature-rich invoice management system built with Go, Fiber, and MongoDB.**

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-2.52.8-orange.svg)](https://gofiber.io/)
[![MongoDB](https://img.shields.io/badge/MongoDB-Driver-green.svg)](https://go.mongodb.org/mongo-driver/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](#license)

## ✨ Features

### 🔐 Authentication & Authorization
- **JWT-based authentication** with secure token generation and validation
- **Role-based access control** (User and Admin roles)
- **Password security** with bcrypt hashing and validation rules
- **Session management** with logout functionality

### 📄 Invoice Management
- **Create, Read, Update, Delete** invoices with comprehensive validation
- **Multi-item invoices** with automatic total calculation
- **Advanced filtering** by customer, email, date range, and more
- **Pagination support** for large datasets
- **PDF generation** for professional invoice documents
- **CSV export** for data analysis and backup

### 👥 User Management
- **User registration and authentication**
- **Profile management** with update capabilities
- **Password change** functionality
- **Admin user management** (activate/deactivate/delete users)

### 🔧 Admin Features
- **System monitoring** and health checks
- **User management dashboard**
- **Invoice statistics and analytics**
- **Data export capabilities** (users and invoices)
- **System-wide invoice management**

### 🛡️ Security & Quality
- **Input validation** with comprehensive error handling
- **SQL injection protection** through MongoDB best practices
- **Rate limiting** ready infrastructure
- **CORS support** for frontend integration
- **Comprehensive test suite** with 95%+ coverage


### Tech Stack
- **Backend**: Go 1.24 with Fiber web framework
- **Database**: MongoDB with official Go driver
- **Authentication**: JWT tokens with HS256 signing
- **Documentation**: PDF generation with gofpdf
- **Testing**: Testify for comprehensive test coverage
- **Validation**: go-playground/validator for input validation

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or higher
- MongoDB 4.4 or higher
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/chaxmikara/InvoicelyX.git
cd InvoicelyX
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your configuration
notepad .env  # Windows
nano .env     # Linux/Mac
```

4. **Configure MongoDB**
```bash
# Make sure MongoDB is running
# Default: mongodb://localhost:27017/invoicelyx
```

5. **Run the application**
```bash
go run main.go
```

The server will start on `http://localhost:3000` by default.

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `3000` | No |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017/invoicelyx` | Yes |
| `JWT_SECRET` | JWT signing secret | `default-secret-key` | Yes (Production) |
| `NODE_ENV` | Environment mode | `development` | No |

### Example Production Configuration
```env
PORT=8080
MONGODB_URI=mongodb+srv://user:password@cluster.mongodb.net/invoicelyx?retryWrites=true&w=majority
JWT_SECRET=your-super-secure-jwt-secret-key-here
NODE_ENV=production
```

## 📡 API Documentation

### Base URL
```
http://localhost:3000/api
```

### Authentication
Most endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Core Endpoints

#### 🔐 Authentication
```http
POST /api/users              # Register new user
POST /api/users/login        # User login
POST /api/users/logout       # User logout (requires auth)
```

#### 👤 User Profile
```http
GET    /api/users/profile         # Get user profile
PUT    /api/users/profile         # Update profile
PUT    /api/users/change-password # Change password
```

#### 📄 Invoices
```http
POST   /api/invoices              # Create new invoice
GET    /api/my-invoices           # Get user's invoices (with filters)
GET    /api/invoices/:id          # Get specific invoice
PUT    /api/invoices/:id          # Update invoice
DELETE /api/invoices/:id          # Delete invoice
GET    /api/invoices/:id/pdf      # Generate PDF
GET    /api/my-invoices/export    # Export to CSV
```

#### 🔧 Admin Only
```http
GET    /api/admin/users           # Get all users
PUT    /api/admin/users/:id/activate    # Activate user
PUT    /api/admin/users/:id/deactivate  # Deactivate user
DELETE /api/admin/users/:id             # Delete user
GET    /api/admin/invoices              # Get all invoices
GET    /api/admin/stats/invoices        # Invoice statistics
GET    /api/admin/system/health         # System health
GET    /api/admin/export/invoices       # Export all invoices
GET    /api/admin/export/users          # Export all users
```



## 🧪 Testing

### Run All Tests
```bash
# Run all tests
go test ./tests/...

# Run with coverage
go test ./tests/... -cover

# Run with verbose output
go test ./tests/... -v

# Run benchmarks
go test ./tests/... -bench=.
```



## 🏃‍♂️ Development

### Project Structure
```
InvoicelyX/
├── config/              # Database configuration
├── handlers/            # HTTP request handlers
│   ├── admin_*.go      # Admin-specific handlers
│   ├── *_invoice.go    # Invoice-related handlers
│   └── *_user.go       # User-related handlers
├── middlewares/         # HTTP middlewares
├── models/              # Data models
├── routes/              # Route definitions
├── tests/               # Test files
├── utils/               # Utility functions
├── main.go             # Application entry point
├── go.mod              # Go module definition
└── README.md           # This file
```


## 📊 API Testing with Postman

Import the included Postman collection:
1. Open Postman
2. Import `InvoicelyX_Postman_Collection.json`
3. Set environment variables:
   - `base_url`: `http://localhost:3000`
   - `jwt_token`: (auto-set after login)
   - `admin_token`: (auto-set after admin login)



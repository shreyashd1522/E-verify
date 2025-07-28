# E-verify 🔐

A modern, secure full-stack authentication system built with **Go** backend and **React** frontend. Features comprehensive user authentication with email verification, password reset functionality, and a beautiful animated UI.

![E-verify](https://img.shields.io/badge/Go-1.24.5-blue)
![React](https://img.shields.io/badge/React-18.3.1-blue)
![MongoDB](https://img.shields.io/badge/MongoDB-5.0+-green)
![Material-UI](https://img.shields.io/badge/Material--UI-7.2.0-orange)

## ✨ Features

### 🔒 Security Features
- **SHA-256 Password Hashing** - Secure password storage using Go's crypto library
- **Email Verification** - Gmail SMTP integration with expiring verification tokens
- **Password Reset** - Secure token-based password reset functionality
- **CORS Protection** - Cross-origin request handling for frontend-backend communication
- **Token Expiration** - Time-based token expiration for enhanced security

### 🎨 Frontend Features
- **Modern React UI** - Built with React 18 and Material-UI 7
- **Animated Components** - Smooth animations using Framer Motion
- **Responsive Design** - Mobile-friendly interface
- **Lottie Animations** - Beautiful loading and success/error animations
- **React Router** - Client-side routing for seamless navigation

### 🚀 Backend Features
- **Go HTTP Server** - High-performance backend with Go 1.24.5
- **MongoDB Integration** - NoSQL database for user data storage
- **RESTful API** - Clean API endpoints for authentication operations
- **Template Rendering** - HTML templates for email verification pages
- **Environment Configuration** - Secure configuration management

## 🏗️ Architecture

```
E-verify/
├── frontend/                 # React frontend application
│   ├── public/              # Static assets and animations
│   ├── src/                 # React source code
│   │   ├── components/      # Reusable UI components
│   │   ├── pages/          # Page components
│   │   └── App.js          # Main application component
│   └── package.json        # Frontend dependencies
├── templates/               # HTML templates for email verification
├── main.go                 # Go backend server
├── go.mod                  # Go module dependencies
└── config.env              # Environment configuration
```

## 🛠️ Tech Stack

### Backend
- **Go 1.24.5** - High-performance programming language
- **MongoDB** - NoSQL database for user data
- **Gmail SMTP** - Email service for verification
- **SHA-256** - Password hashing algorithm
- **JWT-like tokens** - Secure token generation

### Frontend
- **React 18.3.1** - Modern JavaScript framework
- **Material-UI 7.2.0** - React UI component library
- **Framer Motion 12.23.9** - Animation library
- **React Router DOM 7.7.0** - Client-side routing
- **React Lottie Player 2.1.0** - Lottie animation support

## 🚀 Quick Start

### Prerequisites
- Go 1.24.5 or higher
- Node.js 16 or higher
- MongoDB instance
- Gmail account with App Password

### 1. Clone the Repository
```bash
git clone https://github.com/shreyashd1522/E-verify.git
cd E-verify
```

### 2. Backend Setup

#### Install Go Dependencies
```bash
go mod download
```

#### Configure Environment Variables
Create a `config.env` file in the root directory:
```env
GMAIL_USER=your-email@gmail.com
GMAIL_PASS=your-app-password
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
MONGODB_URI=mongodb://localhost:27017
```

#### Run the Backend Server
```bash
go run main.go
```
The server will start on `http://localhost:8080`

### 3. Frontend Setup

#### Install Dependencies
```bash
cd frontend
npm install
```

#### Start Development Server
```bash
npm start
```
The React app will start on `http://localhost:3000`

## 📋 API Endpoints

### Authentication Endpoints
- `POST /verify` - User registration and email verification
- `GET /verify` - Email verification page
- `POST /check-verified` - Check user verification status
- `POST /forgot-password` - Initiate password reset
- `POST /reset-password` - Complete password reset
- `POST /resend-verification` - Resend verification email

### Frontend Routes
- `/` - Main application
- `/verify` - Email verification page
- `/forgot-password` - Password reset request
- `/reset-password` - Password reset form
- `/resend-verification` - Resend verification email

## 🔧 Configuration

### Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `GMAIL_USER` | Gmail email address | Required |
| `GMAIL_PASS` | Gmail app password | Required |
| `SMTP_SERVER` | SMTP server address | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |

### Gmail Setup
1. Enable 2-factor authentication on your Gmail account
2. Generate an App Password for this application
3. Use the App Password in the `GMAIL_PASS` environment variable

## 🎯 Features in Detail

### User Registration
- Email and password validation
- SHA-256 password hashing
- Verification token generation
- Email verification sent via Gmail SMTP

### Email Verification
- Secure token-based verification
- Token expiration (24 hours)
- Resend verification functionality
- Beautiful verification success/error animations

### Password Reset
- Secure token generation
- Email-based password reset
- Token expiration handling
- Password strength validation

### Frontend Animations
- Loading animations with Lottie
- Success/error state animations
- Smooth page transitions
- Material-UI component animations

## 🔒 Security Considerations

- **Password Hashing**: SHA-256 hashing for password storage
- **Token Security**: Cryptographically secure random token generation
- **Token Expiration**: Time-based token expiration
- **CORS Protection**: Proper CORS headers for frontend-backend communication
- **Input Validation**: Server-side validation for all inputs
- **Error Handling**: Secure error messages without information leakage

## 🧪 Testing

### Backend Testing
```bash
go test ./...
```

### Frontend Testing
```bash
cd frontend
npm test
```

## 📦 Deployment

### Backend Deployment
1. Build the Go binary: `go build -o e-verify main.go`
2. Set up environment variables on your server
3. Run the binary: `./e-verify`

### Frontend Deployment
1. Build the React app: `npm run build`
2. Deploy the `build` folder to your web server
3. Configure your web server to serve the React app

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Material-UI](https://mui.com/) for the beautiful UI components
- [Framer Motion](https://www.framer.com/motion/) for smooth animations
- [MongoDB](https://www.mongodb.com/) for the database
- [Go](https://golang.org/) for the high-performance backend

## 📞 Support

If you have any questions or need help with the project, please open an issue on GitHub.

---

**Made with ❤️ by [shreyashd1522](https://github.com/shreyashd1522)**

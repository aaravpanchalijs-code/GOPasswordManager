# 🔐 SecureVault

A secure full-stack Password Manager built using **Go**, **MongoDB**, **JWT Authentication**, and **AES-GCM Encryption**.

SecureVault allows users to safely store, retrieve, update, delete, and share passwords while ensuring sensitive data remains encrypted in the database.

---

## Features

- User Signup & Login
- JWT Authentication
- Password Hashing using bcrypt
- AES-GCM Encryption for stored passwords
- Add Passwords
- View Stored Passwords
- Show/Hide Passwords
- Edit Passwords
- Delete Passwords
- Share Passwords with another user
- Search Stored Passwords
- Dashboard Statistics
- Modern Responsive UI

---

## Tech Stack

### Backend
- Go
- MongoDB Atlas
- JWT
- bcrypt
- AES-GCM Encryption

### Frontend
- HTML5
- CSS3
- JavaScript (Vanilla)

---

## 📂 Project Structure

```
secure-password-manager/
│
├── cmd/
│   └── main.go
│
├── database/
│   └── mongo.go
│
├── handlers/
│   ├── auth.go
│   └── vault.go
│
├── middleware/
│   └── auth.go
│
├── models/
│   ├── user.go
│   └── vault.go
│
├── utils/
│   ├── encrypt.go
│   ├── jwt.go
│   └── response.go
│
├── frontend/
│   ├── css/
│   ├── js/
│   ├── assets/
│   ├── login.html
│   ├── signup.html
│   └── dashboard.html
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/aaravpanchalijs-code/GOPasswordManager.git
```

```bash
cd GOPasswordManager
```

---

### 2. Install dependencies

```bash
go mod tidy
```

---

### 3. Create a `.env`

Create a `.env` file in the project root.

Example:

```env
MONGO_URI=your_mongodb_connection_string

JWT_SECRET=your_jwt_secret

ENCRYPTION_KEY=your_32_byte_encryption_key
```

---

### 4. Run the server

```bash
go run ./cmd
```

Server starts on

```
http://localhost:8080
```

---

## 🔐 Security

- Passwords are hashed using **bcrypt**
- Stored passwords are encrypted using **AES-GCM**
- JWT-based authentication protects all vault routes
- Middleware validates every protected request

---

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|---------|----------|-------------|
| POST | `/signup` | Register user |
| POST | `/login` | Login user |

---

### Vault

| Method | Endpoint | Description |
|---------|----------|-------------|
| POST | `/vault/add` | Add Password |
| GET | `/vault/get` | Get Passwords |
| PUT | `/vault/update` | Update Password |
| DELETE | `/vault/delete` | Delete Password |

---

## Learning Objectives

This project was built to learn:

- Backend Development in Go
- REST API Design
- MongoDB Integration
- Authentication & Authorization
- Cryptography
- Middleware
- Frontend-Backend Communication
- Secure Password Storage

---


## Author

**Aarav Panchal**

GitHub: https://github.com/aaravpanchalijs-code


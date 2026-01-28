# Secure-Question-Paper-DIstribution-Portal

A command-line application implementing comprehensive security concepts for secure exam paper management, developed as part of FOUNDATIONS OF CYBER SECURITY lab evaluation.

## Project Overview

This system demonstrates the practical implementation of core security principles including authentication, authorization, encryption, digital signatures, and secure encoding for managing exam question papers in an academic environment.

## Security Features Implemented

### 1. Authentication (Multi-Factor)
- Single-factor authentication using username and password
- Multi-factor authentication with time-limited OTP (One-Time Password)
- Password hashing using bcrypt with custom salt
- SHA-256 pre-hashing to handle unlimited password lengths
- OTP validity period: 5 minutes
- OTP single-use enforcement

### 2. Authorization (Access Control)
- Role-Based Access Control (RBAC) with three roles: Faculty, Exam Cell, Student
- Access Control Matrix implementation with granular permissions
- Permission enforcement before all sensitive operations
- Audit logging for security-critical actions

### 3. Encryption (Hybrid Approach)
- AES-256-GCM encryption for question paper content
- RSA-2048 for secure key exchange
- Random AES key generation per document
- AES key encrypted with recipient's RSA public key
- Hybrid encryption combining speed of AES with security of RSA

### 4. Digital Signatures
- SHA-256 hashing of document content
- RSA-based digital signature creation using faculty's private key
- Signature verification using faculty's public key
- Integrity and authenticity verification
- Tamper detection capability

### 5. Encoding
- Base64 encoding for binary-safe database storage
- Encoding applied to encrypted content, keys, and signatures
- Clear distinction between encoding (format conversion) and encryption (confidentiality)

## System Architecture

### Roles and Permissions

**Faculty:**
- Upload question papers
- Encrypt papers automatically during upload
- Sign papers with private key
- View own uploaded papers

**Exam Cell:**
- View all encrypted question papers
- Decrypt papers using private key
- Verify digital signatures
- Manage exam sessions

**Student:**
- View exam schedule (read-only access)
- No access to question papers
- No decryption capabilities

### Access Control Matrix

| Role      | Question Paper    | Encryption Key | Exam Session |
|-----------|-------------------|----------------|--------------|
| Faculty   | Create, Encrypt   | Generate       | View         |
| Exam Cell | Read, Decrypt     | Decrypt        | Manage       |
| Student   | None              | None           | View         |

## Technical Stack

- Language: Go 1.21+
- Database: MySQL 8.0+
- Cryptography: Go standard library (crypto/*)
- Password Hashing: bcrypt with cost factor 12
- Symmetric Encryption: AES-256-GCM
- Asymmetric Encryption: RSA-2048
- Hash Function: SHA-256
- Encoding: Base64 (standard encoding)

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- GCC (for MySQL driver compilation)

## Installation

### 1. Clone Repository
```bash
git clone https://github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal.git
cd Secure-Question-Paper-Distribution-Portal
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3.Configure Database and env

Setup mysql server and enter the details in the .env file as follows :

```bash
DB_USER=userName
DB_PASSWORD=password
DB_HOST=host (i.e, localhost)
DB_PORT=portNo
DB_NAME=databaseName
# optional
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_USER=your-email@gmail.com
# SMTP_PASS=your-app-password
# SMTP_FROM=youremail@gmail.com
```

## Usage Flow

### First-Time Setup

1. Start the application
2. Register users for each role:
   - Faculty user (receives RSA key pair automatically)
   - Exam Cell user (receives RSA key pair automatically)
   - Student user (no keys required)

### Faculty Workflow

1. Login with credentials
2. Complete MFA with OTP
3. Upload question paper:
   - Provide paper title and subject
   - Specify exam date
   - Enter file path (PDF or TXT)
4. System automatically:
   - Generates random AES-256 key
   - Encrypts paper with AES-GCM
   - Encrypts AES key with Exam Cell's RSA public key
   - Creates digital signature with faculty's RSA private key
   - Encodes all data in Base64
   - Stores in database

### Exam Cell Workflow

1. Login with credentials
2. Complete MFA with OTP
3. View all encrypted papers
4. Select paper to decrypt
5. System automatically:
   - Retrieves encrypted paper and key
   - Decodes Base64 data
   - Decrypts AES key using Exam Cell's RSA private key
   - Decrypts paper content using AES key
   - Verifies digital signature using faculty's RSA public key
   - Displays decrypted content if signature valid

### Student Workflow

1. Login with credentials
2. Complete MFA with OTP
3. View exam schedule (limited access)
4. Access to question papers blocked by ACL

## Database Schema

### Key Tables

**users**: Stores user credentials, roles, and RSA keys
**otp_sessions**: Manages OTP tokens for MFA
**question_papers**: Stores encrypted papers, keys, and signatures
**exam_sessions**: Manages exam scheduling
**access_control**: Defines ACL permissions
**audit_log**: Tracks security-relevant actions

## Security Considerations

### Password Security
- Bcrypt hashing with cost factor 12
- Random 32-byte salt per user
- SHA-256 pre-hashing to bypass bcrypt's 72-byte limit
- No plaintext passwords stored

### Key Management
- RSA-2048 key pairs generated during registration
- Private keys stored in PEM format (in production, encrypt these)
- Public keys distributed for encryption and verification
- Separate key pairs for Faculty and Exam Cell roles

### Encryption Details
- AES-256-GCM provides authenticated encryption
- Unique AES key per document
- Nonce generated using crypto/rand
- RSA PKCS1v15 for key encryption

### Digital Signature Process
1. Compute SHA-256 hash of plaintext document
2. Sign hash with faculty's RSA private key using PKCS1v15
3. Signature verified during decryption
4. Failed verification indicates tampering

### Attack Mitigation

| Attack Type         | Mitigation Strategy                                    |
|---------------------|--------------------------------------------------------|
| Brute Force         | bcrypt adaptive cost + MFA                             |
| Rainbow Tables      | Random salt per user                                   |
| Man-in-the-Middle   | End-to-end encryption                                  |
| Tampering           | Digital signatures with integrity verification         |
| Unauthorized Access | ACL enforcement with database-level permissions        |
| Key Compromise      | Role-based key separation                              |
| Replay Attacks      | OTP expiration and single-use enforcement              |

## Testing

### Test Scenarios

1. User registration with password validation
2. Login with invalid credentials (should fail)
3. Login with correct password but wrong OTP (should fail)
4. Faculty paper upload and encryption
5. Exam Cell paper decryption and verification
6. Student access denial (ACL enforcement)
7. Signature verification after data tampering (should fail)


## Project Structure
```
secure-exam-system/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/
│   │   ├── registration.go     # User registration logic
│   │   ├── login.go            # Authentication and MFA
│   │   └── otp.go              # OTP generation and verification
│   ├── crypto/
│   │   ├── aes.go              # AES encryption/decryption
│   │   ├── rsa.go              # RSA key operations
│   │   ├── hashing.go          # Password hashing
│   │   ├── signature.go        # Digital signatures
│   │   └── encoding.go         # Base64 encoding
│   ├── acl/
│   │   └── permissions.go      # Access control logic
│   ├── models/
│   │   └── user.go             # Data models
│   ├── database/
│   │   ├── db.go               # Database connection
│   │   └── schema.go           # Schema definitions
│   └── services/
│       └── paper_service.go    # Business logic
├── pkg/
│   ├── email/
│   │   └── sender.go           # Email simulation
│   └── utils/
│       └── input.go            # User input utilities
├── storage/
│   ├── keys/                   # RSA key storage
│   └── exam.db                 # Database file
├── .env                        # Environment configuration
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
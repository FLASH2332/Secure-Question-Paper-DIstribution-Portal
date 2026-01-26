package database

const Schema = `
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(64) NOT NULL,
    role ENUM('Faculty', 'ExamCell', 'Student') NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    public_key TEXT,
    private_key_encrypted TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- OTP sessions table
CREATE TABLE IF NOT EXISTS otp_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_otp (user_id, is_used)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Question papers table
CREATE TABLE IF NOT EXISTS question_papers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    subject VARCHAR(100) NOT NULL,
    faculty_id INT NOT NULL,
    encrypted_content TEXT NOT NULL,
    encrypted_aes_key TEXT NOT NULL,
    digital_signature TEXT NOT NULL,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    exam_date DATE,
    status ENUM('pending', 'approved', 'published') DEFAULT 'pending',
    FOREIGN KEY (faculty_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_faculty (faculty_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Exam sessions table
CREATE TABLE IF NOT EXISTS exam_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    paper_id INT NOT NULL,
    session_name VARCHAR(100) NOT NULL,
    scheduled_time DATETIME NOT NULL,
    duration_minutes INT NOT NULL,
    status ENUM('scheduled', 'active', 'completed') DEFAULT 'scheduled',
    created_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (paper_id) REFERENCES question_papers(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_paper (paper_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Access control list
CREATE TABLE IF NOT EXISTS access_control (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role ENUM('Faculty', 'ExamCell', 'Student') NOT NULL,
    object_type ENUM('QuestionPaper', 'EncryptionKey', 'ExamSession') NOT NULL,
    can_create BOOLEAN DEFAULT FALSE,
    can_read BOOLEAN DEFAULT FALSE,
    can_update BOOLEAN DEFAULT FALSE,
    can_delete BOOLEAN DEFAULT FALSE,
    can_encrypt BOOLEAN DEFAULT FALSE,
    can_decrypt BOOLEAN DEFAULT FALSE,
    UNIQUE KEY unique_role_object (role, object_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Audit log
CREATE TABLE IF NOT EXISTS audit_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    action VARCHAR(100) NOT NULL,
    object_type VARCHAR(50) NOT NULL,
    object_id INT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    success BOOLEAN NOT NULL,
    details TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_action (user_id, action),
    INDEX idx_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

const ACLData = `
-- Insert default ACL permissions

-- Faculty permissions
INSERT INTO access_control (role, object_type, can_create, can_read, can_update, can_delete, can_encrypt, can_decrypt) VALUES
('Faculty', 'QuestionPaper', TRUE, TRUE, FALSE, FALSE, TRUE, FALSE),
('Faculty', 'EncryptionKey', TRUE, FALSE, FALSE, FALSE, FALSE, FALSE),
('Faculty', 'ExamSession', FALSE, TRUE, FALSE, FALSE, FALSE, FALSE);

-- ExamCell permissions
INSERT INTO access_control (role, object_type, can_create, can_read, can_update, can_delete, can_encrypt, can_decrypt) VALUES
('ExamCell', 'QuestionPaper', FALSE, TRUE, TRUE, FALSE, FALSE, TRUE),
('ExamCell', 'EncryptionKey', FALSE, FALSE, FALSE, FALSE, FALSE, TRUE),
('ExamCell', 'ExamSession', TRUE, TRUE, TRUE, TRUE, FALSE, FALSE);

-- Student permissions (very limited)
INSERT INTO access_control (role, object_type, can_create, can_read, can_update, can_delete, can_encrypt, can_decrypt) VALUES
('Student', 'QuestionPaper', FALSE, FALSE, FALSE, FALSE, FALSE, FALSE),
('Student', 'EncryptionKey', FALSE, FALSE, FALSE, FALSE, FALSE, FALSE),
('Student', 'ExamSession', FALSE, TRUE, FALSE, FALSE, FALSE, FALSE);
`

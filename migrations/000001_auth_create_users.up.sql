CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);
-- Create unique index on the email for gorm:uniqueIndex
CREATE UNIQUE INDEX idx_users_email ON users(email);
-- Create index for deleted_at for gorm:index
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

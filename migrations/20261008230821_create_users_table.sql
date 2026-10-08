-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    role_id INTEGER NOT NULL DEFAULT 2, -- Default to a 'User' role
    deleted_at DATETIME DEFAULT NULL,
    --Foreign Key Constraints
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);
-- Create unique index on the email for gorm:uniqueIndex
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
-- Create index for deleted_at for gorm:index
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

INSERT OR IGNORE INTO users (id,first_name,last_name,email,role_id) VALUES
(1,'Admin','User','admin@example.com',1);
-- +goose Down
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP TABLE IF EXISTS users;
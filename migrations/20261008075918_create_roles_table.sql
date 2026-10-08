-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);

INSERT OR IGNORE INTO roles (id,role_name) VALUES
(1,'Super Admin'),
(2,'Admin'),
(3,'User');
-- +goose Down
DROP INDEX IF EXISTS idx_roles_deleted_at;
DROP TABLE IF EXISTS roles;
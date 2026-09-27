CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);
--Seed the default roles
INSERT OR IGNORE INTO roles(id,role_name) VALUES
(1,'Admin'),
(2,'User'),
(3,'Guest');

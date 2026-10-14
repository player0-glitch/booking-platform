-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    first_name VARCHAR(255) NULL,
    last_name VARCHAR(255) NULL,
    email VARCHAR(255) NULL UNIQUE,
    password_hash TEXT NULL, -- nullable so we can have temporary guest users
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    role_id INTEGER NOT NULL DEFAULT 3, -- Default to a 'Customer' role
    deleted_at TIMESTAMPTZ NULL,
    --Foreign Key Constraints
    CONSTRAINT fk_users_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);
-- Create unique index on the email for gorm:uniqueIndex

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;

-- Create index for deleted_at for gorm:index
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
--
-- INSERT INTO users (first_name,last_name,email,role_id) VALUES
-- ('Super','User','super@example.com',1),
-- ('Admin','User','admin@example.com',2),
-- ('Customer','User','customer@example.com',3),
-- ('Guest','User','guest@example.com',4)
-- ON CONFLICT (email) WHERE deleted_at IS NULL DO NOTHING;
--
-- +goose Down
DROP INDEX IF EXISTS idx_users_roles;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP TABLE IF EXISTS users;

-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    first_name VARCHAR(255) NULL,
    last_name VARCHAR(255) NULL,
    email VARCHAR(255) NULL,
    account_status VARCHAR(20) NOT NULL DEFAULT 'registerd'
        CHECK( account_status IN ('registered','suspended','guest')),
    password_hash TEXT NULL, -- nullable so we can have temporary guest users
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- role_id INTEGER NOT NULL DEFAULT 3, -- Default to a 'Customer' role
    deleted_at TIMESTAMPTZ NULL
    --Foreign Key Constraints
    -- CONSTRAINT users_account_status_check 
        -- CHECK( account_status IN ('registered','suspended','guest'))
);

-- Create index for deleted_at for gorm:index
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

CREATE UNIQUE INDEX uq_user_email_active ON users (LOWER(email))
    WHERE email IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- Anonymous sessions allows for us to have a 'Guest' identity
-- like a window shoppe
-- ============================================
-- +goose Down
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS uq_user_email_active; 
DROP TABLE IF EXISTS users;

-- +goose Up
-- ============================================
-- Platform user (each entry) can be any orgranisation, hairstylist, freelancer
-- restaurant. Everything safely leads back here allowing this 
-- booking engine to be used by several types businesses
-- timezones will be calculated from UTC offset
-- ============================================
CREATE TABLE IF NOT EXISTS organisations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE, -- preferably URL friendly
    timezone TEXT NOT NULL DEFAULT 'UTC',
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active','inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP 
);
-- +goose Down
DROP TABLE IF EXISTS organisations;
-- +goose Up
-- ============================================
-- THE PERSON BEING SERVED, this is scoped to an organisation. A freelance
-- developer is not going to be serving a haircut customer
-- some business have customer that will never ever be users (like a hairt stylist)
-- so 'user_id' will eventually be null. 'user_id' can be linked to a customer
-- once login later
-- ============================================
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    -- give this customer an identity on the platform
    user_id UUID NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email TEXT,
    phone VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    -- we need to be able to reach out to a customer
    CONSTRAINT customers_contactable 
        CHECK (email IS NOT NUTT OR phone IS NOT NULL)
);
--One customer per email per business, case-insensitive. 
CREATE UNIQUE INDEX customers_org_email_uniq ON
    customers (organisation_id,lower(email)) WHERE email IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_user_id_active ON 
    customers (user_id)
    WHERE user_id NOT NULL AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers 
    (deleted_at);
-- +goose Down
DROP INDEX IF NOT EXISTS customers_org_email_uniq;
DROP INDEX IF NOT EXISTS idx_customers_deleted_at;
DROP INDEX IF NOT EXISTS idx_customers_user_id_active;
DROP TABLE IF NOT EXISTS customers;
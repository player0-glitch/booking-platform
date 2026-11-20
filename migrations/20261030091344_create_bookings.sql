-- +goose Up
-- ===========================================
-- What the customer seees; "Thandi, Haircut, Saturday 10:00-11:30,confirmed,
-- R350"
-- This is the root of the booking engine (allocations, payments, reschedules, events) hang off of this
-- 'start_time'/'ends_at': Customer facing times and do not include buffer times
-- time locked to a booking can be found in allocations
-- 'total_price_minor': is a SNAPSHOT, changing a service's price tomorrow
-- should not override yesterday's bookings
-- 'version': optimistic lock
-- ===========================================

CREATE TABLE IF NOT EXISTS bookings(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    organisation_id UUID NOT NULL REFERENCES organisations (id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    service_id UUID NOT NULL REFERENCES services(id),
    start_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    duration_minutes INT NOT NULL 
        CHECK(duration_minutes > 0),
    -- becomes useful when booking a hotel room or restaurant table
    party_size INT NOT NULL DEFAULT 1
        CHECK (party_size >= 1),
    -- not every organisations will use all these
    status TEXT  NOT NULL 
        CHECK (status IN ('pending','confirmed','declined','cancelled',
                'completed','no_show','expired')),
    notes TEXT NOT NULL DEFAULT '',
    total_price_minor INT NOT NULL
        CHECK(total_price_minor >=0),
    currency CHAR(3) NOT NULL ,
    version int NOT NULL DEFAULT 1 CHECK (version >= 1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ,

-- because time moves forward duh
    CONSTRAINT booking_time_order CHECK (ends_at> start_at)
);
CREATE INDEX IF NOT EXISTS  bookings_org_start_idx ON bookings (organisation_id,start_at);

CREATE INDEX IF NOT EXISTS  bookings_customer_starts_idx ON bookings (customer_id,start_at);
-- +goose Down
DROP INDEX IF EXISTS bookings_org_start_idx ;
DROP INDEX IF EXISTS bookings_customer_starts_idx;
DROP TABLE IF EXISTS bookings;
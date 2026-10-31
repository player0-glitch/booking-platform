-- +goose Up
-- ============================================
-- Resource is 'the thing that gets booked' (restaurant table, hair appointment)
-- things of that nature.
-- 'booking_mode': 
--              'exclusive'-> 1 booking at a time. (stylist, table, hotel room)
--              'shared'   -> many bookings can overlap creating a 'crowd' 
--              situation NOTE!  this may not get implemented
-- 'capacity': maximum people booking to be part of a group
-- ============================================
CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    organisation_id UUID NOT NULL REFERENCES organisations (id),
    kind TEXT NOT NULL,
    name TEXT NOT NULL,
    booking_mode TEXT NOT NULL DEFAULT 'exclusive' 
            -- NO PLANS TO IMPLEMENT THIS
            CHECK (booking_mode in ('exclusive','shared')),
    capacity INT NOT NULL DEFAULT 1 CHECK (capacity>=1),
    status TEXT NOT NULL DEFAULT 'active'
            CHECK (status in ('active','inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (organisation_id,name)
);

-- +goose Down
DROP TABLE IF EXISTS resources;
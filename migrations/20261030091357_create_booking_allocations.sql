-- +goose Up
-- ===========================================
-- The time that is actually locked int on a specific resource, buffers included
-- this is where 'no overlap' is enforced at a db level
-- Separation because;
--  * one booking can lock several resources (stylist+chair etc)
-- * the locked range can be wider than the booked range (buffers)
-- * a hold can be released or moved without touching what the customer sees.
-- We have half-open ranges [blocked_from,blocked_until) that way
-- 10:00-11:00 and 11:00-12:30 do not collide
-- 'exclusive': copied from resources.booking_mode at insert time,
--          because a constraint cannot look at another table.
-- 'status': 'held'-> (pending approval) and 'active' block the slot.
--       'released'-> (cancelled/declined/expired) frees it but keeps the history:
-- ===========================================
-- teaches GiST indexes to handle plain equality (resource_id WITH = )
-- alongside range overlap (&&)
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS booking_allocations(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    booking_id UUID NOT NULL 
        REFERENCES bookings(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources (id),
    blocked_from TIMESTAMPTZ NOT NULL,
    blocked_until TIMESTAMPTZ NOT NULL,
    exclusive boolean NOT NULL,
    quantity INT NOT NULL DEFAULT 1
        CHECK (quantity >= 1),
    status TEXT NOT NULL 
        CHECK(status IN ('held','active','released')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT booking_allocations_valid_range 
      CHECK (blocked_from < blocked_until),
-- RULE:"the same resourcre cannot be in two places at once."
-- Constraint says: reject any row that, compared with another row, has the 
-- SAME resource_id AND an OVERLAPPINGA time_range, but only for rows that are exclusive
-- and still holding the slot;
-- WHERE (...) = only rows that currently block the slot. 'released' rows
-- (cancelled/declined/expired) are history and never collidd.
-- The constraint creates a GiST index that also speeds up "is this range free?" queries.
--
  CONSTRAINT booking_allocation_no_overlap 
    EXCLUDE using gist(
      resource_id WITH =,
      tstzrange(blocked_from,blocked_until,'[)') WITH &&
      )
      WHERE (exclusive AND status IN ('held','active'))
); 

CREATE INDEX IF NOT EXISTS  booking_allocations_booking_idx ON booking_allocations (booking_id);
-- +goose Down
DROP INDEX IF EXISTS  booking_allocations_booking_idx;
DROP TABLE IF EXISTS booking_allocations;
--Leaving btree_gist extensions installed for later use

-- +goose U
-- ============================================
-- If a Resource is 'what/who' is being booked, a Service is for 'how long/
-- for how much?'. This can be used for variable-length bookings
-- 'duration_mode':
--      'fixed' -> fixed amount of time (1hr for a haircut,2 hrs for a 
-- consultation)
--      'variable' -> customer can pick slices of time for the service
--       between a min and max (freelancer can do: 1hr,2hr,3hrs)
-- 'buffer' -> cleanup time. THIS IS NOT CHARGED TO THE CUSTOMER
-- 'requires_approval: some booking require confirmation or cancellation by
-- the service provider
-- ============================================
CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    organisation_id UUID NOT NULL REFERENCES organisations (id),
    name TEXT NOT NULL,
    duration_mode TEXT NOT NULL DEFAULT 'fixed'
        CHECK (duration_mode IN ('fixed','variable')),
    duration_minutes INT NOT NULL 
        CHECK (duration_minutes > 0),
    min_duration_minutes INT,
    max_duration_minutes INT,
    duration_step INT,
    -- can depend on organisation
    buffer_before_minute INT NOT NULL DEFAULT 0
        CHECK (buffer_before_minute >= 0),
    buffer_after_minute INT NOT NULL DEFAULT 0
        CHECK (buffer_after_minute>= 0),

    -- this is supposed to represent money
    price_minor  INT NOT NULL DEFAULT 0 
        CHECK (price >= 0),
    pricing_unit TEXT NOT NULL DEFAULT 'per_booking'
        CHECK(pricing_unit IN 'per_hour','per_booking')
    currency CHAR(3) NOT NULL DEFAULT 'ZAR',
    requires_approval BOOLEAN NOT NULL DEFAULT false,
    status TEXT NOT NULL DEFAULT 'active'
        CHECK(status IN 'active','inactive'),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    UNIQUE(organisation_id,name),
-- A service with a variable time needs bounds [start_time,end_time)
-- therefore DB will not store variable service without bounds
CONSTRAINT services_variable_needs_bounds 
    CHECK ( duration_mode = 'fixed' OR (
        min_duration_minutes IS NOT NULL AND
        max_duration_minutes IS NOT NULL AND
        duration_step IS NOT NULL AND duration_step > 0 AND
        min_duration_minutes <= max_duration_minutes
    ))
);

-- ============================================
-- Resource_service: Defines "who can do what".Joining Table
-- A hair stylist can perform many services (nails, hair, makeup) and many 
-- services can be performed by many hair stylists
-- 'price_override_minor': let's a higher ranking stylist to charge 'their price'
-- This elimates duplicating a db service entry all because of a new price
CREATE TABLE IF NOT EXISTS resource_services (
    resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    service_id  UUID REFERENCES services(id) ON DELETE CASCADE
    price_override_minor INT 
        CHECK( price_override_minor >= 0),
    PRIMARY KEY (resource_id,service_id)
);
-- PRIMARY KEY only covers 'the services for a resource'
-- This covers "resources for a service"
CREATE INDEX resource_services_service_idx ON resource_services (service_id);
-- ============================================
-- +goose Down
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS resource_services;

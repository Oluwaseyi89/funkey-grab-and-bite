-- address keeps the street line; the rest of the postal address gets its own
-- columns so clients can format it and build map links and structured data.
ALTER TABLE business_settings
    ADD COLUMN IF NOT EXISTS city VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS state VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS postal_code VARCHAR(20) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS country VARCHAR(100) NOT NULL DEFAULT 'Nigeria';

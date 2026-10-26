ALTER TABLE business_settings
    DROP COLUMN IF EXISTS country,
    DROP COLUMN IF EXISTS postal_code,
    DROP COLUMN IF EXISTS state,
    DROP COLUMN IF EXISTS city;

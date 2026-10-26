-- Lets the admin switch pay on delivery (and pay at pickup) on or off. It
-- defaults to on because cash orders were always accepted before.
ALTER TABLE business_settings
    ADD COLUMN IF NOT EXISTS pay_on_delivery_enabled BOOLEAN NOT NULL DEFAULT true;

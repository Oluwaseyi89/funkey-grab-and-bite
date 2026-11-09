ALTER TABLE orders
    DROP COLUMN IF EXISTS payment_review_reason,
    DROP COLUMN IF EXISTS payment_received_kobo;

DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS payment_references;

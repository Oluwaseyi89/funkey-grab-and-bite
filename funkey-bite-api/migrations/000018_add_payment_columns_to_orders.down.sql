DROP INDEX IF EXISTS idx_orders_payment_reference;

ALTER TABLE orders
    DROP COLUMN IF EXISTS payment_paid_at,
    DROP COLUMN IF EXISTS payment_expires_at,
    DROP COLUMN IF EXISTS payment_bank_name,
    DROP COLUMN IF EXISTS payment_account_name,
    DROP COLUMN IF EXISTS payment_account_number,
    DROP COLUMN IF EXISTS payment_reference,
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS payment_method;

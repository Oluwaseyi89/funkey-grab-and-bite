ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_method VARCHAR(20) NOT NULL DEFAULT 'cash',
    ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) NOT NULL DEFAULT 'not_required',
    ADD COLUMN IF NOT EXISTS payment_reference VARCHAR(100),
    ADD COLUMN IF NOT EXISTS payment_account_number VARCHAR(20),
    ADD COLUMN IF NOT EXISTS payment_account_name VARCHAR(200),
    ADD COLUMN IF NOT EXISTS payment_bank_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS payment_expires_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS payment_paid_at TIMESTAMP;

-- Paystack references must be unique so a webhook can never be misattributed
-- to the wrong order; NULL is allowed (cash orders never get one).
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_payment_reference
    ON orders(payment_reference) WHERE payment_reference IS NOT NULL;

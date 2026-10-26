DROP INDEX IF EXISTS idx_orders_payment_status_method;
DROP INDEX IF EXISTS idx_orders_invoice_number;

ALTER TABLE orders
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS payment_collected_via,
    DROP COLUMN IF EXISTS discount_amount,
    DROP COLUMN IF EXISTS tax_amount,
    DROP COLUMN IF EXISTS delivery_fee,
    DROP COLUMN IF EXISTS subtotal_amount,
    DROP COLUMN IF EXISTS invoice_number,
    DROP COLUMN IF EXISTS delivery_address;

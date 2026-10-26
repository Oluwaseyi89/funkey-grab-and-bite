-- delivery_address: where a delivery order goes (required for delivery orders).
-- invoice_number: raised for pay-on-delivery orders.
-- subtotal/delivery fee/tax/discount: the breakdown behind total_amount, kept
-- so an invoice shows what was charged even after the fee settings change.
-- payment_collected_via: how a pay-on-delivery payment was collected.
-- updated_at: status updates and cancellations already write it, but the
-- column was never created, so those updates failed.
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS delivery_address TEXT,
    ADD COLUMN IF NOT EXISTS invoice_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS subtotal_amount DECIMAL(10,2),
    ADD COLUMN IF NOT EXISTS delivery_fee DECIMAL(10,2),
    ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(10,2),
    ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(10,2),
    ADD COLUMN IF NOT EXISTS payment_collected_via VARCHAR(20),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_invoice_number
    ON orders(invoice_number) WHERE invoice_number IS NOT NULL;

-- The admin "awaiting payment" view filters on these.
CREATE INDEX IF NOT EXISTS idx_orders_payment_status_method
    ON orders(payment_status, payment_method);

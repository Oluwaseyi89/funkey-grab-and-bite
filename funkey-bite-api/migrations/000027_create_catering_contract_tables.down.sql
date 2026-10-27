DROP INDEX IF EXISTS idx_catering_requests_invoice_number;
ALTER TABLE catering_requests
    DROP COLUMN IF EXISTS paid_in_full_at,
    DROP COLUMN IF EXISTS contracted_at,
    DROP COLUMN IF EXISTS invoice_number,
    DROP COLUMN IF EXISTS accepted_quote_id;

DROP TABLE IF EXISTS catering_status_history;
DROP TABLE IF EXISTS catering_receipts;
DROP TABLE IF EXISTS catering_invoices;
DROP TABLE IF EXISTS catering_payment_charges;
DROP TABLE IF EXISTS catering_payments;
DROP TABLE IF EXISTS catering_quotes;

ALTER TABLE business_settings DROP COLUMN IF EXISTS catering_terms;

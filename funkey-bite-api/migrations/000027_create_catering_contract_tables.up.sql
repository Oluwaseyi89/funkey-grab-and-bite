-- Catering contracts: the admin sends quotes, the customer accepts one, and
-- payments are recorded until the booking is paid in full. Invoices and the
-- receipt are stored as issued so they can be reproduced exactly.

-- Default terms copied onto each new quote; the admin can edit them per quote.
ALTER TABLE business_settings
    ADD COLUMN IF NOT EXISTS catering_terms TEXT NOT NULL DEFAULT
'1. Booking: your booking is confirmed once the deposit shown on your invoice is paid.
2. Balance: the remaining balance is due before the event date unless agreed otherwise in writing.
3. Changes: changes to the number of guests or the menu must be requested at least 72 hours before the event and may change the price.
4. Cancellation: if you cancel 7 or more days before the event, payments are refunded less any costs already incurred; cancellations within 7 days of the event are not refundable.
5. Venue: you provide access to the venue, a suitable serving area and power where needed.
6. Events outside our control: we are not liable for delays or losses caused by events outside our reasonable control.';

CREATE TABLE IF NOT EXISTS catering_quotes (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL REFERENCES catering_requests(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    line_items JSONB NOT NULL DEFAULT '[]',
    deposit_percent INTEGER NOT NULL CHECK (deposit_percent BETWEEN 1 AND 100),
    terms TEXT NOT NULL,
    message TEXT,
    valid_until DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'sent'
        CHECK (status IN ('sent', 'accepted', 'superseded', 'withdrawn')),
    -- SHA-256 of the one-time link token emailed with the quote.
    acceptance_token_hash VARCHAR(64),
    accepted_at TIMESTAMP,
    accepted_via VARCHAR(20) CHECK (accepted_via IN ('email_link', 'status_page')),
    created_by INTEGER REFERENCES admin_users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (request_id, version)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_catering_quotes_one_accepted
    ON catering_quotes(request_id) WHERE status = 'accepted';
CREATE UNIQUE INDEX IF NOT EXISTS idx_catering_quotes_token
    ON catering_quotes(acceptance_token_hash) WHERE acceptance_token_hash IS NOT NULL;

-- The payment ledger: amounts paid and balances are always summed from here.
CREATE TABLE IF NOT EXISTS catering_payments (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL REFERENCES catering_requests(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    method VARCHAR(20) NOT NULL
        CHECK (method IN ('cash', 'bank_transfer', 'paystack_transfer', 'card')),
    reference VARCHAR(100),
    note TEXT,
    recorded_by INTEGER REFERENCES admin_users(id) ON DELETE SET NULL,
    paid_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_catering_payments_request ON catering_payments(request_id);
-- Online payments are confirmed by webhooks that may be delivered twice.
CREATE UNIQUE INDEX IF NOT EXISTS idx_catering_payments_online_reference
    ON catering_payments(reference) WHERE method IN ('paystack_transfer', 'card');

-- Online payment requests (Paystack dedicated accounts) awaiting the webhook.
CREATE TABLE IF NOT EXISTS catering_payment_charges (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL REFERENCES catering_requests(id) ON DELETE CASCADE,
    reference VARCHAR(100) NOT NULL UNIQUE,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    account_number VARCHAR(20),
    account_name VARCHAR(200),
    bank_name VARCHAR(100),
    expires_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Every invoice version as issued: on acceptance, then after each part payment.
CREATE TABLE IF NOT EXISTS catering_invoices (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL REFERENCES catering_requests(id) ON DELETE CASCADE,
    quote_id INTEGER NOT NULL REFERENCES catering_quotes(id),
    invoice_number VARCHAR(40) NOT NULL,
    version INTEGER NOT NULL,
    reason VARCHAR(20) NOT NULL CHECK (reason IN ('accepted', 'payment')),
    total DECIMAL(12,2) NOT NULL,
    deposit_due DECIMAL(12,2) NOT NULL,
    amount_paid DECIMAL(12,2) NOT NULL,
    balance DECIMAL(12,2) NOT NULL,
    issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (invoice_number, version)
);

-- One receipt per booking, issued when it is paid in full.
CREATE TABLE IF NOT EXISTS catering_receipts (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL UNIQUE REFERENCES catering_requests(id) ON DELETE CASCADE,
    receipt_number VARCHAR(40) NOT NULL UNIQUE,
    total_paid DECIMAL(12,2) NOT NULL,
    issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Who changed a request's status, when and why.
CREATE TABLE IF NOT EXISTS catering_status_history (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL REFERENCES catering_requests(id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    actor_type VARCHAR(20) NOT NULL CHECK (actor_type IN ('admin', 'customer', 'system')),
    actor_id INTEGER,
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_catering_status_history_request ON catering_status_history(request_id);

ALTER TABLE catering_requests
    ADD COLUMN IF NOT EXISTS accepted_quote_id INTEGER REFERENCES catering_quotes(id),
    ADD COLUMN IF NOT EXISTS invoice_number VARCHAR(40),
    ADD COLUMN IF NOT EXISTS contracted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS paid_in_full_at TIMESTAMP;
CREATE UNIQUE INDEX IF NOT EXISTS idx_catering_requests_invoice_number
    ON catering_requests(invoice_number) WHERE invoice_number IS NOT NULL;

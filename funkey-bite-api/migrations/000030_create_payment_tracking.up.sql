-- Payment tracking for online transfers.

-- Every Paystack reference issued for an order. A customer who got new
-- payment details may still pay into an earlier account, so a payment is
-- matched to its order by any of them, not just the latest.
CREATE TABLE IF NOT EXISTS payment_references (
    reference VARCHAR(100) PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_references_order_id ON payment_references (order_id);

INSERT INTO payment_references (reference, order_id)
SELECT payment_reference, id FROM orders WHERE payment_reference IS NOT NULL
ON CONFLICT DO NOTHING;

-- Paystack webhook events as received. An event is processed once: a repeat
-- delivery of the same event for the same reference is recognised and
-- ignored. Outcome records what the event did.
CREATE TABLE IF NOT EXISTS payment_events (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL DEFAULT 'paystack',
    event VARCHAR(50) NOT NULL,
    reference VARCHAR(100) NOT NULL,
    amount_kobo BIGINT,
    currency VARCHAR(10),
    outcome VARCHAR(30) NOT NULL,
    order_id INTEGER REFERENCES orders(id) ON DELETE SET NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider, event, reference)
);

CREATE INDEX IF NOT EXISTS idx_payment_events_received_at ON payment_events (received_at DESC);

-- A payment that arrived but can't simply be accepted (the wrong amount, or
-- after the order was cancelled for non-payment) waits for the admin with
-- payment_status 'needs_review'; these record what was received and why.
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_received_kobo BIGINT,
    ADD COLUMN IF NOT EXISTS payment_review_reason VARCHAR(50);

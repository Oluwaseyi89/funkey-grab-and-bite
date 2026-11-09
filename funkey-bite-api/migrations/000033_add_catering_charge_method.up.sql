-- Catering online payments can be made by card in Paystack's popup as well
-- as by bank transfer. Card charges carry no account details and are never
-- reused, since the popup's access code isn't stored.
ALTER TABLE catering_payment_charges
    ADD COLUMN IF NOT EXISTS method VARCHAR(20) NOT NULL DEFAULT 'paystack_transfer'
        CHECK (method IN ('paystack_transfer', 'card'));

-- One-time codes emailed to customers: to prove they own the email on a
-- guest record before an account takes over its order history, and to reset
-- a forgotten password. Only a keyed hash of each code is stored; a code
-- expires, can be used once, and allows a few wrong tries.
CREATE TABLE IF NOT EXISTS account_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    purpose VARCHAR(30) NOT NULL CHECK (purpose IN ('claim_account', 'password_reset')),
    email VARCHAR(200) NOT NULL,
    code_hash VARCHAR(64) NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_account_codes_user_purpose ON account_codes (user_id, purpose, created_at DESC);

-- Customers are matched by their number whatever its format (0803..., +234803...),
-- using its last ten digits.
CREATE INDEX IF NOT EXISTS idx_users_phone_digits ON users (right(regexp_replace(phone, '\D', '', 'g'), 10));

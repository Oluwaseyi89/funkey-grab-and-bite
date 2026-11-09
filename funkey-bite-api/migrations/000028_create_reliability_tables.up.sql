-- Reliability foundations: idempotent client requests, an append-only audit
-- log and rate-limit counters shared by every API instance.

-- One row per client request made with an Idempotency-Key. A retry with the
-- same key replays the stored result instead of repeating the change.
CREATE TABLE IF NOT EXISTS idempotency_keys (
    -- Method and route, e.g. "POST /api/v1/orders", so a key only ever
    -- matches the endpoint it was first used on.
    scope VARCHAR(200) NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL,
    -- SHA-256 of the request body; reusing a key for a different request is rejected.
    request_hash VARCHAR(64) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'in_progress'
        CHECK (status IN ('in_progress', 'completed')),
    -- Set inside the transaction that makes the change, so a key is only
    -- committed when the change is.
    entity_type VARCHAR(50),
    entity_id VARCHAR(100),
    response_status INTEGER,
    response_body BYTEA,
    -- A request still in progress after this time crashed before committing,
    -- so a retry may take it over.
    locked_until TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (scope, idempotency_key)
);

CREATE INDEX IF NOT EXISTS idx_idempotency_keys_expires_at ON idempotency_keys (expires_at);

-- Who changed what and when. Rows are only ever inserted, in the same
-- transaction as the change they describe.
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_type VARCHAR(20) NOT NULL
        CHECK (actor_type IN ('admin', 'customer', 'system', 'paystack')),
    actor_id VARCHAR(100),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100),
    request_id VARCHAR(100),
    idempotency_key VARCHAR(255),
    ip_address VARCHAR(64),
    before_state JSONB,
    after_state JSONB,
    metadata JSONB
);

CREATE INDEX IF NOT EXISTS idx_audit_log_entity ON audit_log (entity_type, entity_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_occurred_at ON audit_log (occurred_at DESC);

-- Request counts per bucket (e.g. a phone number on a route group) per
-- fixed window, shared by every API instance and kept across restarts.
CREATE TABLE IF NOT EXISTS rate_limit_counters (
    bucket VARCHAR(255) NOT NULL,
    window_start TIMESTAMPTZ NOT NULL,
    request_count INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (bucket, window_start)
);

CREATE INDEX IF NOT EXISTS idx_rate_limit_counters_window_start ON rate_limit_counters (window_start);

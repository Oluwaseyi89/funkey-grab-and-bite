-- Background jobs with bounded retries, without an external queue. Jobs are
-- inserted in the same transaction as the change that needs them (an order,
-- a payment), so a job exists only if that change was saved. A worker in each
-- API instance claims due jobs with FOR UPDATE SKIP LOCKED, so no two
-- instances run the same job.
CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    kind VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'running', 'succeeded', 'dead')),
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5 CHECK (max_attempts > 0),
    run_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- A running job whose lock expired was abandoned (the instance stopped)
    -- and may be claimed again.
    locked_until TIMESTAMPTZ,
    last_error TEXT,
    -- Stops the same job being queued twice, e.g. one expiry check per order.
    dedupe_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_jobs_due ON jobs (run_at) WHERE status IN ('pending', 'running');
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs (status, updated_at DESC);

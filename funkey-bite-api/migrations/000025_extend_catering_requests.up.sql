-- Catering requests are either for a package or custom. These columns hold
-- what is needed to price and deliver the event: where it is, how long, how
-- it is served, the food chosen, dietary needs, and the admin's quote.
ALTER TABLE catering_requests
    ADD COLUMN IF NOT EXISTS reference VARCHAR(30),
    ADD COLUMN IF NOT EXISTS request_type VARCHAR(10) NOT NULL DEFAULT 'custom',
    ADD COLUMN IF NOT EXISTS package_id INTEGER REFERENCES catering_packages(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS organization VARCHAR(200),
    ADD COLUMN IF NOT EXISTS event_type_other VARCHAR(100),
    ADD COLUMN IF NOT EXISTS duration_hours DECIMAL(4,1),
    ADD COLUMN IF NOT EXISTS venue_name VARCHAR(200),
    ADD COLUMN IF NOT EXISTS venue_address TEXT,
    ADD COLUMN IF NOT EXISTS venue_area VARCHAR(100),
    ADD COLUMN IF NOT EXISTS venue_state VARCHAR(100),
    ADD COLUMN IF NOT EXISTS venue_notes TEXT,
    ADD COLUMN IF NOT EXISTS service_style VARCHAR(20),
    ADD COLUMN IF NOT EXISTS selected_items JSONB NOT NULL DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS menu_notes TEXT,
    ADD COLUMN IF NOT EXISTS dietary_requirements JSONB NOT NULL DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS dietary_notes TEXT,
    ADD COLUMN IF NOT EXISTS estimated_total DECIMAL(12,2),
    ADD COLUMN IF NOT EXISTS quoted_amount DECIMAL(12,2),
    ADD COLUMN IF NOT EXISTS admin_notes TEXT,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Budgets for large events can exceed DECIMAL(10,2).
ALTER TABLE catering_requests ALTER COLUMN budget TYPE DECIMAL(12,2);

-- Earlier requests get a reference so customers and admins can quote one.
UPDATE catering_requests
SET reference = 'CAT-' || lpad(id::text, 6, '0')
WHERE reference IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_catering_requests_reference ON catering_requests(reference);
CREATE INDEX IF NOT EXISTS idx_catering_requests_status ON catering_requests(status);

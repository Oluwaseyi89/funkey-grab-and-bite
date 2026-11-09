-- Jobs whose payload holds something private, such as an email carrying a
-- customer's one-time quote link. Their payload is wiped once they succeed,
-- and job payloads are never shown in the admin.
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS sensitive BOOLEAN NOT NULL DEFAULT false;

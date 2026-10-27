DROP INDEX IF EXISTS idx_catering_requests_status;
DROP INDEX IF EXISTS idx_catering_requests_reference;

ALTER TABLE catering_requests
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS admin_notes,
    DROP COLUMN IF EXISTS quoted_amount,
    DROP COLUMN IF EXISTS estimated_total,
    DROP COLUMN IF EXISTS dietary_notes,
    DROP COLUMN IF EXISTS dietary_requirements,
    DROP COLUMN IF EXISTS menu_notes,
    DROP COLUMN IF EXISTS selected_items,
    DROP COLUMN IF EXISTS service_style,
    DROP COLUMN IF EXISTS venue_notes,
    DROP COLUMN IF EXISTS venue_state,
    DROP COLUMN IF EXISTS venue_area,
    DROP COLUMN IF EXISTS venue_address,
    DROP COLUMN IF EXISTS venue_name,
    DROP COLUMN IF EXISTS duration_hours,
    DROP COLUMN IF EXISTS event_type_other,
    DROP COLUMN IF EXISTS organization,
    DROP COLUMN IF EXISTS package_id,
    DROP COLUMN IF EXISTS request_type,
    DROP COLUMN IF EXISTS reference;

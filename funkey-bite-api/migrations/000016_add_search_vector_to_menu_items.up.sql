ALTER TABLE menu_items
    ADD COLUMN IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(tags::text, '')), 'C')
    ) STORED;

CREATE INDEX IF NOT EXISTS idx_menu_items_search
    ON menu_items USING GIN(search_vector);

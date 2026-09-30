DROP INDEX IF EXISTS idx_menu_items_search;

ALTER TABLE menu_items
    DROP COLUMN IF EXISTS search_vector;

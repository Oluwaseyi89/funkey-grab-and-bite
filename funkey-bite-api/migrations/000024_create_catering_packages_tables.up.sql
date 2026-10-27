-- Catering packages are managed by the admin. A package sets a per-person
-- price, the guest range it serves, how it is served and how much notice it
-- needs; customers pick dish_choices dishes from the menu items linked to it.
CREATE TABLE IF NOT EXISTS catering_packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price_per_person DECIMAL(12,2) NOT NULL CHECK (price_per_person >= 0),
    min_guests INTEGER NOT NULL CHECK (min_guests >= 1),
    max_guests INTEGER CHECK (max_guests IS NULL OR max_guests >= min_guests),
    service_style VARCHAR(20) NOT NULL CHECK (service_style IN ('drop_off', 'buffet', 'full_service')),
    min_notice_days INTEGER NOT NULL DEFAULT 3 CHECK (min_notice_days >= 0),
    inclusions JSONB NOT NULL DEFAULT '[]',
    dish_choices INTEGER NOT NULL DEFAULT 0 CHECK (dish_choices >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- The menu items a customer may choose from for a package.
CREATE TABLE IF NOT EXISTS catering_package_items (
    package_id INTEGER NOT NULL REFERENCES catering_packages(id) ON DELETE CASCADE,
    menu_item_id INTEGER NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    PRIMARY KEY (package_id, menu_item_id)
);

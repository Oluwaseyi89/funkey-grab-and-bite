CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES menu_categories(id),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    image_url VARCHAR(500),
    is_available BOOLEAN DEFAULT true,
    is_pre_order BOOLEAN DEFAULT false,
    preparation_time INTEGER,
    tags JSONB DEFAULT '[]'::jsonb,
    nutritional_info JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

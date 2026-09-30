CREATE TABLE IF NOT EXISTS inventory_history (
    id SERIAL PRIMARY KEY,
    inventory_item_id INTEGER REFERENCES inventory_items(id),
    previous_stock INTEGER NOT NULL,
    new_stock INTEGER NOT NULL,
    change INTEGER NOT NULL,
    operation VARCHAR(20) NOT NULL,
    reason VARCHAR(200) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

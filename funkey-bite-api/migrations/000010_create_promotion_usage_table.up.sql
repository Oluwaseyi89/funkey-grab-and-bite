CREATE TABLE IF NOT EXISTS promotion_usage (
    id SERIAL PRIMARY KEY,
    promotion_id INTEGER REFERENCES promotions(id),
    order_id INTEGER REFERENCES orders(id),
    customer_id INTEGER REFERENCES users(id),
    discount_applied DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

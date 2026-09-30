CREATE TABLE IF NOT EXISTS business_settings (
    id SERIAL PRIMARY KEY,
    business_name VARCHAR(200) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(200) NOT NULL,
    address TEXT NOT NULL,
    opening_hours TEXT NOT NULL,
    delivery_fee DECIMAL(10,2) DEFAULT 2.99,
    min_order_amount DECIMAL(10,2) DEFAULT 10.00,
    tax_rate DECIMAL(5,2) DEFAULT 8.5,
    is_delivery_open BOOLEAN DEFAULT true,
    is_pickup_open BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

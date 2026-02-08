-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(200) UNIQUE,
    full_name VARCHAR(200) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order number sequence for pattern-based IDs
CREATE SEQUENCE order_number_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 999999
    CYCLE;

-- Update orders table to link to users
ALTER TABLE orders 
ADD COLUMN user_id INTEGER REFERENCES users(id),
ADD COLUMN customer_id VARCHAR(50) GENERATED ALWAYS AS (
    'CUST-' || LPAD(id::text, 6, '0')
) STORED;

-- Create index for performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_customer_phone ON orders(customer_phone);
CREATE INDEX idx_users_phone_email ON users(phone, email);
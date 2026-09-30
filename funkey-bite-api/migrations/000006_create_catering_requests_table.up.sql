CREATE TABLE IF NOT EXISTS catering_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    event_name VARCHAR(200),
    contact_name VARCHAR(200) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    contact_email VARCHAR(200),
    event_date DATE NOT NULL,
    event_time TIME,
    guest_count INTEGER,
    event_type VARCHAR(100),
    budget DECIMAL(10,2),
    special_requests TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    seller_id INT REFERENCES sellers(id),
    message VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS favorites (
    PRIMARY KEY (user_id, product_id),
    user_id INT REFERENCES users(id),
    product_id INT REFERENCES products(id),
    created_at TIMESTAMP DEFAULT NOW()
);
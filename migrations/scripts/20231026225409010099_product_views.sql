CREATE TABLE product_views (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    user_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);
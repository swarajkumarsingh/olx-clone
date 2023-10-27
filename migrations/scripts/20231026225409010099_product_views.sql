CREATE TABLE product_views (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    user_id INT REFERENCES users(id),
    view_time TIMESTAMP DEFAULT NOW()
);
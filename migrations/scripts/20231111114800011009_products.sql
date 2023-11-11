CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	location TEXT NOT NULL,
	coordinates TEXT NOT NULL,
    views BIGINT DEFAULT 0,
    price VARCHAR(100) NOT NULL,
    -- category INT REFERENCES category(id),
    seller_id INT REFERENCES sellers(id),
    created_at TIMESTAMP DEFAULT NOW()
);
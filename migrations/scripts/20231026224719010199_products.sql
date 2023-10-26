DROP TABLE products;

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    images TEXT[],
	location TEXT,
	coordinates TEXT,
    views BIGINT DEFAULT 0,
    price VARCHAR(100) NOT NULL,
    -- category INT REFERENCES category(id),
    seller_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);
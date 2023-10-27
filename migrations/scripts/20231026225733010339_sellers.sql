CREATE TABLE sellers (
	id SERIAL NOT NULL PRIMARY KEY,
	username UNIQUE VARCHAR(500) NOT NULL,
	fullname VARCHAR(100) NOT NULL,
	email UNIQUE VARCHAR(100) NOT NULL,
	password VARCHAR(200) NOT NULL,
	number UNIQUE VARCHAR(10) NOT NULL,
    rating NUMERIC(3, 2) CHECK (rating >= 0 AND rating <= 5),
    description TEXT,
	avatar TEXT,
	location TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
	coordinates TEXT,
	created_at TIMESTAMP DEFAULT NOW()
);
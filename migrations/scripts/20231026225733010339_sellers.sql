CREATE TABLE IF NOT EXISTS sellers (
	id SERIAL NOT NULL PRIMARY KEY,

	avatar TEXT,
    description TEXT,
	password VARCHAR(200) NOT NULL,
	fullname VARCHAR(100) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
	email VARCHAR(100) UNIQUE NOT NULL,
	username VARCHAR(500) UNIQUE NOT NULL,

	location TEXT,
	coordinates TEXT,
	phone VARCHAR(10) NOT NULL,

	city VARCHAR(50) NOT NULL,
	state VARCHAR(50) NOT NULL,
	country VARCHAR(50) NOT NULL
	zip_code VARCHAR(10) NOT NULL,

	created_at TIMESTAMP DEFAULT NOW()
    rating NUMERIC(3, 2) CHECK (rating >= 0 AND rating <= 5),
	account_status ENUM('Active', 'Suspended', 'Banned') NOT NULL,
);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL NOT NULL PRIMARY KEY,
	username VARCHAR(500) UNIQUE NOT NULL,
	fullname VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL,
	password VARCHAR(200) NOT NULL,
	phone VARCHAR(12) NOT NULL,
	avatar TEXT,
	location TEXT,
	coordinates TEXT,
	created_at TIMESTAMP DEFAULT NOW(),
	otp TEXT,
	otp_expiration TIMESTAMP
);
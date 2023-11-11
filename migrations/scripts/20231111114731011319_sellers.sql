CREATE TABLE IF NOT EXISTS sellers (
	id SERIAL PRIMARY KEY,

    description TEXT,
	password VARCHAR(200) NOT NULL,
	fullname VARCHAR(100) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
	email VARCHAR(100) UNIQUE NOT NULL,
	username VARCHAR(500) UNIQUE NOT NULL,
	avatar TEXT DEFAULT 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTOc9VDs02ZrmIC7pS3WzBTvXl8UrI3jwAOVQ&usqp=CAU',

	location TEXT,
	coordinates TEXT,
	phone VARCHAR(10) NOT NULL,

	city VARCHAR(50),
	state VARCHAR(50),
	country VARCHAR(50),
	zip_code VARCHAR(10),

	created_at TIMESTAMP DEFAULT NOW(),
    rating TEXT DEFAULT '0',
	account_status TEXT DEFAULT 'active'
);

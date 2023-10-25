		CREATE TABLE IF NOT EXISTS product (
			id SERIAL NOT NULL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(200) NOT NULL,
			number VARCHAR(10) NOT NULL,
			avatar TEXT,
			address TEXT,
			created_at DATE
		);
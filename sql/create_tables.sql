CREATE TABLE IF NOT EXISTS users (
	id SERIAL NOT NULL,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,
	username VARCHAR(50),
	email VARCHAR(50),
	password VARCHAR(100),
	firstname VARCHAR(50)
) 
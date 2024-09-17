CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	user_identifier varchar(50) NOT NULL UNIQUE,
	first_name varchar(50) NOT NULL,
	last_name varchar(50) NULL DEFAULT '',
	email varchar(100) NOT NULL UNIQUE ,
	avatar_url varchar(500)  NULL,
	description text  NULL,
	password varchar(255) NOT NULL,
	is_verify bool DEFAULT false,
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
    updated_date timestamptz(0) NULL,
    deleted_date timestamptz(0) NULL
);


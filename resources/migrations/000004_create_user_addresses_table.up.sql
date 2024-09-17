CREATE TABLE user_addresses (
	id SERIAL PRIMARY KEY,
	user_id bigint references users (id) ON DELETE CASCADE ON UPDATE CASCADE, 
    address_identifier varchar(50) NOT NULL UNIQUE,
	title varchar(100) NULL,
	full_address text  NULL,
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
    updated_date timestamptz(0) NULL,
    deleted_date timestamptz(0) NULL
);
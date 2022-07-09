CREATE TABLE customer (
	id VARCHAR(36) PRIMARY KEY,
	version INT NOT NULL,
	first_name VARCHAR(50),
	last_name VARCHAR(50),
	email VARCHAR(255)
);

CREATE UNIQUE INDEX customer_email_uk ON customer (email);

INSERT INTO customer (id, version, first_name, last_name, email) values('b4f990a2-707b-41ae-aa19-9224fb8d2d7e', 1, 'Paulo', 'Pereira', 'paulo.pereira@mail.com');

CREATE TABLE registration (
	id VARCHAR(36) PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	verified BOOLEAN
);

CREATE UNIQUE INDEX registration_email_uk ON registration (email);

CREATE TABLE outbox (
	id BIGSERIAL PRIMARY KEY,
	kind VARCHAR(50) NOT NULL,
	payload BYTEA,
	consumed BOOLEAN
);

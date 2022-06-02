CREATE TABLE customer (
	id VARCHAR(36),
	version INT,
	first_name VARCHAR(50),
	last_name VARCHAR(50),
	email VARCHAR(255),
	PRIMARY KEY(id)
);

CREATE UNIQUE INDEX email_uk ON customer (email);

INSERT INTO customer (id, version, first_name, last_name, email) values('b4f990a2-707b-41ae-aa19-9224fb8d2d7e', 1, 'Paulo', 'Pereira', 'paulo.pereira@mail.com');

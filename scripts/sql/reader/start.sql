DROP DATABASE IF EXISTS reader;
CREATE DATABASE reader;

\connect reader

CREATE TABLE table_with_pk
(
	a SERIAL,
	b VARCHAR(30),
	c TIMESTAMP NOT NULL,
	PRIMARY KEY(a, c)
);

CREATE TABLE table_without_pk
(
	a SERIAL,
	b NUMERIC(5,2),
	c TEXT
);

SELECT 'init' FROM pg_create_logical_replication_slot('td',  'test_decoding');

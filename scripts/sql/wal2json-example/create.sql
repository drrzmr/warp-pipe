DROP DATABASE IF EXISTS wal2json;
CREATE DATABASE wal2json;

\connect wal2json

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

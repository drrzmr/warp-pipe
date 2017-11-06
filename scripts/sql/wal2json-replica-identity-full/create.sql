DROP DATABASE IF EXISTS wal2json_ri_full;
CREATE DATABASE wal2json_ri_full;

\connect wal2json_ri_full

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

ALTER TABLE table_with_pk REPLICA IDENTITY FULL;
ALTER TABLE table_without_pk REPLICA IDENTITY FULL;

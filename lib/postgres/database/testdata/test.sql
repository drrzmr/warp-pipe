CREATE TABLE test
(
  a SERIAL,
  b VARCHAR(30),
  c TIMESTAMP NOT NULL,
  PRIMARY KEY (a)
);

BEGIN;
INSERT INTO test (b, c) VALUES ('test 1', now());
INSERT INTO test (b, c) VALUES ('test 2', now());
INSERT INTO test (b, c) VALUES ('test 3', now());
COMMIT;
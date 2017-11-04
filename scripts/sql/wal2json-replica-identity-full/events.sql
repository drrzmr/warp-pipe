\connect wal2json_ri_full

BEGIN;
INSERT INTO table_with_pk(b, c) VALUES('Backup and Restore', now());
INSERT INTO table_with_pk (b, c) VALUES('Tuning', now());
INSERT INTO table_with_pk (b, c) VALUES('Replication', now());
DELETE FROM table_with_pk WHERE a < 3;

INSERT INTO table_without_pk (b, c) VALUES(2.34, 'Tapir');
UPDATE table_without_pk SET c = 'Anta' WHERE c = 'Tapir';
COMMIT;

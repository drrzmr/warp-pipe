\connect wal2json_ri_full

SELECT * FROM pg_logical_slot_get_changes('w2j', NULL, NULL, 'pretty-print', '1', 'write-in-chunks', '0');
SELECT * FROM pg_logical_slot_get_changes('td', NULL, NULL);

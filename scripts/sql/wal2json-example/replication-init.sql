\connect wal2json

SELECT 'init' FROM pg_create_logical_replication_slot('td',  'test_decoding');
SELECT 'init' FROM pg_create_logical_replication_slot('w2j', 'wal2json');

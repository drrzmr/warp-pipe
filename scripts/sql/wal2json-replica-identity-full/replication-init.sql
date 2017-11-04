\connect wal2json_ri_full

SELECT 'init' FROM pg_create_logical_replication_slot('td',  'test_decoding');
SELECT 'init' FROM pg_create_logical_replication_slot('w2j', 'wal2json');

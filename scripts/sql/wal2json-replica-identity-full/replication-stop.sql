\connect wal2json_ri_full

SELECT 'stop' FROM pg_drop_replication_slot('td');
SELECT 'stop' FROM pg_drop_replication_slot('w2j');

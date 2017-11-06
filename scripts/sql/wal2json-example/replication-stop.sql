\connect wal2json

SELECT 'stop' FROM pg_drop_replication_slot('td');
SELECT 'stop' FROM pg_drop_replication_slot('w2j');

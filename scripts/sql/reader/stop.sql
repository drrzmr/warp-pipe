\connect reader

SELECT 'stop' FROM pg_drop_replication_slot('td');

\connect postgres
DROP DATABASE IF EXISTS reader;

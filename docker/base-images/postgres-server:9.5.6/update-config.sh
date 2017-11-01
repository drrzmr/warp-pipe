#!/bin/bash -eux

###################
# postgresql.conf #
###################

# enable conf.d directory
mkdir -p ${PGDATA}/conf.d
sed -i "s/^#include_dir/include_dir/" ${PGDATA}/postgresql.conf

# enable wal replication slots
cat << EOF > ${PGDATA}/conf.d/wal.conf
max_wal_senders = 2
max_replication_slots = 2
wal_level = logical
EOF

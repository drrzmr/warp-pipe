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


###############
# pg_hba.conf #
###############

# open replication protocol, except for ipv6
cat << EOF >> ${PGDATA}/pg_hba.conf

###########################
# open replication access #
###########################
local replication postgres           trust
host  replication postgres 0.0.0.0/0 trust
EOF

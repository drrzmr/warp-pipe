FROM postgres:9.5.6

ENV WAL2JSON_DEPENDENCIES "git build-essential postgresql-server-dev-9.5 ca-certificates"

RUN set -eux; \
	DEBIAN_FRONTEND=noninterative; \
	USE_PGXS=1; \
	\
	apt-get -q update; \
	apt-get -qy install --no-upgrade --no-install-recommends ${WAL2JSON_DEPENDENCIES}; \
	\
	cd /tmp; \
	git clone --depth 1 https://github.com/eulerto/wal2json.git; \
	cd wal2json; \
	make; \
	make install; \
	\
	apt-get -qy remove --purge ${WAL2JSON_DEPENDENCIES}; \
	apt-get -qy autoremove --purge; \
	apt-get -qy clean; \
	rm -rf \
		/var/lib/apt/lists/* \
		/tmp/*

COPY docker/base-images/postgres-server:9.5.6/update-config.sh /docker-entrypoint-initdb.d/

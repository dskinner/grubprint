#! /usr/bin/env sh

if [ ! -d "/data/postgres/" ]; then
	mkdir -p /data/postgres
	chown postgres /data/postgres
	su postgres -c "/usr/lib/postgresql/9.3/bin/pg_ctl initdb -D /data/postgres -o '--locale=en_US.utf8 -E UTF8'"
	echo host all all 0.0.0.0 0.0.0.0 md5 >> /data/postgres/pg_hba.conf
	echo "listen_addresses='*'" >> /data/postgres/postgresql.conf
	su postgres -c "/usr/lib/postgresql/9.3/bin/pg_ctl start -D /data/postgres -c -w -l /dev/null"
	su -c "psql -c \"ALTER USER postgres with encrypted password 'postgres';\" template1" postgres
	su -c "psql -c \"CREATE DATABASE food;\" template1" postgres
	su postgres -c "/usr/lib/postgresql/9.3/bin/pg_ctl stop -D /data/postgres -c -w -l /dev/null"
fi

su postgres -c "/usr/lib/postgresql/9.3/bin/pg_ctl start -D /data/postgres -c -w"
/usr/sbin/sshd -D

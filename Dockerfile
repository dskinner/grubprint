FROM ubuntu

MAINTAINER Daniel Skinner <daniel@dasa.cc>

RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get -y update
RUN apt-get -y upgrade

RUN locale-gen en_US.UTF-8

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install vim.tiny openssh-server postgresql
RUN mkdir -p /var/run/sshd
RUN echo "root:root" | chpasswd

RUN mkdir /data
RUN chown postgres /data
RUN su postgres -c "/usr/lib/postgresql/9.1/bin/pg_ctl initdb -D /data -o '--locale=en_US.utf8 -E UTF8'"
RUN echo host all all 0.0.0.0 0.0.0.0 md5 >> /data/pg_hba.conf
RUN echo "listen_addresses='*'" >> /data/postgresql.conf
RUN su postgres -c "/usr/lib/postgresql/9.1/bin/pg_ctl start -D /data -c -w -l /dev/null" && su -c "psql -c \"ALTER USER postgres with encrypted password 'postgres';\" template1" postgres && su -c "psql -c \"CREATE DATABASE food;\" template1" postgres

CMD su postgres -c "/usr/lib/postgresql/9.1/bin/pg_ctl start -D /data -c -w" && /usr/sbin/sshd -D

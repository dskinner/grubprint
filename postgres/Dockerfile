FROM postgres:9.4

ENV POSTGRES_PASSWORD=mysecretpassword 
ENV PGDATA /data/postgres

RUN mkdir -p /data
ADD usda/ /data/usda/
ADD schema.sql /data/
ADD import.sql /data/

ADD init.sh /docker-entrypoint-initdb.d/

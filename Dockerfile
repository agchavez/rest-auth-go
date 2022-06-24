FROM postgres:9.6


COPY ./migrations/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
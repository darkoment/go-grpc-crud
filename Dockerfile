# Create MySQL Image for GO-GRPC-CRUD-API
FROM mysql:8
MAINTAINER darkoment@yandex.ru

ENV MYSQL_ROOT_PASSWORD=TestOnGo 

COPY ./sql.db/dump.sql /docker-entrypoint-initdb.d

#ADD dump.sql /docker-entrypoint-initdb.d

EXPOSE 3306
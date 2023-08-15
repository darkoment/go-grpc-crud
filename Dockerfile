# Create MySQL Image for GO-GRPC-CRUD-API
FROM mysql
MAINTAINER darkoment@yandex.ru

ENV MYSQL_ROOT_PASSWORD TestOnGo 

ADD dump.sql /docker-entrypoint-initdb.d

EXPOSE 3306
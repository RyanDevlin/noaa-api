#!/usr/bin/env bash

# NOTE: Paths are relative for container volumes so the startup script must stay in the correct directory (at least until I decide to make this fancy). Both volume directories should be in the gitignore and therefore won't end up online. If you want to run this locally, create these directories on your machine and the container will populate them.

# Standup mysql DB (obviously change this to store the password in a file. Also use a better password :) )
podman run --rm --name wordpress-sqldb -e MYSQL_ROOT_PASSWORD=yeet --net slirp4netns:allow_host_loopback=true -p 3306:3306 -v ./wordpress-db:/var/lib/mysql -d mysql:8.0

# Standup the wordpress frontend
podman run --rm --name wordpress-frontend -e WORDPRESS_DB_HOST=10.0.2.2:3306 -e WORDPRESS_DB_USER=wordpress -e WORDPRESS_DB_PASSWORD=yeet -e WORDPRESS_DB_NAME=wordpress -e WORDPRESS_TABLE_PREFIX=wp_ --net slirp4netns:allow_host_loopback=true -p 8080:80 -v ./wp-content:/var/www/html/wp-content -d wordpress:5.8

# To log into the mysql instance, use this command:

# mysql -h localhost -P 3306 --protocol=tcp -u root -p 


# Mysql commands to configure the DB

# CREATE DATABASE wordpress;
# CREATE USER 'wordpress'@'10.0.2.100' IDENTIFIED BY 'yeet';
# GRANT ALL PRIVILEGES ON wordpress.* TO "wordpress"@"10.0.2.100";
# FLUSH PRIVILEGES;
# EXIT

# Note: 10.0.2.100 is some weird address podman is using to reference localhost. I cannot replace this address with 'localhost' or '127.0.0.1'.

# Note: 10.0.2.2 is the address podman uses to reference the host machine when --net slirp4netns:allow_host_loopback=true is used. This is used to tell the wordpress container to connect to the database on the local host at default port 3306.

# Because the database for the mysql container is mounted as a local volume, the data and configurations persist after container restart.

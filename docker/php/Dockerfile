FROM php:8.3-apache

RUN apt-get update && apt-get install -y libpq-dev \
    && docker-php-ext-install pdo pdo_pgsql \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY public/ /var/www/html/

COPY docker/php/config.php /var/www/html/config.php

EXPOSE 80

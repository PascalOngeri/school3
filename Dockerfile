# Use the official PHP image with FPM (FastCGI Process Manager)
FROM php:7.4-fpm

# Set the working directory inside the container
WORKDIR /var/www/html

# Copy source code from local machine to the container
COPY . /var/www/html

# Copy the .env config file into the container
COPY configs/.env /var/www/html/.env

# Install PHP extensions if needed
RUN docker-php-ext-install pdo pdo_mysql

# Set permissions
RUN chmod -R 755 /var/www/html

# Expose port 9000 for PHP-FPM
EXPOSE 9000

# Start PHP-FPM when the container runs
CMD ["php-fpm"]
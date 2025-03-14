#!/bin/sh
until mysqladmin ping -h mysql_clean-arch -u root -proot --protocol=tcp --silent; do
  echo "Waiting for MySQL..."
  sleep 2
done
echo "MySQL is ready!"
exec "$@"
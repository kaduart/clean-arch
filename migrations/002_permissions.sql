CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY 'root'; -- Cria o usuário
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION; -- Concede privilégios
FLUSH PRIVILEGES;
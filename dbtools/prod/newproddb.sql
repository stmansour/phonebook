DROP DATABASE IF EXISTS accord; CREATE DATABASE accord; USE accord;
source accorddb.sql
GRANT ALL PRIVILEGES ON accord TO 'ec2-user'@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON accord TO 'adbuser'@'localhost' WITH GRANT OPTION;

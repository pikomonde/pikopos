sudo apt-get update
sudo apt-get install mysql-server

## Init MySQL
`mysql -h localhost -u root -p < ./setup/deploy_00.00.001_init_schema.sql`

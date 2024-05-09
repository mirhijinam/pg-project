# PGStart Trainee

Приложение, предоставляющее REST API для запуска команд - bash-скриптов. 
Приложение позволяет параллельно запускать произвольное количество команд.

## 🚀 Инструкция по запуску
Для запуска приложения необходимо:
1. ```git clone git@github.com:mirhijinam/pg-project.git```
2. ```cd pg-start```
3. ```touch .env```
4. ```bash
   cat <<EOF > .env
   MAXCOUNT=10

   PGUSER=pguser
   PGPASSWORD=pgpass
   PGHOST=pghost
   PGPORT=5432
   PGDATABASE=pgdb
   PGSSLMODE=disable

   HTTP_PORT=7070
   SERVER_ENDPOINT=localhost:7070/
   EOF
   ```

5. ```go mod download```
6. ```make up```

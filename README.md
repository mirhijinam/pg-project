# PGStart Trainee

Приложение, предоставляющее REST API для запуска команд - bash-скриптов. 
Приложение позволяет параллельно запускать произвольное количество команд.

## 🚀 Инструкция по запуску
Для запуска приложения необходимо:
1. Клонирование репозитория.
   
   ```git clone git@github.com:mirhijinam/pg-project.git```
   
2. Переход в директорию проекта.
   
   ```cd pg-start```
   
3. Создание файла окружения с необходимыми переменными.
   
   ```touch .env```

     Можете использовать готовый пример: 
   ```bash
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
   TIMEOUT="5s"
   IDLE_TIMEOUT="30s"

   ADMIN_TOKEN="admin"

   LOGENV="prod"
   EOF
   ```
   
4. Загрузка зависимостей.
   
   ```go mod download```
   
5. Запуск приложения.
   
    ```make up```


## 📚 Дополнительно
Были реализованы:

   - Использование токена в хедерах запросов для создания sudo команд
   - Долгие команды с сохранением вывода в БД по мере выполнения
   - POST-метод для остановки команд
   - Миграция
   - Логгирование методов через Middleware
   - Документирование с помощью Swagger и формата ADR


## 📬 Примеры запросов и ответов
### Остановка команды:
Ее условием является то, чтобы останавливаемая команда еще либо не успела записать логи из "короткой" команды, либо не успела дописать логи "долгой". Пример:
   ```
   curl -X POST http://localhost:7070/stop_cmd/3
   ```
   *Ответ:*
   ```
   {
           "id of stopped command": 3
   }
   ```
   Таким образом, обращаясь к записям остановленной команды в БД, получаем:
   ```
   {
           "requested command": {
                   "id": 3,
                   "name": "echo",
                   "raw": "echo z1; sleep 2; echo z2; sleep 2; echo z3; sleep 2; echo z4; sleep 2; echo z5; sleep 5",
                   "status": "stopped",
                   "error_msg": "signal: killed",
                   "logs": "z1\nz2\n",
                   "created_at": "2024-05-13T04:19:50.404658Z",
                   "updated_at": "2024-05-13T04:19:54.424816Z"
           }
   }
   ```
### Остальные примеры запросов с более простой логикой перечислены ниже:

- Запуск sudo-команды с токеном админа:
   ```
   curl -X POST "http://localhost:7070/create_cmd" \
        -H "Content-Type: application/json" \
        -H "token: admin" \
        -d '{
              "cmd_raw": "sudo echo z1; sleep 1; echo z2; sleep 1; echo z3; sleep 1; echo z4; sleep 1; echo z5; sleep 1",
              "is_long_cmd": false
            }'
   ```
   *Ответ:*
   ```
   {
           "created command": {
                   "id": 1,
                   "name": "sudo echo",
                   "raw": "sudo echo z1; sleep 1; echo z2; sleep 1; echo z3; sleep 1; echo z4; sleep 1; echo z5; sleep 1",
                   "created_at": "2024-05-13T04:14:42.54506Z"
           }
   }
   ```
- Запуск sudo-команды без токена админа:
   ```
   curl -X POST "http://localhost:7070/create_cmd" \
        -H "Content-Type: application/json" \
        -d '{
              "cmd_raw": "sudo echo z1; sleep 1; echo z2; sleep 1; echo z3; sleep 1; echo z4; sleep 1; echo z5; sleep 1",
              "is_long_cmd": false
            }'
   ```
   *Ответ:*
   ```
   {
           "answer": "error! user has no access"
   }
   
   ```

- Запуск "долгой" команды:
   ```
   curl -X POST "http://localhost:7070/create_cmd" \
        -H "Content-Type: application/json" \
        -d '{
              "cmd_raw": "echo z1; sleep 2; echo z2; sleep 2; echo z3; sleep 2; echo z4; sleep 2; echo z5; sleep 5",
              "is_long_cmd": true
            }'
   ```
   *Ответ:*
   ```
   {
           "created command": {
                   "id": 2,
                   "name": "echo",
                   "raw": "echo z1; sleep 2; echo z2; sleep 2; echo z3; sleep 2; echo z4; sleep 2; echo z5; sleep 5",
                   "created_at": "2024-05-13T04:15:05.736278Z"
           }
   }
   ```
- Получение списка запущенных команд
   ```
   curl http://localhost:7070/cmd_list/
   ```
   *Ответ:*
   ```
   {
           "list of executed commands": [
                   {
                           "id": 1,
                           "name": "sudo echo",
                           "raw": "sudo echo z1; sleep 1; echo z2; sleep 1; echo z3; sleep 1; echo z4; sleep 1; echo z5; sleep 1",
                           "status": "success",
                           "logs": "z2\nz3\nz4\nz5\n",
                           "created_at": "2024-05-13T04:14:42.54506Z",
                           "updated_at": "2024-05-13T04:14:47.601284Z"
                   },
                   {
                           "id": 2,
                           "name": "echo",
                           "raw": "echo z1; sleep 2; echo z2; sleep 2; echo z3; sleep 2; echo z4; sleep 2; echo z5; sleep 5",
                           "status": "success",
                           "logs": "z1\nz2\nz3\nz4\nz5\n",
                           "created_at": "2024-05-13T04:15:05.736278Z",
                           "updated_at": "2024-05-13T04:15:18.763786Z"
                   }
           ]
   }
   ```

- Получение определенной команды из списка
   ```
   curl -X POST http://localhost:7070/stop_cmd/1
   ```
   *Ответ:*
   ```
   {
           "requested command": {
                   "id": 1,
                   "name": "sudo echo",
                   "raw": "sudo echo z1; sleep 1; echo z2; sleep 1; echo z3; sleep 1; echo z4; sleep 1; echo z5; sleep 1",
                   "status": "success",
                   "logs": "z2\nz3\nz4\nz5\n",
                   "created_at": "2024-05-13T04:14:42.54506Z",
                   "updated_at": "2024-05-13T04:14:47.601284Z"
           }
   }
   ```

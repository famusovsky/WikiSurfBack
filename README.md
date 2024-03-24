### Выполнил Степанов Алексей Александрович

## Запуск:

Запуск dev-среды с помощью docker-compose:

```bash
docker-compose up
```
Также доступен запуск с помощью Dockerfile, в данном случае среда, в которой происходит запуск, должна иметь переменные окружения:
- DB_HOST - хост базы данных.
- DB_PORT - порт для доступа к базе данных.
- DB_USER - логин пользователя БД.
- DB_PASSWORD - пароль пользователя БД.
- DB_NAME - ИМЯ пользователя БД.

Запуск с помощью go run:

```bash
go run ./cmd/web/main.go
# Флаги:
# -override_tables=true - запуск с автоматической перезаписью таблиц в БД
# -addr=:8080 - выбор порта, с которым будет работать сервер
# -dsn="postgres://user:qwerty@localhost:8888/my_db?sslmode=disable" - выбор data source name для подключения к БД.
```

## PostgreSQL Query для создания таблиц в БД вручную:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);	
CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    start TEXT NOT NULL,
    finish TEXT NOT NULL,
    creator_id INTEGER NOT NULL,
    CONSTRAINT start_finish UNIQUE (start, finish),
    FOREIGN KEY (creator_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS sprints (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    length_time INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    route_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    path TEXT ARRAY NOT NULL,
    FOREIGN KEY (route_id) REFERENCES routes(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS tournaments (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    pswd TEXT NOT NULL,
    private BOOLEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS tournament_users (
    tour_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (tour_id, user_id)
);
CREATE TABLE IF NOT EXISTS tournament_creators (
    tour_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (tour_id, user_id)
);
CREATE TABLE IF NOT EXISTS tournament_routes (
    tour_id INTEGER NOT NULL,
    route_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (route_id) REFERENCES routes(id),
    PRIMARY KEY (tour_id, route_id)
);
```

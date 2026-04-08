# FIAS Search Service

Микросервис для поиска адресов по базе ФИАС.

---

## Стек технологий

- **Backend (Processing):** Go (Golang) + pgx  
- **Backend (API):** PHP 8.x  
- **Database:** PostgreSQL 15+ (с расширением `pg_trgm`)  
- **Frontend:** Vanilla JS, CSS3, HTML5  

---

## Установка окружения

Установите:

- Go: <https://go.dev/dl/>  
- PHP: <https://www.php.net/downloads>  
- PostgreSQL: <https://www.postgresql.org/download/>  

## Проверка окружения

```bash
go version
php -v
psql --version
```

### 2. Клонирование и зависимости

```bash
git clone [https://github.com/Drengis/loveru_test](https://github.com/Drengis/loveru_test)
cd loveru_test
go mod download
```

## Настройка БД

Создать базу в Postgres.
Прописать доступы в public/config.php и в load_bd.go и build_addresses.go

## Миграции и индексы

```bash
php backend-php/migrate.php
```

## Импорт данных (Go)

```bash
go run backend-go/load_bd/load_bd.go
go run backend-go/build_addresses/build_addresses.go
```

## Запуск сервера

```bash
php -S localhost:8000 -t public
```

URL: <http://localhost:8000/index.php>

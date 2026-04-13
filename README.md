# FIAS Search Service

Микросервис для поиска адресов по базе ФИАС.

---

## Стек технологий

- **Backend (Processing):** Go (Golang) + pgx
- **Backend (API):** PHP 8.3 + Apache
- **Database:** PostgreSQL 17 (с расширением `pg_trgm`)
- **Frontend:** Vanilla JS, CSS3, HTML5
- **Инфраструктура:** Docker, Docker Compose

---

## Быстрый старт (Docker)

Требуется только [Docker Desktop](https://www.docker.com/products/docker-desktop/).

```bash
git clone https://github.com/Drengis/loveru_test
cd loveru_test
docker compose up --build
```

### Что произойдёт автоматически:

| Шаг | Контейнер | Описание |
|---|---|---|
| 1 | `postgres` | Запуск PostgreSQL |
| 2 | `migrate` | Создание таблиц и индексов |
| 3 | `load-bd` | Загрузка данных ФИАС (~час) |
| 4 | `build-addresses` | Сборка полных адресов |
| — | `web` | PHP сервер поиска (стартует сразу) |

После завершения загрузки данных сервис доступен по адресу:

**http://localhost:8080**

> Загрузка данных ФИАС (`load-bd`) занимает продолжительное время — идёт скачивание и парсинг большого ZIP-архива с сайта nalog.ru. Поисковый интерфейс будет доступен сразу, но результаты появятся только после завершения `load-bd` и `build-addresses`.

---

## Структура проекта

```
loveru_test/
├── backend-go/
│   ├── load_bd/          # Загрузка данных ФИАС в БД
│   └── build_addresses/  # Сборка полных адресных строк
├── backend-php/
│   └── migrate.php       # Создание таблиц и индексов
├── public/               # Фронтенд + PHP API поиска
│   ├── index.html
│   ├── index.php
│   ├── script.js
│   └── config.php
└── docker/
    ├── php/Dockerfile
    └── go/Dockerfile
```

---

## Ручной запуск (без Docker)

Установите: Go, PHP, PostgreSQL.

Пропишите параметры подключения в `public/config.php` и в константах `load_bd.go` / `build_addresses.go`.

```bash
# Миграции
php backend-php/migrate.php

# Загрузка данных ФИАС
go run backend-go/load_bd/load_bd.go

# Сборка адресов
go run backend-go/build_addresses/build_addresses.go

# Сервер
php -S localhost:8000 -t public
```

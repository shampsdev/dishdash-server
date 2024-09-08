# DishDash Server

Стек: Go, Gin, Socket.IO, PostgreSQL, Docker, Traefik (на сервере)

[Development](#Development)

Серверная часть проекта для Creative Space Hackathon \
Занял первое место. [Сайт](https://dishdash.ru)

<img src="https://github.com/shampiniony/dishdash-server/assets/79862574/0d0a7d7b-13d1-4a37-9c26-abb1c844b335" width="300">

Выбирай с друзьями где поесть с помощью свайпов. \
Люди в одном лобби свайпают заведения, пока не произойдет match.

https://github.com/shampiniony/dishdash-server/assets/79862574/6987d0d0-3d09-4c2e-83fb-64256e7fff13

Prod: https://dishdash.ru/api/v1

### Development

Скопировать .env

```
cp .env.example .env
```
В частности в .env нужен `TWOGIS_API_KEY`

Есть два варианта запуска:

#### 1. Full compose
Подойдёт, если не планируете активно изменять код

```
make compose-up
make db-default-data
```

#### 2. Compose db + go run
Подойдёт, если занимаетесь разработкой
```
make db-compose-up
make db-default-data
make run
```

#### Детали

---
- Рекомендуемая версия `go1.22.6`
- `make help` выдаёт много полезной информации
- При `POSTGRES_AUTOMIGRATE=True` бек поднимет миграции при подключении, если это не удастся, он упадёт с ошибкой
- Документация rest по адресу http://localhost:8000/api/v1/swagger/index.html
- `make db-default-data` добавит в базу некий сет тегов. С пустой базой тегов работать не будет. Требуется установленный `psql` (`postgresql-client-16`). Можно и другим образом занести эти данные из [migrations/data/default.sql](migrations/data/default.sql)
- на `1000` порту запускается adminer (см. `docker-compose.yml`). Креды для доступа (в ui adminer) при параметрах из .env.example: `Engine=PostgreSQL Host=database, User=root Password=root Database=root`
- Инструменты для разработки (golangci-lint, swag, migrate) автоматически устанавливаются в `bin/` папку в корне проекта. Нужен установленный go.

#### Документация

[Протокол](https://linear.app/shampiniony/document/protocol-documentation-9118e6048e55)

[База данных](https://linear.app/shampiniony/document/dishdash-database-d785527280fb)
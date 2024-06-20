# DishDash

Стек: Go, Gin, Socket.IO, PostgreSQL, Docker, Traefik (на сервере)

Проект для Creative Space Hackathon \
Занял первое место

<img src="https://github.com/shampiniony/dishdash-server/assets/79862574/0d0a7d7b-13d1-4a37-9c26-abb1c844b335" width="300">

Выбирай с друзьями где поесть с помощью свайпов. \
Люди в одном лобби свайпают заведения, пока не произойдет match.

https://github.com/shampiniony/dishdash-server/assets/79862574/6987d0d0-3d09-4c2e-83fb-64256e7fff13

### Deployment

Скопировать .env

```
cp .env.example .env
```

Есть два варианта запуска:

#### 1. Full compose

```
make compose-up
```

#### 2. Compose db + go run

```
make db-compose-up
make run
```

#### Детали

---

- Бек поднимет миграции при подключении, если это не удастся, он упадёт с ошибкой
- Документация rest по адресу http://localhost:8000/api/v1/swagger/index.html
- Документация socket.io далее

### Socket.io

---

Подключение к socket.io: http://localhost:8000/socket.io

### Events !!В РАЗРАБОТКЕ!!

---

__joinLobby__ : client->server

Первый запрос при подключении.

__request__:

```
{
    "id": "AV1S2",
}
```

__response__:

```
{
    "status": "connected",
    "lobby": TODO
}
```

---

__changeSettings__: client->server

```
{
    "settings"
}
```

__response__:

```
{
    "status": "connected",
}
```
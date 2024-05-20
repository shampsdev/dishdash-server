# DishDash
Стек: Go, Gin, Socket.IO, PostgreSQL, Docker, Traefik (на сервере)

Проект для Creative Space Hackathon \
Занял первое место

<img src="https://github.com/shampiniony/dishdash-server/assets/79862574/0d0a7d7b-13d1-4a37-9c26-abb1c844b335" width="300">

Выбирай с друзьями где поесть с помощью свайпов. \
Люди в одном лобби свайпают заведения, пока не произойдет match.


https://github.com/shampiniony/dishdash-server/assets/79862574/6987d0d0-3d09-4c2e-83fb-64256e7fff13



### Deployment
Create .env file
```
cp .env.example .env
```

Run database and adminer (localhost:1000)
```
docker-compose up --build -d
```

Run
```
go run cmd/server/main.go
```

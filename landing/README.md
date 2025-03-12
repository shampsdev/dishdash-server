## Landing DishDash

По умолчанию при билде с помощью `Dockerfile` фронтенд запускается на `:80` порту.

Пример `docker-compose.yml`:
```
version: '3.8'

services:
  frontend:
    ports:
      - ${NGINX_PORT-3000}:80
    build:
      context: .
      dockerfile: Dockerfile
```
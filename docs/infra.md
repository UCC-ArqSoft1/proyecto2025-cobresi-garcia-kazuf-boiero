# Infraestructura Docker

Esta configuración corresponde a la branch `chore/devops-docker-compose` y deja listo el proyecto “Gestión de Actividades Deportivas” para levantar MySQL + backend Go con un `docker compose up`.

## Servicios en `docker-compose.yml`
- `mysql`: imagen `mysql:8`, credenciales desde `.env`, charset `utf8mb4`, volumen persistente `mysql_data`.
- `backend`: build desde `Dockerfile` en la raíz, expone `8080`, depende de MySQL, toma variables de entorno desde `.env`.
- `frontend` (comentado): bloque guía para agregar más adelante la SPA en React (build + Nginx) dependiendo del backend.

## Variables de entorno
Usa `.env` (creado a partir de `.env.example`) con:
- Base de datos: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`.
- App: `SERVER_PORT`, `JWT_SECRET`.
- MySQL: `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE`, `MYSQL_USER`, `MYSQL_PASSWORD`.
Dentro de Docker, el backend se conecta a la DB con `DB_HOST=mysql` y `DB_PORT=3306`.

## Cómo levantar todo
```bash
cp .env.example .env
docker compose up -d --build
```
La API quedará en `http://localhost:8080` y MySQL en `localhost:3306`.

### Conexión backend → MySQL en la red de Docker
- Host: `mysql`
- Puerto: `3306`
- Usuario/clave/base: tomados de `.env`

## Ramas y módulos relacionados
- Infra: `chore/devops-docker-compose`.
- Próximos módulos backend: `feat/backend-models-database`, `feat/backend-auth-users`, `feat/backend-activities-public`, `feat/backend-enrollments`, `feat/backend-activities-admin`.
- Otros docs: `docs/api-contract`, `docs/architecture.md`.

## Futuro frontend
Cuando exista la SPA en React, descomenta y ajusta el servicio `frontend` en `docker-compose.yml` para servirla (puerto 5173 en dev o 80 en prod con Nginx).

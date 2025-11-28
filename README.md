# Gestión de Actividades Deportivas

Backend base para un sistema de administración de actividades deportivas de un gimnasio (socios y administradores). Incluye API REST en Go + Gin, persistencia con MySQL vía GORM, autenticación JWT (pendiente TODO) y despliegue local mediante Docker Compose.

## Stack
- Go 1.22
- Gin HTTP Framework
- GORM + MySQL 8
- JWT para autenticación
- Docker & docker-compose

## Puesta en marcha
### Local
```bash
go run cmd/server/main.go
```
_Antes de ejecutar exporta las variables de entorno o crea un `.env` basado en `.env.example`._

### Docker Compose
```bash
docker compose up -d --build
```
Esto levanta `mysql`, `backend` y `frontend` conectados con las variables definidas en `docker-compose.yml`. El frontend queda disponible en `http://localhost:5173` y consume la API publicada por el backend (`http://localhost:8080/api`). Para detenerlos ejecuta `docker compose down` desde la misma carpeta.

## Variables de entorno
Revisa `.env.example` para conocer los valores mínimos:
- `SERVER_PORT`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET`

## Modelo de datos
1. `users`: socios/administradores con rol y hash de contraseña.
2. `activities`: catálogo de clases deportivas (día, horario, cupo, instructor, etc.).
3. `enrollments`: relación usuario-actividad con estado e índice único `(user_id, activity_id)`.

Las estructuras se definen en `models/` y las migraciones se ejecutan automáticamente en `database.InitDB()`.

## Plan de branches
- `main` (estable) y `develop` (integración).
- Evoluciones planificadas en `docs/branches-plan.md` (`feat/backend-*`, `chore/devops-docker-compose`, `docs/api-contract`, etc.). Cada módulo se implementará en su branch y luego se integrará en `develop` antes de llegar a `main`.

## Documentación
Consulta la carpeta `docs/` para el contrato de API, descripción de modelos, arquitectura y roadmap de branches. Un frontend en React consumirá esta API en fases posteriores.

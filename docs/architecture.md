# Arquitectura

- **Capas principales:**
  - `handlers/` (capa HTTP) utiliza Gin para recibir requests, parsear payloads y delegar al dominio.
  - `services/` concentra la lógica de negocio (autenticación, actividades, inscripciones) apoyándose en GORM.
  - `database/` expone la inicialización y migraciones automáticas hacia MySQL.
  - `models/` define las entidades persistidas.
- **Middlewares:**
  - `AuthMiddleware` valida el JWT (pendiente TODO) y carga `userID` + `role` en el contexto.
  - `AdminMiddleware` verifica que el rol sea `admin` antes de llegar a los handlers restringidos.
- **Persistencia:** MySQL 8 se conecta mediante GORM. Cada arranque ejecuta `AutoMigrate` para asegurar el esquema base.
- **Infraestructura:** Docker Compose orquesta `mysql` y `backend`. Variables de entorno viven en `.env`/`.env.example`.
- **Frontend futuro:** Una SPA en React consumirá los endpoints definidos en `docs/api-contract.md`. No se incluye código frontend en esta fase pero el backend ya expone contratos estables para integrarse después.

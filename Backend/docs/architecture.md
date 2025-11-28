# Arquitectura de la solución

## Visión general
- **Frontend:** SPA creada con React + Vite (`Frontend/`). Consume la API a través del cliente común `src/services/apiClient.js`, que toma la base URL de `VITE_API_BASE_URL` y adjunta el header `Authorization` cuando existe un JWT guardado en `localStorage` (`AuthContext.jsx`).
- **Backend:** API REST en Go (Gin). Capas principales:
  - `handlers/`: recibe las peticiones HTTP, valida payloads y arma las respuestas (incluye endpoints públicos, protegidos y de administración).
  - `services/`: encapsula la lógica de negocio (auth/JWT, actividades, inscripciones, usuarios).
  - `middlewares/`: autenticación JWT (`AuthMiddleware`), control de rol (`AdminMiddleware`) y CORS (`CORSMiddleware`) para permitir el origen del frontend (`http://localhost:5173` durante el desarrollo).
  - `database/`: inicializa GORM, ejecuta migraciones y semillas (`database/seed.go`) en entornos `APP_ENV=dev`.
  - `models/`: entidades persistidas.
- **Base de datos:** MySQL 8.0. El DSN se construye con las variables `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`. Las migraciones se ejecutan automáticamente al iniciar el backend.

## Flujo Frontend → Backend → MySQL
1. El usuario se autentica desde `Login.jsx`, que invoca `AuthContext.login`. Éste llama a `POST /api/auth/login`, almacena el JWT y los datos del usuario en `localStorage` y notifica al `ActivitiesContext` para que refresque las inscripciones (`GET /api/me/activities`).
2. Las pantallas públicas (`Home`, `Activities`) cargan el listado mediante `GET /api/activities`. El detalle (`ActivityDetail.jsx`) consulta `GET /api/activities/:id` cuando la actividad no está cacheada.
3. Al presionar “Inscribirme” se ejecuta `POST /api/activities/:id/enroll`. El backend valida cupos, actividad activa y duplicados antes de crear el registro en `enrollments`.
4. La sección “Mis actividades” (`MyActivities.jsx`) consume `GET /api/me/activities` para renderizar el DTO que arma el handler (`enrollments_handler.go`).
5. Los formularios administrativos (`AddActivity.jsx` y `EditActivity.jsx`) invocan las operaciones CRUD de `/api/admin/activities`. Todos esos endpoints pasan primero por `AuthMiddleware` y luego por `AdminMiddleware` para asegurar que el rol sea `admin`.

## Orquestación con Docker
```
┌───────────┐    HTTP     ┌─────────────┐     TCP      ┌─────────────┐
│ React SPA │ <---------> │ Go Backend  │ <--------->  │ MySQL 8     │
│ (Vite)    │ 5173        │ Gin + GORM  │ 8080         │ 3306 (3307) │
└───────────┘             └─────────────┘              └─────────────┘
```
- `docker-compose.yml` ahora expone tres servicios: `mysql`, `backend` y `frontend`. El servicio `frontend` levanta la SPA (Vite dev server dentro del contenedor, puerto `5173`) y usa `VITE_API_BASE_URL=http://localhost:8080/api`, que sigue siendo alcanzable desde el navegador aun cuando la app se sirve desde Docker. Para desarrollo basta con `docker compose up -d --build` (desde `Backend/`) y acceder a `http://localhost:5173`.
- Variables compartidas viven en `Backend/.env` y se consumen también por el Compose (usuario/contraseña de MySQL, `JWT_SECRET`, etc.). El frontend únicamente necesita `VITE_API_BASE_URL`, que ya viene seteada en el servicio.

## Consideraciones adicionales
- **CORS:** `middlewares/CORSMiddleware` habilita los métodos `GET, POST, PUT, DELETE, OPTIONS` y los headers `Content-Type, Authorization`. Hoy se permite cualquier `Origin` para simplificar el desarrollo; en producción se recomienda restringirlo.
- **Seguridad:** Las contraseñas se almacenan con `bcrypt` (helpers en `security/password.go`) y los JWT se firman con HS256 usando `JWT_SECRET`. El middleware de autenticación vuelve a consultar el usuario para reconstruir el rol antes de permitir el acceso.
- **Semillas:** Con `APP_ENV=dev` se crean usuarios de prueba (`admin@example.com`, `socia@example.com`, ambos con `contra123`) y actividades de ejemplo. Esto permite probar el flujo full-stack sin pasos manuales adicionales.

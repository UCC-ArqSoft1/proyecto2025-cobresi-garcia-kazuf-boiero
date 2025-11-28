# API Contract

La API se expone bajo HTTP/JSON y es consumida por el frontend React (Vite) configurado con la variable `VITE_API_BASE_URL` (por defecto `http://localhost:8080/api`). Cuando los servicios se levantan con Docker Compose (`Backend/docker-compose.yml`) la URL sigue siendo `http://localhost:8080/api` desde el host y `http://backend:8080/api` desde otros contenedores.

- Todas las respuestas exitosas utilizan el envoltorio `APIResponse` `{ "success": true, "message": "opcional", "data": <payload> }`, salvo los listados públicos (`GET /api/activities`) que devuelven directamente un arreglo.
- Todas las respuestas de error usan `APIError` `{ "success": false, "error": "...", "code": "opcional", "details": "debug" }`.
- Los tokens JWT tienen una vigencia de 1 hora, deben enviarse en `Authorization: Bearer <token>` y transportan `user_id` + `role` (`socio` o `admin`).

## Endpoints

> Para cada endpoint se incluye qué parte del frontend lo consume (archivo relativo a `Frontend/src`).

### Healthcheck
- **GET `/api/health`**
  - **Descripción:** indica si la API está viva. Público.
  - **Headers:** ninguno.
  - **Respuesta 200:** `{ "status": "ok" }`.
  - **Frontend:** usado por herramientas externas / scripts de despliegue (no consumido directamente por la SPA).

### Autenticación

#### POST `/api/auth/login`
- **Descripción:** autentica usuarios (`socio` o `admin`) y devuelve un JWT.
- **Body (JSON):**
  ```json
  { "email": "admin@example.com", "password": "changeme" }
  ```
- **Respuesta 200:**
  ```json
  {
    "success": true,
    "message": "Login exitoso",
    "data": {
      "token": "<jwt>",
      "user": { "id": 1, "name": "Admin", "email": "admin@example.com", "role": "admin" }
    }
  }
  ```
- **Errores frecuentes:** `401 UNAUTHORIZED` (credenciales inválidas), `400 VALIDATION_ERROR` (payload incorrecto).
- **Frontend:** `pages/Login.jsx` via `contexts/AuthContext.jsx` → `services/authService.js`.

#### POST `/api/auth/register`
- **Descripción:** registra un nuevo socio (rol fijo `socio`). No devuelve token.
- **Body:**
  ```json
  { "name": "Socia Demo", "email": "socia@example.com", "password": "changeme" }
  ```
- **Respuesta 201:** `data` contiene el usuario creado.
- **Errores:** `409 VALIDATION_ERROR` si el email ya existe.
- **Frontend:** `pages/Signup.jsx` (llama a `register` y luego realiza `login` automáticamente).

### Actividades públicas

#### GET `/api/activities`
- **Descripción:** lista actividades activas. Acepta filtros opcionales `?q=<texto>` (coincide contra título/descripción), `?category=<categoria>` y `?day=<0-6>` para día de la semana.
- **Auth:** público.
- **Respuesta 200:**
  ```json
  [
    {
      "id": 1,
      "title": "Yoga Sunrise",
      "description": "Clase matinal",
      "category": "yoga",
      "day_of_week": 1,
      "start_time": "07:30",
      "end_time": "08:30",
      "capacity": 20,
      "instructor": "Lucia Perez",
      "image_url": "",
      "is_active": true,
      "available_slots": 18,
      "enrolled_count": 2
    }
  ]
  ```
- **Frontend:** `pages/Activities.jsx` (búsqueda/listado) y precarga en `contexts/ActivitiesContext.jsx`.

#### GET `/api/activities/:id`
- **Descripción:** devuelve el detalle completo de una actividad.
- **Respuesta 200:** objeto `Activity` completo, incluyendo `available_slots` y `enrolled_count`.
- **Errores:** `400 VALIDATION_ERROR` si `:id` no es numérico, `404` si no existe.
- **Frontend:** `pages/ActivityDetail.jsx` y `pages/EditActivity.jsx` (mediante `ActivitiesContext.loadActivityById`). También usado indirectamente tras crear/editar para refrescar.

### Inscripciones y perfil del socio

#### POST `/api/activities/:id/enroll`
- **Descripción:** inscribe al usuario autenticado. Requiere que la actividad esté activa y con cupo disponible.
- **Auth:** `Authorization: Bearer <token>`.
- **Respuesta 201:** `data` contiene la inscripción (`Enrollment`).
- **Errores:** `404 ACTIVITY_NOT_FOUND`, `400 ACTIVITY_INACTIVE`, `409 ALREADY_ENROLLED`, `409 NO_CAPACITY`, `409 SCHEDULE_CONFLICT` (si ya existe una actividad con el mismo día y horarios solapados) y `401 UNAUTHORIZED` si falta token.
  - Ejemplo de solapamiento:
    ```json
    {
      "success": false,
      "error": "La actividad se solapa en dia y horario con otra inscripcion activa",
      "code": "SCHEDULE_CONFLICT"
    }
    ```
- **Frontend:** botón “Inscribirme” en `pages/ActivityDetail.jsx` mediante `ActivitiesContext.enrollInActivity`.

#### DELETE `/api/activities/:id/enroll`
- **Descripción:** desinscribe al usuario autenticado de la actividad indicada. Cambia el `status` de la inscripción a `cancelado` y libera el cupo.
- **Auth:** `Authorization: Bearer <token>`.
- **Respuesta 200:** `{ "success": true, "message": "Te desinscribiste de la actividad" }`.
- **Errores:** `404 ENROLLMENT_NOT_FOUND` si el usuario no estaba inscripto, `401 UNAUTHORIZED` por token faltante/ inválido.
- **Frontend:** botones “Desinscribirme” en `pages/MyActivities.jsx` y `pages/ActivityDetail.jsx` (`ActivitiesContext.unenrollFromActivity`).

#### GET `/api/me/activities`
- **Descripción:** lista las actividades vigentes del usuario logueado (solo actividades con inscripción `status = inscripto`).
- **Auth:** `Authorization: Bearer <token>`.
- **Respuesta 200:** `data` es un arreglo con `id`, `title`, `description`, `category`, `day_of_week`, `start_time`, `end_time`, `instructor`.
- **Frontend:** `pages/MyActivities.jsx` y verificación de inscripciones en `pages/ActivityDetail.jsx` vía `ActivitiesContext`.

### Administración de actividades (rol `admin`)
Todas requieren `Authorization: Bearer <token>` y rol `admin` (middleware `AdminMiddleware`).

#### GET `/api/admin/activities`
- **Descripción:** listado completo (activos e inactivos). Filtros: mismos que públicos + `is_active=true|false`.
- **Respuesta 200:** `APIResponse` con arreglo de actividades.
- **Frontend:** usado indirectamente al crear/editar (el contexto refresca el listado general). Para paneles más avanzados se puede reutilizar en `pages/AddActivity.jsx` o vistas futuras.

#### POST `/api/admin/activities`
- **Descripción:** crea una actividad. Todos los campos son obligatorios salvo `image_url` e `is_active` (por defecto `true`).
- **Body:**
  ```json
  {
    "title": "Funcional",
    "description": "Entrenamiento de fuerza",
    "category": "fuerza",
    "day_of_week": 2,
    "start_time": "18:00",
    "end_time": "19:00",
    "capacity": 20,
    "instructor": "Carlos Diaz",
    "image_url": "",
    "is_active": true
  }
  ```
- **Respuesta 201:** actividad creada (incluye `available_slots` y `enrolled_count` iniciales).
- **Errores:** `400 VALIDATION_ERROR` (horarios inválidos, `capacity <= 0`, etc.).
- **Frontend:** formulario `pages/AddActivity.jsx` → `ActivitiesContext.createActivity`.

#### PUT `/api/admin/activities/:id`
- **Descripción:** actualiza completamente una actividad.
- **Body:** mismo schema que `POST`.
- **Respuesta 200:** actividad actualizada con los nuevos `available_slots` calculados en base a las inscripciones activas.
- **Errores:** `404 NOT_FOUND` si la actividad no existe, `400 VALIDATION_ERROR` para datos inválidos.
- **Frontend:** `pages/EditActivity.jsx` → `ActivitiesContext.updateActivity`.

#### DELETE `/api/admin/activities/:id`
- **Descripción:** desactiva la actividad (soft delete, `is_active=false`).
- **Respuesta 200:** `{ "success": true, "message": "Actividad desactivada" }`.
- **Errores:** `404 NOT_FOUND` si el id no existe.
- **Frontend:** botón “Eliminar” en `pages/ActivityDetail.jsx` cuando el usuario es admin (`ActivitiesContext.deleteActivity`). Después se navega al listado y el contexto elimina la actividad del estado local.

### Resumen de cabeceras y puertos
| Contexto | URL base | Notas |
| --- | --- | --- |
| Desarrollo local (Go sin Docker) | `http://localhost:8080/api` | `SERVER_PORT` configurable vía `.env`. |
| Docker Compose | `http://localhost:8080/api` desde el host, `http://backend:8080/api` entre contenedores | CORS habilitado para cualquier origen (se recomienda ajustar en producción). |
| Frontend React (Vite) | `VITE_API_BASE_URL` (definido en `.env.example`) | El cliente HTTP (`src/services/apiClient.js`) agrega `Authorization` automáticamente si el usuario inició sesión. |

Los métodos permiten los encabezados `Content-Type` y `Authorization` y aceptan verbos `GET/POST/PUT/DELETE/OPTIONS`, por lo que no se requieren configuraciones adicionales al consumirlos desde el navegador.

# API contract

> Todos los endpoints devuelven `application/json`.

## Formatos base
- `APIResponse`: `{ "success": true, "message": "opcional", "data": {} }`
- `APIError`: `{ "success": false, "error": "descripcion legible", "code": "opcional", "details": "opcional" }`

## Modelos (contrato parcial)
### User (JSON)
Campos expuestos: `id`, `name`, `email`, `role`, `created_at`, `updated_at`. No se expone `password_hash`.

### Activity (JSON)
Campos expuestos: `id`, `title`, `description`, `category`, `day_of_week` (0=domingo .. 6=sabado), `start_time` (`HH:MM`), `end_time` (`HH:MM`), `capacity`, `instructor`, `image_url`, `is_active`, `created_at`, `updated_at`.

### Enrollment (JSON)
Campos expuestos: `id`, `user_id`, `activity_id`, `status`, `created_at`, `updated_at`, y opcionalmente los objetos embebidos `user` o `activity` segun uso de `Preload`.

## Endpoints actuales
### Auth
#### /api/auth/login (POST)
- Descripcion: inicio de sesion de socios o administradores.
- Body:
```json
{ "email": "ale@example.com", "password": "string" }
```
- Response 200:
```json
{
  "success": true,
  "message": "Login exitoso",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": 1,
      "name": "Ale",
      "email": "ale@example.com",
      "role": "socio"
    }
  }
}
```
- Error 401 (credenciales invalidas):
```json
{
  "success": false,
  "error": "Credenciales inválidas",
  "code": "UNAUTHORIZED"
}
```

#### /api/auth/register (POST)
- Descripcion: crea un nuevo usuario con rol `socio`.
- Body:
```json
{ "name": "Ale", "email": "ale@dominio.com", "password": "123456" }
```
- Response 201:
```json
{
  "success": true,
  "message": "Registro exitoso",
  "data": {
    "id": 10,
    "name": "Ale",
    "email": "ale@dominio.com",
    "role": "socio"
  }
}
```
- Error 409 (email duplicado):
```json
{
  "success": false,
  "error": "El email ya está registrado",
  "code": "VALIDATION_ERROR"
}
```

### /api/health (GET)
- Descripcion: verifica que la API esta viva.
- Response 200: `{ "status": "ok" }`.

### /api/activities (GET)
- Descripcion: lista actividades publicas con filtros opcionales `q`, `category`, `day`.
- Response 200:
```json
[
  {
    "id": 1,
    "title": "Funcional",
    "category": "fuerza",
    "day_of_week": 2,
    "start_time": "18:00",
    "end_time": "19:00",
    "capacity": 20,
    "instructor": "Luca",
    "is_active": true
  }
]
```

### /api/activities/:id (GET)
- Descripcion: detalle de una actividad especifica.
- Response 200: actividad completa.

### /api/activities/:id/enroll (POST)
- Descripcion: inscribe al socio o admin autenticado en una actividad activa.
- Auth: requiere header `Authorization: Bearer <token>`.
- Path params: `id` (uint) identificador de la actividad.
- Body: sin cuerpo.
- Response 201:
```json
{
  "success": true,
  "message": "Inscripcion exitosa",
  "data": {
    "id": 10,
    "user_id": 5,
    "activity_id": 3,
    "status": "inscripto",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
}
```
- Errores:
  - 404 `APIError` con `code: "ACTIVITY_NOT_FOUND"` si la actividad no existe.
  - 400 `APIError` con `code: "ACTIVITY_INACTIVE"` si la actividad esta inactiva.
  - 409 `APIError` con `code: "ALREADY_ENROLLED"` si el usuario ya esta inscripto.
  - 409 `APIError` con `code: "NO_CAPACITY"` si no hay cupos.
  - 401 `APIError` `code: "UNAUTHORIZED"` si falta o es invalido el token.

### /api/me/activities (GET)
- Descripcion: lista actividades en las que el usuario autenticado esta inscripto (status `inscripto`).
- Auth: requiere header `Authorization: Bearer <token>`.
- Response 200:
```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "title": "Zumba",
      "description": "Clases de ritmos latinos",
      "category": "zumba",
      "day_of_week": 1,
      "start_time": "19:30",
      "end_time": "20:30",
      "instructor": "Maria"
    }
  ]
}
```

## Admin Activities (solo rol admin)
Todos requieren `Authorization: Bearer <token>` y los middlewares `Auth` + `Admin`.

### /api/admin/activities (GET)
- Descripcion: listado completo para administracion, permite filtrar por `q`, `category`, `day` y `is_active=true|false`.
- Response 200:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Funcional",
      "category": "fuerza",
      "day_of_week": 2,
      "start_time": "18:00",
      "end_time": "19:00",
      "capacity": 20,
      "instructor": "Luca",
      "is_active": true
    }
  ]
}
```

### /api/admin/activities (POST)
- Descripcion: crea una nueva actividad.
- Body:
```json
{
  "title": "Pilates",
  "description": "Clases grupales",
  "category": "movilidad",
  "day_of_week": 4,
  "start_time": "09:00",
  "end_time": "10:00",
  "capacity": 12,
  "instructor": "Lola",
  "image_url": "https://cdn.example/pilates.png",
  "is_active": true
}
```
- Response 201:
```json
{
  "success": true,
  "message": "Actividad creada",
  "data": {
    "id": 12,
    "title": "Pilates",
    "category": "movilidad",
    "day_of_week": 4,
    "start_time": "09:00",
    "end_time": "10:00",
    "capacity": 12,
    "instructor": "Lola",
    "is_active": true
  }
}
```
- Error 400 (validacion):
```json
{
  "success": false,
  "error": "capacity debe ser mayor a 0",
  "code": "VALIDATION_ERROR"
}
```

### /api/admin/activities/:id (PUT)
- Descripcion: actualiza completamente una actividad existente.
- Body: mismo formato que `POST`.
- Response 200:
```json
{
  "success": true,
  "message": "Actividad actualizada",
  "data": {
    "id": 12,
    "title": "Pilates avanzado",
    "description": "Nueva descripcion",
    "category": "movilidad",
    "day_of_week": 4,
    "start_time": "09:30",
    "end_time": "10:30",
    "capacity": 15,
    "instructor": "Lola",
    "image_url": "https://cdn.example/pilates-2.png",
    "is_active": true
  }
}
```
- Error 404:
```json
{
  "success": false,
  "error": "Actividad no encontrada",
  "code": "NOT_FOUND"
}
```

### /api/admin/activities/:id (DELETE)
- Descripcion: desactiva (soft-delete) una actividad (`is_active=false`).
- Response 200:
```json
{
  "success": true,
  "message": "Actividad desactivada"
}
```
- Error 404: `APIError` equivalente al de `PUT`.

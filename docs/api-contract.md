# API contract

> Todos los endpoints devuelven `application/json`.

## Formatos base
- `APIResponse`: `{ "success": true, "message": "opcional", "data": {} }`
- `APIError`: `{ "success": false, "error": "descripcion legible", "code": "opcional" }`

## Modelos (contrato parcial)
### User (JSON)
Campos expuestos: `id`, `name`, `email`, `role`, `created_at`, `updated_at`. No se expone `password_hash`.

### Activity (JSON)
Campos expuestos: `id`, `title`, `description`, `category`, `day_of_week` (0=domingo .. 6=sabado), `start_time` (`HH:MM`), `end_time` (`HH:MM`), `capacity`, `instructor`, `image_url`, `is_active`, `created_at`, `updated_at`.

### Enrollment (JSON)
Campos expuestos: `id`, `user_id`, `activity_id`, `status`, `created_at`, `updated_at`, y opcionalmente los objetos embebidos `user` o `activity` segun uso de `Preload`.

## Endpoints actuales
### /api/health (GET)
- Descripcion: verifica que la API esta viva.
- Response 200: `{ "status": "ok" }`.

### /api/auth/login (POST)
- Descripcion: inicio de sesion de socios o administradores.
- Body:
```json
{ "email": "ale@example.com", "password": "string" }
```
- Response 200:
```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "name": "Ale",
    "email": "ale@example.com",
    "role": "socio"
  }
}
```

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
- Descripcion: inscribe al socio autenticado a una actividad.
- Auth: requiere token Bearer.
- Response 201:
```json
{
  "id": 10,
  "user_id": 5,
  "activity_id": 3,
  "status": "inscripto"
}
```

### /api/me/activities (GET)
- Descripcion: lista actividades del socio autenticado.
- Auth: requiere token Bearer.
- Response 200:
```json
[
  {
    "id": 12,
    "user_id": 5,
    "activity_id": 3,
    "status": "inscripto",
    "activity": { "title": "Yoga" }
  }
]
```

### /api/admin/activities (POST)
- Descripcion: crea actividades. Solo admin.
- Body:
```json
{
  "title": "Pilates",
  "description": "",
  "category": "movilidad",
  "day_of_week": 4,
  "start_time": "09:00",
  "end_time": "10:00",
  "capacity": 12,
  "instructor": "Lola",
  "is_active": true
}
```
- Response 201: actividad creada.

### /api/admin/activities/:id (PUT)
- Descripcion: actualiza una actividad existente. Solo admin.
- Response 200: actividad actualizada.

### /api/admin/activities/:id (DELETE)
- Descripcion: desactiva o elimina una actividad. Solo admin.
- Response 204: sin cuerpo.

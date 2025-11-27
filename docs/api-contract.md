# API Contract

> Todos los endpoints devuelven `application/json`.

## /api/health (GET)
- **Descripción:** Verifica que la API esté viva.
- **Request Body:** _N/A_
- **Response 200:**
```json
{ "status": "ok" }
```
- **Errores comunes:** 500 cuando el servidor no puede responder.

## /api/auth/login (POST)
- **Descripción:** Inicia sesión de socios o administradores.
- **Request Body:**
```json
{ "email": "ale@example.com", "password": "string" }
```
- **Response 200:**
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
- **Errores comunes:** 400 (payload inválido), 401 (credenciales inválidas).

## /api/activities (GET)
- **Descripción:** Lista actividades públicas con filtros opcionales `q`, `category`, `day`.
- **Request Body:** _N/A_
- **Response 200:**
```json
[
  {
    "id": 1,
    "title": "Funcional",
    "category": "fuerza",
    "day_of_week": 2,
    "start_time": "18:00:00",
    "end_time": "19:00:00",
    "capacity": 20,
    "instructor": "Lucía",
    "is_active": true
  }
]
```
- **Errores comunes:** 500 (error interno al consultar DB).

## /api/activities/:id (GET)
- **Descripción:** Detalle de una actividad específica.
- **Request Body:** _N/A_
- **Response 200:** Actividad completa.
- **Errores comunes:** 400 (id inválido), 404 (actividad no existe).

## /api/activities/:id/enroll (POST)
- **Descripción:** Inscribe al socio autenticado a una actividad.
- **Auth:** Requiere token Bearer.
- **Request Body:** _N/A_ (el ID se toma de la URL).
- **Response 201:**
```json
{
  "id": 10,
  "user_id": 5,
  "activity_id": 3,
  "status": "inscripto"
}
```
- **Errores comunes:** 400 (actividad inválida o sin cupo), 401 (sin token), 409 (duplicado, TODO), 500 (otros errores).

## /api/me/activities (GET)
- **Descripción:** Lista las actividades a las que está inscripto el socio autenticado.
- **Auth:** Requiere token Bearer.
- **Response 200:**
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
- **Errores comunes:** 401 (sin token), 500 (error interno).

## /api/admin/activities (POST)
- **Descripción:** Crea actividades. Solo admin.
- **Auth:** Token Bearer + rol admin.
- **Request Body:**
```json
{
  "title": "Pilates",
  "description": "",
  "category": "movilidad",
  "day_of_week": 4,
  "start_time": "09:00:00",
  "end_time": "10:00:00",
  "capacity": 12,
  "instructor": "Lola",
  "is_active": true
}
```
- **Response 201:** Actividad creada.
- **Errores comunes:** 400 (datos inválidos), 401/403 (sin permisos), 500 (error DB).

## /api/admin/activities/:id (PUT)
- **Descripción:** Actualiza una actividad existente. Solo admin.
- **Auth:** Token Bearer + rol admin.
- **Response 200:** Actividad actualizada.
- **Errores comunes:** 400 (payload o id inválido), 401/403 (sin permisos), 404 (no existe), 500 (error DB).

## /api/admin/activities/:id (DELETE)
- **Descripción:** Desactiva o elimina una actividad. Solo admin.
- **Auth:** Token Bearer + rol admin.
- **Response 204:** _sin cuerpo_.
- **Errores comunes:** 400 (id inválido), 401/403 (sin permisos), 404 (no existe), 500 (error DB).

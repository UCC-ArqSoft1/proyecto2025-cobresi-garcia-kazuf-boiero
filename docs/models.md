# Modelos de Dominio

## User
| Campo | Tipo Go | Tipo MySQL | Descripción |
| --- | --- | --- | --- |
| ID | `uint` | `bigint unsigned` | PK auto increment |
| Name | `string` | `varchar(255)` | Nombre del socio o admin |
| Email | `string` | `varchar(255)` `UNIQUE NOT NULL` | Email de login |
| PasswordHash | `string` | `varchar(255)` `NOT NULL` | Hash de contraseña |
| Role | `string` | `enum('socio','admin')` | Perfil del usuario |
| CreatedAt | `time.Time` | `datetime` | Timestamp de creación |
| UpdatedAt | `time.Time` | `datetime` | Última actualización |

**Relaciones:** `User` tiene muchas `Enrollment`.

## Activity
| Campo | Tipo Go | Tipo MySQL | Descripción |
| --- | --- | --- | --- |
| ID | `uint` | `bigint unsigned` | PK |
| Title | `string` | `varchar(255)` | Nombre comercial |
| Description | `string` | `text` | Detalle de la clase |
| Category | `string` | `varchar(100)` | Tipo (yoga, fuerza, etc.) |
| DayOfWeek | `int` | `tinyint` | Día 0=dom…6=sáb |
| StartTime | `time.Time` | `time` | Hora inicio |
| EndTime | `time.Time` | `time` | Hora fin |
| Capacity | `int` | `int` | Cupos |
| Instructor | `string` | `varchar(255)` | Nombre instructor |
| ImageURL | `*string` | `varchar(512)` nullable | Imagen opcional |
| IsActive | `bool` | `tinyint(1)` | Activa o no |
| CreatedAt | `time.Time` | `datetime` | Creación |
| UpdatedAt | `time.Time` | `datetime` | Modificación |

**Relaciones:** `Activity` tiene muchas `Enrollment`.

## Enrollment
| Campo | Tipo Go | Tipo MySQL | Descripción |
| --- | --- | --- | --- |
| ID | `uint` | `bigint unsigned` | PK |
| UserID | `uint` | `bigint unsigned` `FK` | Referencia a `users.id` |
| ActivityID | `uint` | `bigint unsigned` `FK` | Referencia a `activities.id` |
| Status | `string` | `enum('inscripto','cancelado')` | Estado |
| CreatedAt | `time.Time` | `datetime` | Creación |
| UpdatedAt | `time.Time` | `datetime` | Modificación |

**Relaciones:** `Enrollment` pertenece a `User` y `Activity`. Restricción única `(user_id, activity_id)` evita duplicados.

# Modelos de dominio

## User
| Campo | Tipo Go | Tipo MySQL | Descripcion |
| --- | --- | --- | --- |
| ID | `uint` | `bigint unsigned` | PK auto increment |
| Name | `string` | `varchar(255) not null` | Nombre del socio o admin |
| Email | `string` | `varchar(255) not null unique` | Correo de login |
| PasswordHash | `string` | `varchar(255) not null` | Hash de contrasena |
| Role | `string` | `varchar(20) not null` | Perfil `socio` o `admin` |
| CreatedAt | `time.Time` | `datetime` | Creacion |
| UpdatedAt | `time.Time` | `datetime` | Ultima actualizacion |

**Relaciones:** `User` tiene muchas `Enrollment` (`enrollments.user_id`).

## Activity
| Campo | Tipo Go | Tipo MySQL | Descripcion |
| --- | --- | --- | --- |
| ID | `uint` | `bigint unsigned` | PK |
| Title | `string` | `varchar(255) not null` | Nombre comercial |
| Description | `string` | `text` | Detalle de la clase |
| Category | `string` | `varchar(100) not null` | Tipo (yoga, fuerza, etc.) |
| DayOfWeek | `int` | `tinyint not null` | Dia 0=domingo .. 6=sabado |
| StartTime | `string` | `varchar(8) not null` | Hora inicio formato `HH:MM` |
| EndTime | `string` | `varchar(8) not null` | Hora fin formato `HH:MM` |
| Capacity | `int` | `int not null` | Cupos totales |
| Instructor | `string` | `varchar(255) not null` | Nombre instructor |
| ImageURL | `string` | `varchar(512)` | Imagen opcional |
| IsActive | `bool` | `tinyint(1)` | Activa o no |
| CreatedAt | `time.Time` | `datetime` | Creacion |
| UpdatedAt | `time.Time` | `datetime` | Modificacion |

**Relaciones:** `Activity` tiene muchas `Enrollment` (`enrollments.activity_id`).

## Enrollment
| Campo | Tipo Go | Tipo MySQL | Descripcion |
| --- | --- | --- | --- |
| ID | `uint` | `bigint unsigned` | PK |
| UserID | `uint` | `bigint unsigned not null FK` | Referencia a `users.id` |
| ActivityID | `uint` | `bigint unsigned not null FK` | Referencia a `activities.id` |
| Status | `string` | `varchar(20) not null default 'inscripto'` | Estado de la inscripcion |
| CreatedAt | `time.Time` | `datetime` | Creacion |
| UpdatedAt | `time.Time` | `datetime` | Modificacion |

**Relaciones:** `Enrollment` pertenece a `User` y `Activity` (OnDelete RESTRICT). Conceptualmente no debe haber dos inscripciones activas para la misma combinacion user/activity.

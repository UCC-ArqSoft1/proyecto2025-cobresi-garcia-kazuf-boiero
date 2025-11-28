# Modelos de dominio

El backend utiliza GORM sobre MySQL 8. A continuación se detallan las tablas, struct en Go y la forma en que cada entidad se serializa hacia el frontend.

## User
Usuarios finales (socios o administradores).

### Esquema MySQL
```sql
CREATE TABLE users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(20) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### Struct Go (models/user.go)
```go
type User struct {
    ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name         string    `gorm:"size:255;not null" json:"name"`
    Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    Role         string    `gorm:"size:20;not null" json:"role"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    Enrollments  []Enrollment `gorm:"foreignKey:UserID" json:"-"`
}
```
Solo se exponen los campos `id`, `name`, `email`, `role` y timestamps; `password_hash` nunca viaja a la API.

### JSON típico
```json
{
  "id": 1,
  "name": "Admin",
  "email": "admin@example.com",
  "role": "admin",
  "created_at": "2024-11-01T15:00:00Z",
  "updated_at": "2024-11-01T15:00:00Z"
}
```

## Activity
Actividades deportivas ofrecidas a los socios.

### Esquema MySQL
```sql
CREATE TABLE activities (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100) NOT NULL,
  day_of_week TINYINT NOT NULL,
  start_time VARCHAR(8) NOT NULL,
  end_time VARCHAR(8) NOT NULL,
  capacity INT NOT NULL,
  instructor VARCHAR(255) NOT NULL,
  image_url VARCHAR(512),
  is_active TINYINT(1) DEFAULT 1,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### Struct Go (models/activity.go)
```go
type Activity struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Title       string    `gorm:"size:255;not null" json:"title"`
    Description string    `gorm:"type:text" json:"description"`
    Category    string    `gorm:"size:100;not null" json:"category"`
    DayOfWeek   int       `gorm:"not null" json:"day_of_week"`
    StartTime   string    `gorm:"size:8;not null" json:"start_time"`
    EndTime     string    `gorm:"size:8;not null" json:"end_time"`
    Capacity    int       `gorm:"not null" json:"capacity"`
    Instructor  string    `gorm:"size:255;not null" json:"instructor"`
    ImageURL    string    `gorm:"size:512" json:"image_url"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
    AvailableSlots int    `gorm:"-" json:"available_slots"`
    EnrolledCount  int    `gorm:"-" json:"enrolled_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Enrollments []Enrollment `gorm:"foreignKey:ActivityID" json:"-"`
}
```

### JSON típico
```json
{
  "id": 3,
  "title": "Spinning",
  "description": "Cardio de alta intensidad en bicicleta fija.",
  "category": "cardio",
  "day_of_week": 4,
  "start_time": "19:30",
  "end_time": "20:15",
  "capacity": 15,
  "instructor": "Agus Flores",
  "image_url": "",
  "is_active": true,
  "available_slots": 12,
  "enrolled_count": 3,
  "created_at": "2024-10-05T12:00:00Z",
  "updated_at": "2024-10-05T12:00:00Z"
}
```
Los listados públicos (`GET /api/activities`) excluyen actividades con `is_active = false`. Las operaciones admin pueden filtrar por ese campo y modificarlo (soft-delete). `available_slots = max(capacity - enrolled_count, 0)` se calcula al vuelo y permite al frontend mostrar cupos dinámicos sin tener que contar inscripciones.

## Enrollment
Relación entre un `User` y una `Activity`.

### Esquema MySQL
```sql
CREATE TABLE enrollments (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  activity_id BIGINT UNSIGNED NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'inscripto',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  CONSTRAINT fk_enrollment_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_enrollment_activity FOREIGN KEY (activity_id) REFERENCES activities(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
```

### Struct Go (models/enrollment.go)
```go
type Enrollment struct {
    ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID     uint      `gorm:"not null;index" json:"user_id"`
    ActivityID uint      `gorm:"not null;index" json:"activity_id"`
    Status     string    `gorm:"size:20;not null;default:'inscripto'" json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`

    User     User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
    Activity Activity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
```

### JSON típico
```json
{
  "id": 22,
  "user_id": 2,
  "activity_id": 3,
  "status": "inscripto",
  "created_at": "2024-11-10T09:15:00Z",
  "updated_at": "2024-11-10T09:15:00Z"
}
```

### Reglas de negocio
- Solo se permite una inscripción activa (`status = 'inscripto'`) por combinación `user_id + activity_id`. El servicio valida duplicados antes de crear un registro nuevo.
- Las actividades inactivas (`is_active = false`) no aceptan nuevas inscripciones.
- El cupo se controla comparando el número de inscripciones activas con `activity.capacity`. Ante overflow se responde con `NO_CAPACITY`. Las desinscripciones actualizan el `status` a `cancelado` para conservar el historial, y solo se contabilizan los registros `inscripto`.
- Un usuario no puede inscribirse en dos actividades que se solapen (mismo `day_of_week` y horarios entrelazados). Ante esta validación se responde con `SCHEDULE_CONFLICT`.
- El endpoint `/api/me/activities` devuelve un DTO liviano que incluye los campos de la actividad asociados a cada inscripción para facilitar el renderizado en React.

# Plan de Branches

## Ramas base
- `main`: rama estable lista para releases.
- `develop`: rama de integración continua del backend.

## Fase 1
- `feat/backend-setup`: esqueleto inicial (este trabajo).

## Fase 2 y 3
- `chore/devops-docker-compose`: mejoras de infraestructura y pipelines.
- `feat/backend-models-database`: evolución de modelos, migraciones y seeds.
- `feat/backend-auth-users`: autenticación completa y gestión de usuarios.
- `feat/backend-activities-public`: lógica avanzada de actividades públicas.
- `feat/backend-enrollments`: reglas de inscripción y consultas de mis actividades.
- `feat/backend-activities-admin`: herramientas de administración para actividades.
- `docs/api-contract`: mantenimiento de la documentación funcional.

> El frontend en React se desarrollará en ramas específicas más adelante. Por ahora el repositorio se enfoca en backend, base de datos, Docker e infraestructura.

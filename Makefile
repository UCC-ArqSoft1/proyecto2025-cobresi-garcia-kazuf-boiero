.PHONY: up down logs ps db clean

# Levanta toda la pila en segundo plano (con build forzado).
up:
	docker compose up -d --build

# Detiene los contenedores.
down:
	docker compose down

# Muestra logs del backend.
logs:
	docker compose logs -f backend

# Estado rápido de los servicios.
ps:
	docker compose ps

# Abre cliente MySQL dentro del contenedor (usa credenciales de .env).
db:
	docker compose exec -it mysql sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE'

# Detiene y elimina contenedores + volumenes (cuidado, borra datos locales).
clean:
	docker compose down -v

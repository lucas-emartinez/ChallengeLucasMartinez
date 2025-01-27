# Variables
SERVICE ?= api

build:
	# Build el servicio especificado (por defecto 'api')
	docker build --build-arg SERVICE=$(SERVICE) -t myapp:$(SERVICE) .

run:
	# Ejecuta el contenedor del servicio especificado (por defecto 'api')
	docker run -p 8080:8080 myapp:$(SERVICE)

run-api:
	# Levanta el contenedor para el servicio 'api'
	$(MAKE) run SERVICE=api

run-worker:
	# Levanta el contenedor para el servicio 'worker'
	$(MAKE) run SERVICE=worker

run-consumer:
	# Levanta el contenedor para el servicio 'consumer'
	$(MAKE) run SERVICE=consumer

# Comando para crear y correr el contenedor con docker-compose
docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down

# ChallengeLucasMartinez
# Desafío Técnico Ualá

Este proyecto es una implementación simplificada de una plataforma similar a Twitter.
Permite a los usuarios publicar tweets, seguir a otros usuarios y visualizar su timeline.

## Tecnologías utilizadas

- **Go**: Lenguaje principal de la aplicación.
- **Kafka**: Broker de mensajes para procesamiento de tweets.
- **Redis**: Almacenamiento en memoria para optimizar las lecturas.
- **Docker & Docker Compose**: Para la contenedorización y fácil despliegue de la aplicación.

## Arquitectura

- Cuando se publica un tweet, se envía un mensaje a un tópico de Kafka.
- Un worker consume los mensajes del tópico y almacena los tweets en Redis.
- Al consultar el timeline de un usuario, se obtienen los tweets de los usuarios seguidos desde Redis.
- Para los follows se guarda en memoria el usuario y los usuarios que sigue. A futuro se podría guardar 
- en una BD de tipo NoSQL como Cassandra, incluso creando una tabla especializada para los follows.
- Se utilizó Redis para almacenar los tweets y los follows, ya que es una base de datos en memoria y es muy rápida para las lecturas.
- Para fines prácticos se utilizó una sola API, pero en un escenario real se podría dividir en microservicios 
- para escalar y mantener la arquitectura desacoplada.


## Endpoints

La API expone los siguientes endpoints en `localhost:8080`:

| Método | Endpoint | Descripción |
|--------|---------|-------------|
| POST   | `/api/tweets` | Permite a los usuarios publicar un tweet. |
| POST   | `/api/follow` | Permite a un usuario seguir a otro usuario. |
| GET    | `/api/timeline/:userID` | Obtiene el timeline de un usuario en base a los usuarios seguidos. |

## Requisitos previos

- Docker y Docker Compose instalados.
- Puerto `8080` libre.

## Instalación y Ejecución

1. Clonar el repositorio:

   ```sh
   git clone https://github.com/lucas-emartinez/ChallengeLucasMartinez.git
    ```
   
2. Hacer el build y levantar los servicios con Docker Compose:

   ```sh
   docker-compose up --build
   ```

3. La aplicación estará disponible en `localhost:8080`.

## Se dejó en el root la colección de Postman para probar los endpoints

La idea sería ejecutar primero el Follow, luego publicar el Tweet desde el usuario seguido y como paso final ver el el tmeline del usuario follower.

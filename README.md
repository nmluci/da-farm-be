# D'Aqua Farm Backend

## Summary
This repository demonstrate a simple backend to manage an aquaculture farms which consist of Farms and Ponds

## Tech Stacks
### This backend uses following stacks

| Stack | Tech | 
| --- | --- | 
| Backend | Go 1.22.5 | 
| Database | PostgreSQL 16.3, alpine docker variant | 
| API Documentation | [Swagger API](http://localhost:7780/api/swagger) |
| Base Image | Linux Alpine |
| Base API | localhost:7780/api/v1/farm, localhost:7780/api/v1/ponds |

## How to run?
### Using docker
Below are the step required to run the backend inside a Docker
1. Make sure docker is installed in your system
2. (optional) change the default configuration for both DB and Backend inside **docker-compose.yml** file, and make a note for each changes
3. Run **`docker compose build backend db`** to build both DB layer and the Backend
4. Run **`docker compose up`** to start the each container
5. Backend can be accessed on **localhost:7780** (if not changed on step 2) with API Documented in the Swagger

### Running directly on Host
To run the backend directly on machine instead of using Docker
1. Make sure docker is installed in your system (since it will be needed to run the DB, unless you have existing storage solution that based on PostgreSQL)
2. Make a new file named **`.env`** inside **config** folder with value describe in [Configuration Section](#configuration). Make sure to set POSTGRES_* variable according to available resource's configuration
3. (optional) Run **`docker compose up -d --build db`** to build and run DB layer for storage solution. *If you have existing PostgreSQL instance, skip this step, and use it instead*
4. Run **`go run cmd/server`** to start the backend on local machine
5. Backend can be accessed on **localhost:7780** (if not changed on step 2) with API Documented in the Swagger

## Configuration 
| ENV Variable | Value | Default |
| ------------ | ----- | ------- |
| SVC_NAME     | Name of service | da-farm-be |
| SVC_ADDRESS  | Address of service | :7780 |
| POSTGRES_ADDRESS | PGSQL address | localhost:5432 or db:5432 |
| POSTGRES_USERNAME | PGSQL Username | postgres |
| POSTGRES_PASSWORD | PGSQL password | postgres |
| POSTGRES_DATABASE | PGSQL database name | aqua_db |
| SWAGGER_HOST | Host Baseapi to be used by Swagger to access API | localhost:7780 |

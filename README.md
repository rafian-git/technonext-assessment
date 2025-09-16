# TechnoNext Assessment

This repository implements the assessment using:
- **Go** , gRPC
- **JWT** for auth 
- **Envoy** gRPC-JSON transcoding to expose REST paths specified in the task
- **PostgreSQL** for persistence
- **go-pg** as ORM
- **Redis** for checking if login token is blacklisted (logout)
- **Docker & docker-compose** for one-command local run

## Endpoints (via Envoy/REST)
[click here](https://documenter.getpostman.com/view/13327243/2sB3HqKJsj#9a2d02a5-d248-4310-9f01-50205a35d6c3) for the Postman Collection .

- `POST /api/v1/login` → issue JWT (token_type `Bearer`)
- `POST /api/v1/orders` → create order (requires `Authorization: Bearer <token>`)
- `GET  /api/v1/orders/all?transfer_status=1&archive=0&limit=5&page=1` → list orders (pagination handled )
- `PUT  /api/v1/orders/{CONSIGNMENT_ID}/cancel` → cancel order 
- `POST /api/v1/logout` → manages logout by blacklisting the token in Redis 
 

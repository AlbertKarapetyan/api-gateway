# API Gateway with Load Balancer

## Description
This API Gateway is built in Go and provides load balancing between multiple backend services. It supports multiple load balancing strategies, including:

- **Round Robin**: Requests are distributed evenly across backend servers.
- **Least Connections**: Requests are directed to the backend with the fewest active connections.

## Features
- Dynamic backend server registration
- Configurable load balancing strategies
- Reverse proxy for handling client requests
- Health checks for backend services
- Test servers for development and debugging
- Dynamic Routing with Auto-Reload: Automatically reloads routing configuration without restarting the server.
- JWT Authentication Middleware for secure access control
- Public Routes support for endpoints that do not require authentication

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) (version 1.18 or later)

### Steps
1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd api-gateway
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Configure the gateway by editing `config.json`:
   ```json
   {
     "load_balancer": "round_robin",
     "health_check_interval": 5,
     
     "servers": {
        "user": [
          "http://localhost:8081",
          "http://localhost:8082"
        ],
        "wallet": [
          "http://localhost:8083",
          "http://localhost:8084"
        ]
      },
      "routes": {
        "user/signin": "/api/auth",
        "user/signup": "/api/register",
        "wallet/get_balance": "/api/get_balance"
      },
      "public_routes": {
        "/user/signin": true,
        "/user/signup": true
      },
      "secret_key": "your-secret-key"
   }
   ```
4. Start the gateway:
   ```sh
   go run main.go
   ```

## Running Test Servers
For testing, you can start mock backend servers:
```sh
cd testServers/server1 && go run main.go &
cd testServers/server2 && go run main.go &
cd testServers/server3 && go run main.go &
```

## Usage
Once running, send requests to the gateway:
```sh
curl http://localhost:8080/user/signin
```
The gateway will forward the request to a backend server based on the load balancing strategy.
---
## JWT Authentication Middleware

The API Gateway includes JWT-based authentication to secure access to protected endpoints. Requests must include a valid JWT token in the `Authorization` header:

### Middleware Implementation
- The middleware extracts the JWT token from the request header.
- It verifies the token signature using the configured secret key.
- If the token is valid, the request proceeds to the backend service.
- If the token is missing or invalid, the request is rejected with a `401 Unauthorized` response.

### Example Request with JWT
```sh
curl -H "Authorization: Bearer <your-jwt-token>" http://localhost:8080/protected-route
```

### Public Routes
Some endpoints, such as authentication-related routes, should be accessible without requiring a JWT token. These routes are specified in the `config.json` file under `public_routes`:

```json
"public_routes": {
  "/user/signin": true,
  "/user/signup": true
}
```

The middleware will bypass authentication for these routes, allowing unauthenticated users to access them.

---
## Dynamic Routing with Auto-Reload
The API Gateway supports **dynamic routing** and **auto-reloading** of the configuration without requiring a server restart. This feature allows you to update the `config.json` file while the server is running, and the gateway will automatically apply the changes.

### How It Works
1. **Dynamic Routing:**
- Routes are defined in the `config.json` file under the `routes` section.
- Each route maps an API gateway path (e.g., `/user/signin`) to a backend service path (e.g., `/api/auth`).
- The gateway dynamically routes requests to the appropriate backend service based on the configuration.

2. **Auto-Reload:**
- The gateway monitors the `config.json` file for changes.
- When a change is detected, the gateway reloads the configuration and updates the routing and server lists.
- No server restart is required—changes take effect within a few seconds.

### Example
1. Update `config.json` to add a new route:
```json
{
  "routes": {
    "user/signin": "/api/auth",
    "user/signup": "/api/register",
    "wallet/get_balance": "/api/get_balance",
    "wallet/transactions": "/api/transactions" // New route
  }
}
```
2. Save the file. The gateway will automatically reload the configuration and start routing requests for `/wallet/transactions` to the specified backend path.

### Benefits
- **Zero Downtime:** Update routes and servers without restarting the gateway.
- **Flexibility:** Easily add, remove, or modify routes and backend servers.
- **Scalability:** Supports multiple services and routes dynamically.

---

## License
MIT License


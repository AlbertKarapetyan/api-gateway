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
     "servers": [
       { "url": "http://localhost:8081" },
       { "url": "http://localhost:8082" }
     ],
     "load_balancer": "round_robin",
     "health_check_interval": 5
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
curl http://localhost:8080/api
```
The gateway will forward the request to a backend server based on the load balancing strategy.

## API Documentation
### Endpoints
#### `GET /api`
- Forwards request to a backend service based on load balancing strategy.
- Example:
  ```sh
  curl http://localhost:8080/api
  ```
---  
## License
MIT License


# IPChecker Service

IPChecker is a Go-based service that checks whether a given IP address belongs to an allowed set of countries. It exposes both HTTP (via the Gin framework) and gRPC endpoints. The service uses a MaxMind GeoLite2 database to map IP addresses to countries.

## Features

### HTTP Endpoint

    POST /api/v1/ip-check accepts a JSON payload containing an IP address and a list of allowed countries.

    Returns whether the IP is allowed (true/false) and the ISO country code (e.g., "US").

### gRPC Service

    ipchecker.v1.IPChecker/CheckIP receives an IP address and allowed countries.

    Returns whether the IP is allowed and the resolved country code.

### Swagger Documentation

    Served at http://<host>:8080/swagger/index.html (by default).

    Provides an interactive UI to test the HTTP endpoint.

### Docker & Kubernetes Ready

    Dockerfile and docker-compose.yml included for containerized local deployment.

    Kubernetes YAML file for easy deployment to an existing cluster.

## Project Structure
```bash
ipchecker/
├── cmd/
│   └── ipchecker/
│       └── main.go                   # Application entrypoint
├── docs/
│   ├── docs.go                       # Swagger documentation initialization
│   ├── swagger.json                  # Generated Swagger documentation (JSON)
│   └── swagger.yaml                  # Generated Swagger documentation (YAML)
├── internal/
│   ├── config/
│   │   └── config.go                 # Application configuration (port, DB path, etc.)
│   ├── dtos/
│   │   └── ip.go                     # Data Transfer Objects (DTOs) for IP checking
│   ├── geo/
│   │   ├── geolookup.go              # GeoLookup service implementation using MaxMind DB
│   │   └── mock_geo.go               # Mock GeoLookup service for unit tests
│   ├── grpcserver/
│   │   ├── ipchecker_grpc.go         # gRPC IPChecker service implementation
│   │   └── ipchecker_grpc_test.go    # gRPC service unit tests
│   ├── handler/
│   │   ├── iphandler.go              # HTTP handler (Gin) for IP checking
│   │   └── iphandler_test.go         # HTTP handler unit tests
│   ├── logger/
│   │   └── logger.go                 # Logger setup using Zap
│   ├── middleware/
│   │   ├── gin_logger.go             # Middleware for HTTP request logging and recovery
│   │   └── grpc_logger.go            # Middleware interceptors for gRPC request logging
│   └── server/
│       ├── appserver.go              # Combined HTTP and gRPC servers with common dependencies
│       ├── grpcserver.go             # gRPC server setup and configuration
│       ├── httpserver.go             # HTTP (Gin) server setup and configuration
│       └── router.go                 # HTTP route definitions and registrations
├── proto/
│   ├── ipchecker.proto               # Protocol Buffers definitions for gRPC service
│   ├── ipchecker.pb.go               # Generated protobuf message types
│   └── ipchecker_grpc.pb.go          # Generated gRPC service bindings
├── Dockerfile                        # Dockerfile for containerizing the application
├── docker-compose.yaml               # Docker Compose file to orchestrate services
├── GeoLite2-Country.mmdb             # MaxMind geo database (do NOT track if license restricts)
├── go.mod                            # Go module dependencies specification
├── go.sum                            # Checksum file for Go module dependencies
└── coverage.out                      # Test coverage report (generated via "go test")
```

## Prerequisites

    Go 1.20+ (or as specified in go.mod)

    GeoLite2-Country.mmdb file from MaxMind

    Docker (optional, if using containers)

    kubectl & a Kubernetes cluster (optional, if deploying to K8s)

## Installation & Running Locally

1. Clone the Repository
```
git clone https://github.com/justfairdev/ipchecker.git
cd ipchecker
```
2. Download the GeoLite2 Database (Already done)

    Ref: https://dev.maxmind.com/geoip/geoip2/geolite2/
3. Build & Run
```
go build -o ipchecker ./cmd/ipchecker
./ipchecker
```

## Usage

### HTTP Endpoint

1. **POST /api/v1/ip-check**

    Request Body (JSON):
    ```
    {
    "ip_address": "128.101.101.101",
    "allowed_countries": ["US", "CA"]
    }
    ```

    Response (JSON):
    ```
    {
    "allowed": true,
    "country": "US"
    }
    ```

2. **Swagger UI**

    Access at http://localhost:8080/swagger/index.html to test the API interactively.

### gRPC Endpoint

    ```
    grpcurl -plaintext -d '{"ip_address":"1.1.1.1","allowed_countries":["US","CA"]}' \
    localhost:50051 ipchecker.v1.IPChecker/CheckIP
    ```
## Testing

1. Test HTTP Handlers
    ```
    go test -v ./internal/handler
    ```
2. Test gRPC
    ```
    go test -v ./internal/grpcserver
    ```
## Docker

Build Locally
```
docker build -t justfairdev/ipchecker .
```
Run Container
```
docker run -p 8080:8080 -p 50051:50051 justfairdev/ipchecker
```
Docker Compose
```
docker compose up -d
```
You can then visit http://localhost:8080/swagger/index.html
To stop the service
```
docker compose down
```

## kubernetes

### Enable kubernetes in docker desktop

In `Setting` enable kubernetes.

### apply kubernetes yaml

```
kubectl apply -f ipchecker-deployment.yaml
```

### check kubernetes

```
kubectl get services
```
```
kubectl port-forward svc/ipchecker-service 8080:8080
```

### AI Coding Assistance

In this project, I utilized ChatGPT to assist with specific coding tasks, particularly for generating partial code snippets related to repetitive and non-logical sections. Sometimes, helped me to solve some errors. However, I structured the entire project and thoroughly reviewed all the generated code to ensure it met the project's standards and functionality.
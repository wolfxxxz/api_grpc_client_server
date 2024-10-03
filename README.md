# api_grpc_client_server
# User Management & Authentication gRPC API 
# Key Features:
    gRPC Communication: Using gRPC for efficient and low-latency communication, avoiding the overhead of REST API routes and long paths.
    Protobuf for Data Serialization: Protocol Buffers (proto) help ensure efficient serialization and cross-language support.
    Redis for Caching: Integration with Redis for storing tokens or caching frequently used data to improve performance.
    MongoDB for Scalable User Storage: Using MongoDB for scalable and flexible user data storage.
    Unit Testing and Mocks: Unit tests for gRPC services and mock implementations to ensure proper testing of gRPC methods.

# Packages Used:
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.1
	github.com/joho/godotenv v1.5.1
	github.com/labstack/gommon v0.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.2
	go.mongodb.org/mongo-driver v1.12.1
	golang.org/x/crypto v0.16.0
	google.golang.org/grpc v1.59.0
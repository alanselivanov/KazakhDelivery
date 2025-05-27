# KazakhDelivery

## Project Overview
KazakhDelivery is a microservice-based e-commerce platform designed to facilitate online shopping and delivery services in Kazakhstan. The system is built with a modern microservices architecture, allowing for scalable and maintainable development of different business domains.

## Technologies Used
- **Go** (Golang) - Main programming language
- **gRPC** - Inter-service communication protocol
- **Protocol Buffers** - Data serialization format
- **MongoDB** - Primary database
- **Redis** - Caching solution
- **NATS** - Messaging system for asynchronous communication
- **Gin** - HTTP web framework for API Gateway
- **Docker** - Containerization

## System Architecture
The project consists of four main microservices:

1. **API Gateway** - Entry point for client applications, routes requests to appropriate services
2. **User Service** - Manages user authentication, registration, and profile management
3. **Inventory Service** - Handles product catalog, categories, and stock management
4. **Order Service** - Processes order creation, management, and tracking

## How to Run Locally

### Prerequisites
- Docker and Docker Compose
- Go 1.18 or higher

### Setup Steps

1. Clone the repository:
```bash
git clone https://github.com/yourusername/KazakhDelivery.git
cd KazakhDelivery
```

2. Start the services using Docker Compose:
```bash
docker-compose up -d
```

3. The services will be available at:
   - API Gateway: http://localhost:8080
   - User Service: localhost:50053 (gRPC)
   - Inventory Service: localhost:50051 (gRPC)
   - Order Service: localhost:50052 (gRPC)

## Tests

### Test Structure
The project uses a comprehensive testing strategy with different test types:
- **Unit Tests** - Test individual components in isolation
- **Integration Tests** - Test interactions between multiple services

### Running Tests

#### Run All Tests
```
go test ./tests/...
```

#### Run Only Unit Tests
```
go test ./tests/unit/...
```

#### Run Only Integration Tests
```
go test ./tests/integration/...
```

#### Run Tests with Verbose Output
```
go test ./tests/... -v
```

### Testing Tools
- **testify** - Assertion library
- **httpexpect** - API testing framework
- **testcontainers-go** - Container management for tests
- **gomock** - Mocking framework for unit tests

## gRPC Endpoints

### User Service
- `RegisterUser` - Register a new user
- `AuthenticateUser` - Authenticate user and generate token
- `GetUserProfile` - Get user profile information
- `UpdateUserProfile` - Update user profile information

### Inventory Service
- `CreateProduct` - Create a new product
- `GetProduct` - Get product details
- `UpdateProduct` - Update product information
- `DeleteProduct` - Delete a product
- `ListProducts` - List products with pagination
- `CreateCategory` - Create a new category
- `GetCategory` - Get category details
- `UpdateCategory` - Update category information
- `DeleteCategory` - Delete a category
- `ListCategories` - List all categories
- `DecreaseStock` - Decrease product stock quantity

### Order Service
- `CreateOrder` - Create a new order
- `GetOrder` - Get order details
- `UpdateOrder` - Update order information
- `ListOrders` - List orders for a user
- `CheckStock` - Check if product is in stock

## Implemented Features

- **User Management**
  - User registration and authentication
  - JWT-based authentication
  - User profile management

- **Product Management**
  - Product CRUD operations
  - Category management
  - Stock management

- **Order Processing**
  - Order creation and management
  - Stock verification
  - Order history

- **System Features**
  - Microservice architecture
  - Inter-service communication via gRPC
  - Caching with Redis
  - Asynchronous messaging with NATS
  - MongoDB for data persistence
  - Graceful shutdown handling
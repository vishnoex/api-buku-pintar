# Buku Pintar API

A RESTful API service for Buku Pintar application built with Go, following Clean Architecture principles and SOLID design patterns.

## Project Structure

```
.
├── cmd/
│   └── api/                 # Application entry point
├── internal/
│   ├── domain/             # Enterprise business rules
│   │   ├── entity/        # Business objects
│   │   ├── repository/    # Repository interfaces
│   │   └── service/       # Service interfaces
│   ├── repository/        # Repository implementations
│   │   └── mysql/        # MySQL repository
│   ├── usecase/          # Application business rules
│   └── delivery/         # Interface adapters
│       └── http/         # HTTP handlers
├── pkg/                   # Public packages
│   ├── config/           # Configuration management
│   └── firebase/         # Firebase integration
├── config.json           # Application configuration
├── firebase-credentials.json  # Firebase service account credentials
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- MySQL 8.0
- Firebase project

## Configuration

The application uses a JSON configuration file (`config.json`) for all settings:

```json
{
    "firebase": {
        "credentials_file": "./firebase-credentials.json"
    },
    "database": {
        "host": "mysql-8",
        "port": "3306",
        "user": "root",
        "password": "bukanpassword",
        "name": "bukupintar",
        "params": "parseTime=true"
    },
    "app": {
        "port": "8080"
    }
}
```

## Firebase Setup

1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com/)
2. Enable Authentication in your Firebase project
3. Go to Project Settings > Service Accounts
4. Click "Generate New Private Key" to download your service account credentials
5. Save the downloaded JSON file as `firebase-credentials.json` in the project root

## Running the Application

### Using Docker

1. Build and start the containers:
```bash
docker-compose up --build
```

2. The API will be available at `http://localhost:8080`

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### Public Endpoints

- `POST /users/register` - Register a new user

### Protected Endpoints (Requires Firebase Authentication)

- `GET /users` - Get user profile
- `PUT /users/update` - Update user profile
- `DELETE /users/delete` - Delete user account

## Authentication

The API uses Firebase Authentication. To access protected endpoints:

1. Include the Firebase ID token in the Authorization header:
```
Authorization: Bearer <firebase_id_token>
```

## SOLID Principles Implementation

### Single Responsibility Principle (SRP)
- Each package and component has a single, well-defined responsibility
- Clear separation between domain logic, data access, and delivery mechanisms
- Repository interfaces and implementations are separated
- Use cases handle specific business operations

### Open/Closed Principle (OCP)
- Repository interfaces allow for different implementations without modifying existing code
- Service interfaces enable extending functionality through new implementations
- Middleware system allows adding new authentication methods without changing existing code

### Liskov Substitution Principle (LSP)
- Repository implementations are interchangeable as long as they satisfy the interface
- Service implementations can be swapped without affecting the rest of the system
- HTTP handlers can be replaced with different implementations while maintaining the same contract

### Interface Segregation Principle (ISP)
- Small, focused interfaces for repositories and services
- Separate interfaces for different types of operations
- Clients only depend on the interfaces they use

### Dependency Inversion Principle (DIP)
- High-level modules (use cases) depend on abstractions
- Low-level modules (repositories) implement these abstractions
- Dependencies are injected through constructors
- Easy to swap implementations for testing or different environments

## Clean Architecture

The project follows Clean Architecture principles with distinct layers:

### 1. Domain Layer (Enterprise Business Rules)
- Contains enterprise-wide business rules
- Independent of other layers
- Includes entities, repository interfaces, and service interfaces
- No dependencies on external frameworks

### 2. Use Case Layer (Application Business Rules)
- Implements application-specific business rules
- Orchestrates the flow of data to and from entities
- Depends only on the domain layer
- Contains use case implementations

### 3. Interface Adapters Layer
- Converts data between the format most convenient for use cases and entities
- Contains controllers, presenters, and gateways
- Implements repository interfaces
- Handles HTTP requests and responses

### 4. Frameworks & Drivers Layer
- Contains frameworks and tools like the database, web framework, etc.
- All frameworks are kept in this layer
- Communicates with the interface adapters layer

## Testing

To run tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 
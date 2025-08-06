# Buku Pintar API

A RESTful API service for Buku Pintar application built with Go, following Clean Architecture principles and SOLID design patterns.

## Project Structure

```
.
├── cmd/
│   └── api/                 # Application entry point
├── internal/
│   ├── constant/          # Constant values
│   ├── domain/            # Enterprise business rules
│   │   ├── entity/        # Business objects
│   │   ├── repository/    # Repository interfaces
│   │   └── service/       # Service interfaces
│   ├── repository/        # Repository implementations
│   │   └── mysql/        # MySQL repository
│   ├── usecase/          # Application business rules
│   └── delivery/         # Interface adapters
│       └── http/         # HTTP handlers
│          └── response/  # HTTP response models
├── pkg/                   # Public packages
│   ├── config/           # Configuration management
│   └── firebase/         # Firebase integration
├── migrations/            # Database migrations
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
- Xendit account (for payment processing)

## Configuration

The application uses a JSON configuration file (`config.json`) for all settings:

```json
{
    "firebase": {
        "credentials_file": "./firebase-credentials.json"
    },
    "payment": {
        "xendit": {
            "key": "your_xendit_api_key"
        }
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
        "port": "8080",
        "environment": "local"
    }
}
```

## Firebase Setup

1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com/)
2. Enable Authentication in your Firebase project
3. Go to Project Settings > Service Accounts
4. Click "Generate New Private Key" to download your service account credentials
5. Save the downloaded JSON file as `firebase-credentials.json` in the project root

## Xendit Setup

1. Create a Xendit account at [Xendit Dashboard](https://dashboard.xendit.co/)
2. Get your API key from the Xendit Dashboard
3. Add the API key to your `config.json` file under `payment.xendit.key`
4. Configure webhook URLs in your Xendit dashboard to point to your callback endpoint

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

2. Run database migrations:
```bash
# If using golang-migrate
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/bukupintar?parseTime=true" up
```

3. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### Public Endpoints

- `POST /users/register` - Register a new user
- `POST /payments/callback` - Xendit payment status callback

### Protected Endpoints (Requires Firebase Authentication)

- `GET /users` - Get user profile
- `PUT /users/update` - Update user profile
- `DELETE /users/delete` - Delete user account
- `POST /payments/initiate` - Initiate a new payment

## Payment Integration

The API integrates with Xendit for payment processing. The payment flow works as follows:

### Initiating a Payment

```bash
POST /payments/initiate
Authorization: Bearer <firebase_id_token>
Content-Type: application/json

{
    "user_id": "user123",
    "amount": 50000,
    "currency": "IDR",
    "description": "Premium subscription"
}
```

Response:
```json
{
    "id": "payment-uuid",
    "user_id": "user123",
    "amount": 50000,
    "currency": "IDR",
    "status": "pending",
    "xendit_reference": "inv_123456789",
    "description": "Premium subscription",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

### Payment Status Callback

Xendit will send payment status updates to the callback endpoint:

```bash
POST /payments/callback
Content-Type: application/json

{
    "external_id": "payment-uuid",
    "status": "PAID",
    "invoice_id": "inv_123456789"
}
```

### Payment Statuses

- `pending` - Payment initiated, waiting for completion
- `paid` - Payment completed successfully
- `failed` - Payment failed
- `expired` - Payment expired

## Authentication

The API uses Firebase Authentication. To access protected endpoints:

1. Include the Firebase ID token in the Authorization header:
```
Authorization: Bearer <firebase_id_token>
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    role VARCHAR(20) NOT NULL,
    avatar TEXT,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Payments Table
```sql
CREATE TABLE payments (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL,
    xendit_reference VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## SOLID Principles Implementation

### Single Responsibility Principle (SRP)
- Each package and component has a single, well-defined responsibility
- Clear separation between domain logic, data access, and delivery mechanisms
- Repository interfaces and implementations are separated
- Use cases handle specific business operations
- Payment service handles only payment-related business logic

### Open/Closed Principle (OCP)
- Repository interfaces allow for different implementations without modifying existing code
- Service interfaces enable extending functionality through new implementations
- Middleware system allows adding new authentication methods without changing existing code
- Payment gateway integration can be extended to support other providers

### Liskov Substitution Principle (LSP)
- Repository implementations are interchangeable as long as they satisfy the interface
- Service implementations can be swapped without affecting the rest of the system
- HTTP handlers can be replaced with different implementations while maintaining the same contract
- Payment service implementations can be swapped for different payment gateways

### Interface Segregation Principle (ISP)
- Small, focused interfaces for repositories and services
- Separate interfaces for different types of operations
- Clients only depend on the interfaces they use
- Payment repository and service interfaces are focused on payment operations

### Dependency Inversion Principle (DIP)
- High-level modules (use cases) depend on abstractions
- Low-level modules (repositories) implement these abstractions
- Dependencies are injected through constructors
- Easy to swap implementations for testing or different environments
- Payment use cases depend on payment service abstractions

## Clean Architecture

The project follows Clean Architecture principles with distinct layers:

### 1. Domain Layer (Enterprise Business Rules)
- Contains enterprise-wide business rules
- Independent of other layers
- Includes entities, repository interfaces, and service interfaces
- No dependencies on external frameworks
- Payment entity and interfaces are domain-driven

### 2. Use Case Layer (Application Business Rules)
- Implements application-specific business rules
- Orchestrates the flow of data to and from entities
- Depends only on the domain layer
- Contains use case implementations
- Payment use cases handle payment flow orchestration

### 3. Interface Adapters Layer
- Converts data between the format most convenient for use cases and entities
- Contains controllers, presenters, and gateways
- Implements repository interfaces
- Handles HTTP requests and responses
- Payment handlers adapt HTTP requests to use cases

### 4. Frameworks & Drivers Layer
- Contains frameworks and tools like the database, web framework, etc.
- All frameworks are kept in this layer
- Communicates with the interface adapters layer
- Xendit SDK integration is isolated in this layer

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

## Payment Integration Notes

- The payment integration uses Xendit's invoice API
- Payment amounts are stored in the smallest currency unit (e.g., cents for USD)
- Webhook callbacks are used to update payment status
- The system supports multiple currencies
- Payment status is synchronized between Xendit and the local database

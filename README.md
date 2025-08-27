# Buku Pintar API

A RESTful API service for Buku Pintar application built with Go, following Clean Architecture principles and SOLID design patterns. The API supports both Firebase Authentication and OAuth2 for multiple providers (Google, GitHub, Facebook).

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
│   │   ├── mysql/         # MySQL repository
│   │   └── redis/         # Redis repository (caching)
│   ├── service/           # Service implementations
│   ├── usecase/          # Application business rules
│   └── delivery/         # Interface adapters
│       └── http/         # HTTP handlers
│          └── response/  # HTTP response models
├── pkg/                   # Public packages
│   ├── config/           # Configuration management
│   ├── firebase/         # Firebase integration
│   └── oauth2/           # OAuth2 service
├── migrations/            # Database migrations
├── seeder/               # Database seeders
├── config.json           # Application configuration
├── firebase-credentials.json  # Firebase service account credentials
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- MySQL 8.0
- Firebase project
- Xendit account (for payment processing)
- OAuth2 provider accounts (Google, GitHub, Facebook)

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
    },
    "oauth2": {
        "google": {
            "client_id": "your_google_client_id",
            "client_secret": "your_google_client_secret",
            "redirect_url": "http://localhost:8080/oauth2/google/redirect"
        },
        "github": {
            "client_id": "your_github_client_id",
            "client_secret": "your_github_client_secret",
            "redirect_url": "http://localhost:8080/oauth2/github/redirect"
        },
        "facebook": {
            "client_id": "your_facebook_client_id",
            "client_secret": "your_facebook_client_secret",
            "redirect_url": "http://localhost:8080/oauth2/facebook/redirect"
        }
    }
}
```

## Firebase Setup

1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com/)
2. Enable Authentication in your Firebase project
3. Go to Project Settings > Service Accounts
4. Click "Generate New Private Key" to download your service account credentials
5. Save the downloaded JSON file as `firebase-credentials.json` in the project root

## OAuth2 Setup

### Google OAuth2 Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Go to Credentials > Create Credentials > OAuth 2.0 Client IDs
5. Configure the OAuth consent screen
6. Set application type to "Web application"
7. Add authorized redirect URIs (e.g., `http://localhost:8080/oauth2/google/redirect`)
8. Copy the Client ID and Client Secret to your `config.json`

### GitHub OAuth2 Setup

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the application details:
   - Application name: Your app name
   - Homepage URL: Your app homepage
   - Authorization callback URL: `http://localhost:8080/oauth2/github/redirect`
4. Click "Register application"
5. Copy the Client ID and Client Secret to your `config.json`

### Facebook OAuth2 Setup

1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create a new app or select an existing one
3. Go to Facebook Login > Settings
4. Add your OAuth redirect URI: `http://localhost:8080/oauth2/facebook/redirect`
5. Copy the App ID and App Secret to your `config.json`

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

### Banner Endpoints

- `GET /banners` - List active banners (paginated)
- `GET /banners/active` - List active banners (paginated)
- `GET /banners/view/{id}` - Get banner by ID
- `POST /banners/create` - Create new banner (protected)
- `PUT /banners/edit/{id}` - Update banner (protected)
- `DELETE /banners/delete/{id}` - Delete banner (protected)

### Category Endpoints

- `GET /categories` - List active categories (paginated)
- `GET /categories/all` - List all categories (paginated)
- `GET /categories/view/{id}` - Get category by ID
- `GET /categories/parent/{parentID}` - List categories by parent (paginated)
- `POST /categories/create` - Create new category (protected)
- `PUT /categories/edit/{id}` - Update category (protected)
- `DELETE /categories/delete/{id}` - Delete category (protected)

### OAuth2 Endpoints

- `POST /oauth2/login` - Initiate OAuth2 login flow
- `POST /oauth2/callback` - Handle OAuth2 callback
- `GET /oauth2/providers` - Get available OAuth2 providers
- `GET /oauth2/{provider}/redirect` - Handle OAuth2 provider redirect

### Protected Endpoints (Requires Firebase or OAuth2 Authentication)

- `GET /users` - Get user profile
- `PUT /users/update` - Update user profile
- `DELETE /users/delete` - Delete user account
- `POST /payments/initiate` - Initiate a new payment

## OAuth2 Authentication Flow

### 1. Initiate OAuth2 Login

```bash
POST /oauth2/login
Content-Type: application/json

{
    "provider": "google"
}
```

Response:
```json
{
    "auth_url": "https://accounts.google.com/oauth2/authorize?...",
    "state": "random_state_string"
}
```

### 2. User Authorization

The user is redirected to the OAuth2 provider's authorization page where they grant permissions.

### 3. OAuth2 Callback

The provider redirects back to your application with an authorization code:

```
GET /oauth2/google/redirect?code=AUTHORIZATION_CODE&state=STATE
```

### 4. Token Exchange

The application exchanges the authorization code for an access token and retrieves user information.

### 5. User Registration/Login

If the user doesn't exist, they are automatically registered. If they exist, they are logged in.

## Payment Integration

The API integrates with Xendit for payment processing. The payment flow works as follows:

### Initiating a Payment

```bash
POST /payments/initiate
Authorization: Bearer <firebase_id_token_or_oauth2_token>
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

The API supports multiple authentication methods:

### Firebase Authentication

Include the Firebase ID token in the Authorization header:
```
Authorization: Bearer <firebase_id_token>
```

### OAuth2 Authentication

Include the OAuth2 access token in the Authorization header:
```
Authorization: Bearer <oauth2_access_token>
```

The middleware automatically detects the token type and validates it accordingly.

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
- OAuth2 service handles only OAuth2 authentication flows

### Open/Closed Principle (OCP)
- Repository interfaces allow for different implementations without modifying existing code
- Service interfaces enable extending functionality through new implementations
- Middleware system allows adding new authentication methods without changing existing code
- Payment gateway integration can be extended to support other providers
- OAuth2 service can be extended to support additional providers

### Liskov Substitution Principle (LSP)
- Repository implementations are interchangeable as long as they satisfy the interface
- Service implementations can be swapped without affecting the rest of the system
- HTTP handlers can be replaced with different implementations while maintaining the same contract
- Payment service implementations can be swapped for different payment gateways
- OAuth2 providers can be swapped without affecting the authentication flow

### Interface Segregation Principle (ISP)
- Small, focused interfaces for repositories and services
- Separate interfaces for different types of operations
- Clients only depend on the interfaces they use
- Payment repository and service interfaces are focused on payment operations
- OAuth2 service interface is focused on OAuth2 operations

### Dependency Inversion Principle (DIP)
- High-level modules (use cases) depend on abstractions
- Low-level modules (repositories) implement these abstractions
- Dependencies are injected through constructors
- Easy to swap implementations for testing or different environments
- Payment use cases depend on payment service abstractions
- OAuth2 handlers depend on OAuth2 service abstractions

## Architecture Overview

The project follows **Clean Architecture** principles with a **UseCase → Service → Repository** pattern, providing clear separation of concerns, improved maintainability, and enhanced performance through intelligent caching strategies.

### Architecture Pattern

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  UseCase        │    │   Service       │
│   (Handlers)    │───▶│  (Orchestration)│───▶│   (Business     │
│                 │    │                 │    │    Logic)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   Repository    │
                                              │   (Data Access) │
                                              │   MySQL + Redis │
                                              └─────────────────┘
```

### Layer Responsibilities

#### 1. **HTTP Layer (Delivery)**
- **Purpose**: Handle HTTP requests and responses
- **Responsibilities**: 
  - Request validation and parsing
  - Response formatting
  - Error handling
  - Authentication middleware
- **Components**: HTTP handlers, middleware, response models

#### 2. **UseCase Layer (Application Business Rules)**
- **Purpose**: Orchestrate application workflows
- **Responsibilities**:
  - Coordinate between handlers and services
  - Data transformation (entities ↔ DTOs)
  - Application-level error handling
  - Business flow orchestration
- **Components**: Use case implementations

#### 3. **Service Layer (Domain Business Logic)**
- **Purpose**: Implement business logic and caching strategies
- **Responsibilities**:
  - Business rule enforcement
  - Cache management (Redis)
  - External service integration
  - Data validation and processing
- **Components**: Service interfaces and implementations

#### 4. **Repository Layer (Data Access)**
- **Purpose**: Abstract data access operations
- **Responsibilities**:
  - Database operations (MySQL)
  - Cache operations (Redis)
  - Data persistence
  - Query optimization
- **Components**: Repository interfaces and implementations

### Clean Architecture Implementation

#### **Domain Layer (Enterprise Business Rules)**
- **Entities**: Core business objects (User, Category, Banner, Payment, etc.)
- **Repository Interfaces**: Data access contracts
- **Service Interfaces**: Business logic contracts
- **Characteristics**: Framework-independent, no external dependencies

#### **Use Case Layer (Application Business Rules)**
- **Orchestration**: Coordinates between layers
- **Data Transformation**: Converts between domain entities and DTOs
- **Error Handling**: Application-level error management
- **Dependencies**: Only depends on domain layer

#### **Interface Adapters Layer**
- **HTTP Handlers**: Request/response handling
- **Repository Implementations**: MySQL and Redis implementations
- **Service Implementations**: Business logic implementations
- **Middleware**: Authentication, logging, etc.

#### **Frameworks & Drivers Layer**
- **Database**: MySQL, Redis
- **External Services**: Firebase, Xendit, OAuth2 providers
- **Web Framework**: Standard HTTP library
- **Configuration**: JSON-based configuration

### Caching Strategy

The application implements a sophisticated caching strategy using Redis:

#### **Cache Patterns**
- **Cache-Aside**: Read from cache first, fallback to database
- **Write-Through**: Update database first, then invalidate cache
- **Cache Invalidation**: Bulk invalidation for data consistency

#### **Cache Keys Structure**
```
banner:list:{limit}:{offset}          # Banner lists
banner:count:total                    # Banner counts
category:list:{limit}:{offset}        # Category lists
category:active:list:{limit}:{offset} # Active category lists
category:count:total                  # Category counts
category:id:{id}                      # Individual categories
```

#### **Cache TTL**
- **Banner Cache**: 5 minutes
- **Category Cache**: 10 minutes
- **Configurable**: TTL can be adjusted per service

#### **Cache Invalidation**
- **Automatic**: Cache cleared on data changes (Create/Update/Delete)
- **Graceful Degradation**: System works without cache
- **Logging**: Cache operations logged for monitoring

### Performance Optimizations

#### **1. Intelligent Caching**
- **Frequently Accessed Data**: Lists, counts, individual items
- **Selective Caching**: Different strategies for different data types
- **Cache Warming**: Pre-populate frequently accessed data

#### **2. Database Optimization**
- **Connection Pooling**: Efficient database connections
- **Query Optimization**: Optimized SQL queries
- **Indexing**: Proper database indexing

#### **3. External Service Integration**
- **Async Operations**: Non-blocking external service calls
- **Circuit Breaker**: Graceful handling of service failures
- **Retry Logic**: Automatic retry for transient failures

### Module-Specific Architecture

#### **Banner Module**
```
UseCase → BannerService → BannerRepository + BannerRedisRepository
```
- **Features**: CRUD operations, active/inactive filtering
- **Caching**: Lists, counts, automatic invalidation
- **Endpoints**: 6 REST endpoints with pagination

#### **Category Module**
```
UseCase → CategoryService → CategoryRepository + CategoryRedisRepository
```
- **Features**: CRUD operations, hierarchical support, parent-child relationships
- **Caching**: Lists, counts, individual items, hierarchical data
- **Endpoints**: 7 REST endpoints with pagination and hierarchy support

#### **User Module**
```
UseCase → UserService → UserRepository
```
- **Features**: Authentication, registration, profile management
- **Integration**: Firebase Auth, OAuth2 providers
- **Security**: Password hashing, token validation

#### **Payment Module**
```
UseCase → PaymentService → PaymentRepository
```
- **Features**: Payment processing, status tracking
- **Integration**: Xendit payment gateway
- **Webhooks**: Payment status callbacks

### Benefits of This Architecture

#### **1. Separation of Concerns**
- **Clear Boundaries**: Each layer has specific responsibilities
- **Maintainability**: Easy to modify individual components
- **Testability**: Each layer can be tested independently

#### **2. Scalability**
- **Horizontal Scaling**: Services can be scaled independently
- **Performance**: Caching reduces database load
- **Flexibility**: Easy to add new features or modify existing ones

#### **3. Reliability**
- **Graceful Degradation**: System works even when components fail
- **Error Handling**: Comprehensive error handling at each layer
- **Monitoring**: Logging and metrics for system health

#### **4. Maintainability**
- **Code Organization**: Clear structure and naming conventions
- **Dependency Management**: Clear dependency flow
- **Documentation**: Comprehensive documentation for each layer

### Future Enhancements

#### **1. Advanced Caching**
- **Cache Warming**: Pre-populate cache with frequently accessed data
- **Cache Tags**: More granular cache invalidation
- **Cache Compression**: Reduce memory usage

#### **2. Performance Monitoring**
- **Metrics**: Prometheus metrics for performance monitoring
- **Tracing**: Distributed tracing for request flows
- **Profiling**: Application performance profiling

#### **3. Advanced Features**
- **Background Jobs**: Async processing for heavy operations
- **Event Sourcing**: Event-driven architecture
- **Microservices**: Service decomposition for large scale

This architecture provides a solid foundation for building scalable, maintainable, and high-performance applications while following clean architecture principles and industry best practices.

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

## OAuth2 Integration Notes

- The OAuth2 integration supports Google, GitHub, and Facebook providers
- OAuth2 tokens are validated with the respective providers
- User information is automatically synchronized between OAuth2 providers and the local database
- The system supports both Firebase and OAuth2 authentication simultaneously
- OAuth2 redirect URLs must be configured in both the application and provider settings
- State parameters are used to prevent CSRF attacks
- Access tokens are returned to the client for subsequent API calls

## Payment Integration Notes

- The payment integration uses Xendit's invoice API
- Payment amounts are stored in the smallest currency unit (e.g., cents for USD)
- Webhook callbacks are used to update payment status
- The system supports multiple currencies
- Payment status is synchronized between Xendit and the local database

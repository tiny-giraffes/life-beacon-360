## Directory Structure

### `cmd/`

- **Purpose**: Entry point for the application.
- **Files**:
  - `main.go`: Initializes the application, loads configuration, connects to the database, and starts the Fiber server.

### `config/`

- **Purpose**: Handles application configuration and environment variables.
- **Files**:
  - `config.go`: Loads environment variables using Viper.
  - `config.yaml` or `.env`: Optional configuration file to store environment-specific settings.

### `internal/`

- **Purpose**: Contains the core application logic, divided into modular packages.

  #### `internal/api/`

  - **Purpose**: Sets up routing and defines HTTP handlers.
  - **Files**:
    - `handlers/`: Contains handler functions for different routes.
      - `location.go`: Handler for location-related endpoints.
      - `auth.go`: Optional handler for authentication-related endpoints.
    - `router.go`: Sets up API routes and applies middleware.

  #### `internal/models/`

  - **Purpose**: Defines data models that map to database tables.
  - **Files**:
    - `location.go`: Model representing location data.
    - `user.go`: Optional user model if authentication is required.

  #### `internal/repository/`

  - **Purpose**: Database interaction layer, providing functions to query and modify data.
  - **Files**:
    - `location_repo.go`: Repository functions for location data.
    - `user_repo.go`: Optional repository for user data.

  #### `internal/services/`

  - **Purpose**: Business logic layer, where core logic is implemented.
  - **Files**:
    - `location_service.go`: Business logic for processing location data.
    - `auth_service.go`: Optional authentication logic.

  #### `internal/utils/`

  - **Purpose**: Utility functions, helpers for response formatting, and error handling.
  - **Files**:
    - `response.go`: Contains helper functions for sending JSON responses and handling errors.

### `pkg/`

- **Purpose**: External or third-party packages that are specific to this application.

  #### `pkg/database/`

  - **Purpose**: Manages the database connection and setup.
  - **Files**:
    - `postgres.go`: Connects to the PostgreSQL database using GORM.

### `Dockerfile`

- **Purpose**: Docker configuration file for building and running the Go Fiber app in a container.

### `docker-compose.yml`

- **Purpose**: Docker Compose file to manage the multi-container setup (e.g., app server, database).

### `go.mod`

- **Purpose**: Go module file, listing dependencies and managing versions.



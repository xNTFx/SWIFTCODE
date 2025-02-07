# Project Setup Instructions

## Docker Setup

### Prerequisites

Ensure you have Docker and Docker Compose installed by running:

```sh
docker --version
docker compose version
```

### Clone the Repository

```sh
git clone <repository_url>
cd <repository_name>
```

### Navigate to the Backend Directory

```sh
cd backend
```

Start Docker Desktop if it is not already running.

### Build and Start the Application in Detached Mode

```sh
docker compose up --build -d
```

### Access the Running Container

```sh
docker exec -it backend_server_ps sh
```

### Run Tests in Docker

#### Run Unit Tests

```sh
go test ./tests/unit_tests/... -v
```

#### Run Integration Tests

```sh
go test ./tests/integration_tests/... -v
```

#### Run All Tests

```sh
go test -v ./tests/...
```

---

## Local Setup

### Prerequisites

Ensure you have **PostgreSQL** and **Go** installed.

### Clone the Repository

```sh
git clone <repository_url>
cd <repository_name>
```

### Install PostgreSQL

Ensure PostgreSQL is installed and running. The database should be named:

```
swift_codes
```

The default connection string is:

```
postgres://postgres:admin@localhost:5432/swift_codes?sslmode=disable
```

The default credentials are:

- **Username:** `postgres`
- **Password:** `admin`

### Changing Database Credentials

To change the database credentials, modify the `.env` file in the `backend` directory. The default `.env` file structure is:

```
POSTGRES_URL=postgres://postgres:admin@localhost:5432/swift_codes?sslmode=disable
SERVER_PORT=8080
ALLOWED_ORIGINS=*
```

Ensure PostgreSQL is installed and running. The database initialization script is located at:

```sh
backend/internal/db/init/database-custom.sql
```

### Install Go

Ensure Go is installed and available in your system's PATH.

### Install Dependencies

Navigate to the backend directory and install the required Go libraries:

```sh
cd backend
go mod tidy
```

### Run the Application

Once the application is running, the backend is available at:

```
http://localhost:8080
```

Start the application:

```sh
go run cmd/server/main.go
```

### Run Tests Locally

#### Run Unit Tests

```sh
go test ./tests/unit_tests/... -v
```

#### Run Integration Tests

```sh
go test ./tests/integration_tests/... -v
```

#### Run All Tests

```sh
go test -v ./tests/...
```
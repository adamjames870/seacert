# SeaCert

SeaCert is a certificate management system designed for maritime certifications. It consists of a Go backend and a React frontend.

You can always access the test server here. Because it is a test server, it may take some time to load (sleeping), and the datatbase may be reset at any time.
- API: https://seacert-api.onrender.com
- Frontend: https://seacert.onrender.com

## Features

- **Certificate Management**: Comprehensive tracking of maritime certifications, including issue dates, certificate numbers, and remarks.
- **Automatic Expiry Tracking**: Intelligent calculation of certificate expiry based on predefined validity periods, with support for manual overrides.
- **Certificate Succession**: Support for superseding existing certificates (e.g., when a certificate is renewed or replaced), maintaining a clear history of certification.
- **Categorization**: Organize certifications by type (e.g., STCW references) and track authorized issuers.
- **Secure Authentication**: Integration with Supabase for robust user authentication and access control.
- **Admin Dashboard**: Specialized endpoints for health monitoring, database statistics, and user management.
- **Type-Safe Backend**: High-performance API built with Go and `sqlc` for reliable data persistence.
- **Modern User Interface**: Responsive React frontend for managing certificates on any device.

## Frontend

The SeaCert frontend is a modern Single Page Application (SPA) built with:

- **React 19**: A powerful library for building user interfaces.
- **TypeScript**: Ensuring type safety and better developer experience.
- **Vite**: A fast build tool and development server.

### Key Frontend Features

- **Certificate Dashboard**: A centralized view of all certificates with sorting capabilities by type, number, issuer, and date.
- **Responsive Design**: Optimized for various screen sizes, ensuring accessibility on mobile and desktop.
- **Real-time Updates**: Manual refresh functionality to ensure the latest data is always visible.
- **Formatted Data Display**: Automatically formats dates and types for better readability.
- **Authentication Flow**: Full integration with Supabase for secure login and account management.

## Project Structure

- `cmd/`: Entry points for the application.
  - `api/`: The main Go backend API server.
  - `get_token/`: Utility for obtaining authentication tokens.
- `internal/`: Internal Go packages.
  - `api/`: HTTP routing and handlers.
  - `domain/`: Core business logic and domain models (e.g., certificates, certificate types).
  - `database/`: Database interaction layer, utilizing `sqlc` for type-safe SQL.
  - `dto/`: Data Transfer Objects for API communication.
- `frontend/`: React frontend built with TypeScript and Vite.
- `sqlc.yaml`: Configuration for `sqlc` code generation.

## API Documentation

The SeaCert API provides endpoints for both administrative tasks and certificate management. All endpoints (except `/admin/healthz`) require authentication via a Bearer token.

### Admin Endpoints

#### Health Check
`GET /admin/healthz`
- **Description**: Verifies if the API server is running.
- **Response**: `OK` (text/plain)

#### Reset Database
`POST /admin/reset`
- **Description**: Resets all tables in the database. Only available when `PLATFORM=dev`.
- **Response**:
  ```json
  { "message": "db reset" }
  ```

#### Database Statistics
`GET /admin/dbstats`
- **Description**: Returns counts of certificates, types, issuers, and users.
- **Response**:
  ```json
  {
    "count-certs": 10,
    "count-cert-types": 5,
    "count-issuers": 3,
    "count-users": 1,
    "user-id": "uuid",
    "user-email": "user@example.com"
  }
  ```

#### User Profile
`GET /admin/users`
- **Description**: Retrieves the profile of the authenticated user.
- **Response**:
  ```json
  {
    "id": "uuid",
    "created-at": "timestamp",
    "updated-at": "timestamp",
    "forename": "John",
    "surname": "Doe",
    "email": "john.doe@example.com",
    "nationality": "British",
    "role": "user"
  }
  ```

#### Update User Profile
`PUT /admin/users`
- **Description**: Updates the profile of the authenticated user.
- **Body**:
  ```json
  {
    "forename": "John",
    "surname": "Doe",
    "nationality": "British"
  }
  ```
- **Response**: Updated user object.

### Certificate API

#### List Certificates
`GET /api/certificates`
- **Description**: Retrieve all certificates for the authenticated user.
- **Parameters**: `id` (optional query param) to retrieve a single certificate.
- **Response**: Array of certificate objects (or a single object if `id` is provided).
  ```json
  [
    {
      "id": "uuid",
      "created-at": "timestamp",
      "cert-type-name": "STCW II/2",
      "cert-number": "12345",
      "issuer-name": "MCA",
      "issued-date": "2023-01-01T00:00:00Z",
      "expiry-date": "2028-01-01T00:00:00Z",
      "remarks": "...",
      "has-successors": false,
      "predecessors": [
        {
          "certificate": { "id": "old-uuid", "cert-number": "..." },
          "reason": "updated"
        }
      ]
    }
  ]
  ```

#### Add Certificate
`POST /api/certificates`
- **Body**:
  ```json
  {
    "cert-type-id": "uuid",
    "cert-number": "12345",
    "issuer-id": "uuid",
    "issued-date": "2023-01-01",
    "alternative-name": "Optional Name",
    "remarks": "Optional Remarks",
    "manual-expiry": "2028-01-01",
    "supersedes": "uuid-of-old-cert",
    "supersede-reason": "updated"
  }
  ```
- **Response**: Created certificate object.

#### Update Certificate
`PUT /api/certificates`
- **Body**:
  ```json
  {
    "id": "uuid",
    "cert-number": "54321",
    "cert-type-id": "uuid",
    "issuer-id": "uuid",
    "issued-date": "2023-02-01",
    "alternative-name": "New Name",
    "remarks": "Updated remarks",
    "manual-expiry": "2028-02-01"
  }
  ```
- **Response**: Updated certificate object.

#### Certificate Types
`GET /api/cert-types` | `POST /api/cert-types`
- **GET Response**: Array of certificate type objects.
- **POST Body**:
  ```json
  {
    "name": "Full Name",
    "short-name": "Short Code",
    "stcw-reference": "A-VI/1",
    "normal-validity-months": 60
  }
  ```

#### Issuers
`GET /api/issuers` | `POST /api/issuers`
- **GET Response**: Array of issuer objects.
- **POST Body**:
  ```json
  {
    "name": "Issuer Name",
    "country": "Country",
    "website": "https://..."
  }
  ```

## Prerequisites

- [Go](https://go.dev/) (version 1.25.4 or later)
- [Node.js](https://nodejs.org/) and npm (for the frontend)
- [PostgreSQL](https://www.postgresql.org/) (or a Supabase instance)

## Getting Started

### Backend

1. Navigate to the root directory.
2. Create a `.env` file based on the required environment variables:
   ```env
   DB_URL=postgres://user:password@host:port/dbname
   PLATFORM=dev # or production
   PORT=8080
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the API server:
   ```bash
   go run ./cmd/api
   ```

### Frontend

1. Navigate to the `frontend/` directory.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the development server:
   ```bash
   npm run dev
   ```

## Database

This project uses `sqlc` to generate Go code from SQL queries.
- SQL schemas and queries are located in `internal/database/queries/`.
- To regenerate the code after modifying SQL files, run:
  ```bash
  sqlc generate
  ```

## Authentication

The project integrates with Supabase for authentication. Check `supainfo.txt` (if available in your local environment) for project-specific Supabase configuration details.

## License

This project is licensed under the **Business Source License 1.1** (BSL 1.1).

- **Licensor**: Adam James
- **Change Date**: 2030-01-01
- **Change License**: Apache License, Version 2.0

Under this license, you are free to use the software for any non-production purposes. For production use, a license is required unless the deployment is for personal, non-commercial use by an individual. On the **Change Date**, the license automatically converts to the permissive **Apache License 2.0**.

See the [LICENSE](LICENSE) file for the full text.

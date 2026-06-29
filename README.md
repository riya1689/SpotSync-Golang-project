# SpotSync Backend API

SpotSync is a robust backend system for managing a smart car parking ecosystem. Built with Go, it provides a seamless API for users to find and reserve parking spots, and for administrators to manage parking zones and monitor reservations. The project ensures data integrity under high concurrency by utilizing row-level database locks.

**Live Demo URL:** [https://spotsync-production-a3e4.up.railway.app/](https://spotsync-production-a3e4.up.railway.app/)

---

## Features

- **User Authentication:** Secure registration and login using JWT (JSON Web Tokens).
- **Role-Based Access Control:** Separate roles for standard `driver` users and `admin` users.
- **Zone Management:** Admins can create and manage parking zones (e.g., general, EV charging, covered).
- **Reservation System:** Users can book parking spots and view their active reservations.
- **Concurrency Safety:** Utilizes PostgreSQL row-level locks (`SELECT ... FOR UPDATE`) inside transactions to prevent double-booking or overbooking of parking spots.
- **Dockerized:** Fully containerized for easy deployment on modern cloud platforms like Railway.

---

## Tech Stack

- **Language:** Go (1.25+)
- **Web Framework:** [Echo v4](https://echo.labstack.com/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL
- **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Validation:** Go Playground Validator (`validator/v10`)
- **Hot Reloading:** [Air](https://github.com/air-verse/air)

---

## Architecture (Layered Design)

SpotSync follows a clean, highly cohesive **Layered Architecture** pattern. This separates concerns, making the codebase maintainable, testable, and scalable.

```mermaid
graph TD
    A[Client Request] -->|HTTP/JSON| B(Handler Layer)
    B -->|DTOs| C(Service Layer)
    C -->|Business Models| D(Repository Layer)
    D -->|SQL/GORM| E[(PostgreSQL Database)]
```

### Layer Breakdown:
1. **Handler Layer (`/internal/handler`):** The entry point for API requests. It parses HTTP requests, validates JSON payloads (using DTOs), calls the appropriate service, and formats the standard JSON response.
2. **Service Layer (`/internal/service`):** Contains the core business logic. It applies rules (e.g., checking if a user already exists, hashing passwords) and delegates data retrieval/storage to the repository.
3. **Repository Layer (`/internal/repository`):** Handles all direct database interactions. It abstracts away GORM and SQL queries, providing a clean interface for the Service layer to request data (e.g., `ReserveSpot` with row locks).
4. **Models & DTOs (`/internal/models` & `/internal/dto`):** Models represent database tables, while DTOs (Data Transfer Objects) define the exact shape of incoming requests and outgoing responses.

---

## ⚙️ Setup & Installation

### Prerequisites
- Go 1.25 or higher
- PostgreSQL database
- Git

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/spotsync-golang-project.git
cd spotsync-golang-project
```

### 2. Environment Variables
Create a `.env` file in the root directory and add the following required variables:
```env
# Database connection string
DATABASE_URL="postgresql://user:password@host:port/dbname?sslmode=require"

# Secret key for signing JWT tokens
JWT_SECRET="your_secret_key"

# Port for the server to run on (Default: 8080)
PORT=8080
```

### 3. Run Locally (with Air)
If you have [Air](https://github.com/air-verse/air) installed for hot-reloading:
```bash
air
```
Alternatively, using standard Go commands:
```bash
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.

---

## Docker Setup

A `Dockerfile` is included for containerized deployment. It uses a multi-stage build to keep the final image incredibly lightweight (Alpine-based).

**Build the image:**
```bash
docker build -t spotSync-golang-Project .
```

**Run the container:**
```bash
docker run -p 8080:8080 --env-file .env spotSync-golang-Project
```

---

## 📡 API Endpoints

### Public Routes
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Health check & Database status |
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive JWT |
| `GET` | `/api/v1/zones` | Get all parking zones |
| `GET` | `/api/v1/zones/:id` | Get details of a specific zone |

### Protected Routes (Requires JWT Bearer Token)
| Method | Endpoint | Role | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/reservations` | `driver`, `admin` | Reserve a parking spot |
| `GET` | `/api/v1/reservations/my-reservations` | `driver`, `admin` | Get logged-in user's reservations |
| `DELETE` | `/api/v1/reservations/:id` | `driver`, `admin` | Cancel a reservation |

### Admin Routes (Requires Admin JWT Bearer Token)
| Method | Endpoint | Role | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/zones` | `admin` | Create a new parking zone |
| `GET` | `/api/v1/reservations` | `admin` | View all system reservations |



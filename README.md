# 🏥 VatanSoft Hospital Management Backend

A backend system for managing hospitals, users, departments, and appointments – built with Go, Gin, GORM, Redis, and PostgreSQL.

---

## 📄 Setup

```bash
docker-compose up --build
```

Then open your browser and visit:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 📁 Project Structure

- `/cmd` – Application entry point (`main.go`)
- `/controllers` – All HTTP handlers (auth, hospital, user, etc.)
- `/models` – GORM models for DB schema
- `/config` – DB and Redis setup
- `/docs` – Swagger-generated files
- `/middlewares` – Auth and other middleware
- `/utils` – Utility functions (e.g., password hashing)

---

## 🔧 Technologies

- **Go** (Golang)
- **Gin** (HTTP Web Framework)
- **GORM** (ORM for Go)
- **PostgreSQL** (Relational Database)
- **Redis** (Caching & temporary storage)
- **Docker** + **Docker Compose**
- **Swagger** (API Documentation)

---

## ✅ Features

- JWT-based authentication
- Admin and staff role system
- Hospital and user management
- Department and doctor linking
- Password reset with Redis
- Swagger UI for live API docs

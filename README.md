# Clean Mobile App (Golang)

A clean architecture template for a mobile application backend using **Go 1.23**. This repository provides a structured approach to build scalable and maintainable backend systems.

---

## 📂 Project Structure

```
clean_mobile_app/
├── cmd/
│   └── web/
│       ├── helpers.go         # Utility functions for the web layer
│       ├── initializer.go     # Application initialization logic
│       ├── main.go            # Entry point of the application
│       ├── middleware.go      # HTTP middlewares
│       └── routes.go          # HTTP route definitions
├── config/
│   └── config.yaml            # YAML file for application configuration
├── db/
│   └── migrations/            # Database migration files
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration loading and handling
│   ├── handlers/              # HTTP request handlers
│   ├── models/                # Application data models
│   ├── repositories/          # Data access layer (DB interactions)
│   └── services/              # Business logic and services
├── go.mod                     # Go module definition
```

---

## 🚀 Features

- **Clean Architecture**: Well-structured layers for separation of concerns.
- **Go 1.23**: Latest version of Go with performance improvements.
- **Scalability**: Designed to handle complex business logic with ease.
- **Config Management**: Centralized and manageable application settings.
- **Database Migrations**: Organized schema changes under the `db/migrations/` directory.

---

## 🛠️ Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ZhumabekovAlim/clean_mobile_app.git
   cd clean_mobile_app
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment**:
   - Update `config/config.yaml` with your application settings.

4. **Run the application**:
   ```bash
   go run cmd/web/main.go
   ```

---

## 🔧 Configuration

The configuration file (`config/config.yaml`) includes parameters like database credentials, server ports, etc. Modify it to suit your environment.

---

## 🗃️ Database Migrations

Use a tool like [golang-migrate](https://github.com/golang-migrate/migrate) to manage your database migrations. Place your `.sql` files in the `db/migrations/` folder.

To run migrations:
```bash
migrate -path db/migrations -database "your-database-url" up
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Open a pull request.


---

## 👨‍💻 Author

- **Alim Zhumabekov** - [ZhumabekovAlim](https://github.com/ZhumabekovAlim)


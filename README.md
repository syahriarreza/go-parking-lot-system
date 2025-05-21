# Go Parking Lot System 🚗

This project is a multi-level parking lot simulator implemented in Go (Golang), featuring a RESTful API, thread-safe access for concurrent gates, dynamic layout configuration, and a built-in HTML UI for simulation.

---

## ✨ Features

- Multi-floor support with customizable layout
- Spot types: B-1 (Bicycle), M-1 (Motorcycle), A-1 (Automobile), X-0 (Inactive)
- Thread-safe parking/unparking mechanism
- Configure layout using JSON or Excel upload
- View floor layout as a readable grid
- Web-based UI for interactive testing
- System reset functionality

---

## 📁 Project Structure

```
go-parking-lot-system/
├── cmd/                 # Main entrypoint
├── config/              # YAML configuration file
├── internal/
│   ├── app/             # Parking logic
│   ├── handler/         # HTTP API handlers
│   ├── service/         # In-memory storage and operations
│   └── utils/           # Excel parser and helper functions
├── static/              # Frontend HTML interface
├── test/                # API test suite
├── Makefile             # Testing and coverage shortcuts
└── README.md
```

---

## 🚀 Getting Started

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Create a config file

Example `config/app.yaml`:

```yaml
server:
  port: 8080
parking:
  floors: 3
  rows: 5
  columns: 5
```

### 3. Run the server

```bash
go run ./cmd/main.go
```

Access the web UI at: [http://localhost:8080](http://localhost:8080)

---

## 🌐 REST API Endpoints

| Method | Endpoint                          | Description                        |
|--------|-----------------------------------|------------------------------------|
| POST   | `/api/v1/park`                    | Assign a spot to a vehicle         |
| POST   | `/api/v1/unpark`                  | Unpark a vehicle                   |
| GET    | `/api/v1/available?vehicleType=M` | List available spots by type       |
| GET    | `/api/v1/search/:vehicleNumber`   | Lookup last known spot of vehicle  |
| POST   | `/api/v1/layout/floor`            | Set layout for a floor (JSON)      |
| POST   | `/api/v1/layout/upload`           | Upload layout using Excel file     |
| GET    | `/api/v1/layout/map?floor=1`      | Get human-readable floor layout    |
| POST   | `/api/v1/reset`                   | Clear all layout and parking state |

---

## 📥 Postman Collection

You can also test all available API endpoints using Postman.

A full collection is included in:

```
docs/Go Parking Lot API.postman_collection.json
```

To use:
1. Open Postman
2. Click "Import"
3. Select the `.json` file above
4. Test endpoints like `POST /park`, `GET /available`, etc.

---

## 🧪 Testing & Coverage

### Run tests

```bash
make test
```

### Run with coverage + HTML report

```bash
make test-coverage
```

---

## 💻 HTML Simulator

Accessible at: [http://localhost:8080](http://localhost:8080)

Supported actions:
- ✅ Park / Unpark vehicle
- 🔍 Search vehicle
- 📄 Upload Excel layout
- ✍️ Set layout manually per floor (JSON)
- 🗺️ View layout grid
- 🔁 Reset system state

---

## ⚙️ Dependencies

- Go 1.20+
- Echo framework (HTTP server)
- testify (unit testing)
- viper (config reader)
- openpyxl (Excel layout support)

---

## 👤 Author

Built by [@syahriarreza](https://github.com/syahriarreza) with ❤️

This project can be extended to include database persistence, authentication, and advanced web UI.

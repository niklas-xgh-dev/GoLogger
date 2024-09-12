<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)

</div>

# GoLogger - System Log Collection Tool

GoLogger is a Go-based system for efficient collection and storage of system logs and metrics.

## Technical Overview

- **Language**: Go 1.20+
- **Storage**: PostgreSQL
- **OS Compatibility**: Linux (Ubuntu), macOS

## Key Features

1. **Comprehensive Log Collection**
   - System metrics (CPU, memory, disk utilization)
   - Network connections and open ports
   - Running services and processes
   - User activities (logins, system users)

2. **Efficient Data Processing**
   - Concurrent log collection and processing
   - Structured log data in JSON format

3. **Modular Architecture**
   - Easily extendable for additional log sources
   - Pluggable storage backend (currently PostgreSQL)

## Main Script Functionality

The core of GoLogger is built around several key functions:

- `collectSecurityLogs()`: Aggregates data from various system sources
- `getOpenPorts()`: Scans for open network ports
- `getRunningServices()`: Lists active system services
- `getSystemUsers()`: Retrieves system user accounts
- `getRecentLogins()`: Fetches recent user login information
- `getSuspiciousProcesses()`: Identifies potentially suspicious running processes

The main loop collects logs at regular intervals, formats them as JSON, and stores them in the configured database.

## Quick Start

1. Clone the repository
2. Ensure Go 1.20+ is installed
3. Run `go mod tidy`
4. Configure PostgreSQL in `config/.env`
5. Build and run: `go build && ./GoLogger`

## Configuration

Set your database URL in `config/.env`:

```
DATABASE_URL=postgres://u:p@localhost:5432/dbname?sslmode=disable
```

## Contributing

Contributions are welcome. Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.
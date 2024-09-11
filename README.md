<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![Elasticsearch](https://img.shields.io/badge/Elasticsearch-005571?style=for-the-badge&logo=elasticsearch&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)

</div>

# GoLogger - Comprehensive Log Collection System

GoLogger is a Go-based log collection system designed for gathering and storage of system logs.

## Recent Refactoring

- Migrated from Python to Go for enhanced performance and concurrency
- Streamlined codebase for improved maintainability and scalability
- Current storage uses PostgreSQL, with planned migration or integration with Elasticsearch
- Extended compatibility to Ubuntu OS
- Maintained the extensible log collection architecture

## Technical Stack

- Go 1.20+ for core functionality (refactored from Python)
- PostgreSQL for current log storage and querying
- Planned Elasticsearch integration for advanced log analytics
- Ubuntu OS compatibility

## Key Features

1. **Extensible Log Collection**
   - System metric collection (CPU, memory, disk utilization)
   - Modular design for integration of additional log sources

2. **High-Performance Data Processing**
   - Optimized log parsing and structuring
   - Concurrent processing for improved throughput

3. **Scalable Storage**
   - PostgreSQL utilized for ACID-compliant log persistence
   - Planned Elasticsearch integration for enhanced search and analytics capabilities

4. **Resource-Efficient Alternative**
   - Designed as a lightweight alternative to Elastic Agent and Metricbeat
   - Focused on core log collection and processing functionalities

## Planned Enhancements

- Full Elasticsearch integration
- Expansion of supported log sources (syslog, journald, etc.)
- Implementation of advanced log parsing and categorization algorithms
- Development of real-time log streaming capabilities
- Enhanced query optimization and visualization tools

## Quick Start Guide

1. Clone the repository
2. Verify Go 1.20 or later is installed
3. Execute `go mod tidy` to resolve and download dependencies
4. Configure PostgreSQL connection parameters in `config/.env`
5. Build and execute: `go build && ./GoLogger`

## Configuration

Define your database connection string in `config/.env`:

```
DATABASE_URL=postgres://u:p@localhost:5432/dbname?sslmode=disable
```

## Execution

Run `./GoLogger` to initiate the log collection process.

## Contributing

Contributions are welcome. Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and contribution procedures.

## License

This project is distributed under the MIT License. See the [LICENSE](LICENSE) file for full details.

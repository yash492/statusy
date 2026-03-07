# statusy

Statusy aggregates incidents and component status from external provider status pages into a single backend.
It stores normalized data in PostgreSQL and serves it for API and UI consumption.

## Dependencies Used

### Runtime & Tooling

- Golang (Go 1.25)
- Bun
- Just (justfile)
- PostgreSQL
- Docker Compose (for local services)

### Code Generation

- TypeSpec compiler (`tsp`)
- oapi-codegen


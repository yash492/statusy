# statusy

> WIP: This project is actively under development.

Statusy aggregates incidents and component status from external provider status pages into a single backend.
It stores normalized data in PostgreSQL and serves it for API and UI consumption.

## Running the Application

Run the backend, frontend, and database services:

```bash
docker compose up --build
```

- **Frontend UI:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8081/api](http://localhost:8081/api)

To stop all services:

```bash
docker compose down
```

## Supported Status Pages

- [Anthropic Claude](https://status.claude.com)
- [CircleCI](https://status.circleci.com)
- [Cloudflare](https://www.cloudflarestatus.com)
- [Cursor](https://status.cursor.com)
- [Datadog](https://status.datadoghq.com)
- [DigitalOcean](https://status.digitalocean.com)
- [Discord](https://discordstatus.com)
- [Dropbox](https://status.dropbox.com)
- [GitHub](https://www.githubstatus.com)
- [New Relic](https://status.newrelic.com)
- [Plivo](https://status.plivo.com)
- [SolarWinds Observability](https://status.cloud.solarwinds.com)
- [Twilio](https://status.twilio.com)
- [Zoom](https://www.zoomstatus.com)

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



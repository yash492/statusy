# statusy

Statusy is a backend service that aggregates incidents, component statuses, and scheduled maintenances from external provider status pages into a single normalized store. It periodically scrapes supported status pages, persists the data in PostgreSQL, and dispatches notifications to configured channels (Slack, Discord, MS Teams, PagerDuty, webhooks, and SolarWinds Incident Response) when things change.

Built with Antigravity & Claude Code.

---

<img width="854" height="512" alt="output" src="https://github.com/user-attachments/assets/9d1c56a5-11a9-408b-8abe-c782dfaba572" />


## What it does

- **Scrapes** status pages from 17 providers on a configurable interval
- **Normalizes** incident and component status data into a common format
- **Stores** everything in PostgreSQL (with separate read/write pool connections)
- **Queues** notification jobs via [PGMQ](https://github.com/tembo-io/pgmq) (Postgres-native message queue)
- **Dispatches** alerts to notification channels when incidents open, update, or resolve
- **Serves** an HTTP API (OpenAPI-generated) for UI and external consumption
- **Ships as a single binary** — PostgreSQL is the only external dependency (PGMQ is installed as a Postgres extension)

## Supported Status Pages

| Provider | Page |
|---|---|
| Anthropic Claude | status.claude.com |
| CircleCI | status.circleci.com |
| Cloudflare | cloudflarestatus.com |
| Cursor | status.cursor.com |
| Datadog | status.datadoghq.com |
| DigitalOcean | status.digitalocean.com |
| Discord | discordstatus.com |
| Dropbox | status.dropbox.com |
| GitHub | githubstatus.com |
| HashiCorp | status.hashicorp.com |
| New Relic | status.newrelic.com |
| OpenAI | status.openai.com |
| Pleo | status.pleo.io |
| Plivo | status.plivo.com |
| SolarWinds Observability | status.cloud.solarwinds.com |
| Twilio | status.twilio.com |
| Zoom | zoomstatus.com |

If the status page you need isn't listed, feel free to open a pull request adding it or open up an issue.

## Notification Channels

Slack · Discord · Microsoft Teams · PagerDuty · SolarWinds Incident Response · Generic Webhook

## Running Locally

Start the backend, database, and PGMQ installer with Docker Compose:

```bash
docker compose up --build
```

- **Backend API:** http://localhost:8081/api

To stop:

```bash
docker compose down
```

For hot-reload development (requires [air](https://github.com/air-verse/air) and [just](https://github.com/casey/just)):

```bash
just air
```

## Tech Stack

| Layer | Tool |
|---|---|
| Language | Go 1.26 |
| Database | PostgreSQL 18 |
| Message Queue | PGMQ (Postgres extension) |
| HTTP router | chi |
| DB driver | pgx/v5 |
| API spec | TypeSpec + oapi-codegen |
| Migrations | goose |
| HTTP client | resty v3 |
| Frontend | Bun (see `_ui/`) |

## Architecture

The project follows hexagonal architecture (ports and adapters). Dependencies flow inward toward the core domain:

```
cmd/               → entry point, wires everything
internal/
  domain/          → core models and port interfaces
  command/         → use cases / business logic
  adapter/pgx/     → PostgreSQL adapters (read/write pool)
  adapter/collector/ → status page scrapers
  adapter/notification/ → channel dispatchers
  port/            → inbound HTTP handlers (OpenAPI strict)
  applications/    → runtime lifecycle (server, scraper, dispatcher)
  common/          → shared helpers (queue, IDs, errors)
```

## License

[AGPL-3.0](./LICENSE)

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
| 1Password | status.1password.com |
| Acquia | status.acquia.com |
| Affirm | status.affirm.com |
| Afterpay | status.afterpay.com |
| Airbrake | airbrake.statuspage.io |
| Aircall | status.aircall.com |
| Airtable | status.airtable.com |
| Aiven | status.aiven.io |
| Akamai | akamaistatus.com |
| Alchemy | status.alchemy.com |
| Amazon Web Services (AWS) | health.aws.amazon.com/health/status |
| Amplitude | amplitude.statuspage.io |
| Anthropic Claude | status.claude.com |
| Apache Software Foundation | apache.statuspage.io |
| Asana | status.asana.com |
| Ashby | status.ashbyhq.com |
| AssemblyAI | status.assemblyai.com |
| Atlan | status.atlan.com |
| Atlassian | status.atlassian.com |
| Bitbucket | bitbucket.status.atlassian.com |
| BitGo | status.bitgo.com |
| Bitpanda | status.bitpanda.com |
| Bitrise | status.bitrise.io |
| Bitso | status.bitso.com |
| Blackboard | status.blackboard.com |
| Blockchain | status.blockchain.com |
| Boomi | status.boomi.com |
| Box | status.box.com |
| Brave | brave.statuspage.io |
| Brex | status.brex.com |
| BrowserStack | status.browserstack.com |
| Bubble | status.bubble.io |
| Buffer | status.buffer.com |
| Buildkite | buildkitestatus.com |
| Calendly | calendlystatus.com |
| Calm | status.calm.com |
| Canva | canvastatus.com |
| Cash App | status.cash.app |
| Chime | status.chime.com |
| Circle | circle.statuspage.io |
| CircleCI | status.circleci.com |
| Clerk | status.clerk.com |
| ClickUp | status.clickup.com |
| Close | status.close.com |
| Cloudflare | cloudflarestatus.com |
| Cloudinary | status.cloudinary.com |
| Cloudways | status.cloudways.com |
| Clover | status.clover.com |
| CockroachDB | status.cockroachlabs.cloud |
| Cohere | status.cohere.com |
| Coinbase | coinbase.statuspage.io |
| Confluent | status.confluent.cloud |
| Contentsquare | contentsquare.statuspage.io |
| Copper | status.copper.com |
| Courier | status.courier.com |
| Coursera | status.coursera.org |
| Crisp | crisp.statuspage.io |
| Criteo | criteo.statuspage.io |
| Cursor | status.cursor.com |
| Dailymotion | status.dailymotion.com |
| Datadog | status.datadoghq.com |
| Datadog AP1 | status.ap1.datadoghq.com |
| Datadog AP2 | status.ap2.datadoghq.com |
| Datadog EU | status.datadoghq.eu |
| Datadog GovCloud | status.ddog-gov.com |
| Datadog US3 | status.us3.datadoghq.com |
| Datadog US5 | status.us5.datadoghq.com |
| dbt | status.getdbt.com |
| DigiCert | digicert.statuspage.io |
| DigitalOcean | status.digitalocean.com |
| Discord | discordstatus.com |
| Docker | dockerstatus.com |
| Dropbox | status.dropbox.com |
| Duolingo | status.duolingo.com |
| Elastic | status.elastic.co |
| Elastic Email | elasticemail.statuspage.io |
| ElevenLabs | status.elevenlabs.io |
| Epic Games | status.epicgames.com |
| Etsy | etsystatus.com |
| Expo | status.expo.dev |
| Factorial | factorial.statuspage.io |
| Figma | status.figma.com |
| Fin | finstatus.com |
| Fireworks AI | status.fireworks.ai |
| Fivetran Activations (Census) | status.getcensus.com |
| Flickr | flickr.statuspage.io |
| Fly.io | status.flyio.net |
| Frontegg | frontegg.statuspage.io |
| Gandi | gandi.statuspage.io |
| Gemini | status.gemini.com |
| Ghost | ghoststatus.org |
| Giphy | status.giphy.com |
| GitHub | githubstatus.com |
| GoDaddy | godaddy.statuspage.io |
| Grafana | status.grafana.com |
| Grammarly | status.grammarly.com |
| Greenhouse | status.greenhouse.io |
| Gusto | status.gusto.com |
| Harness | status.harness.io |
| HashiCorp | status.hashicorp.com |
| Heap | status.heap.io |
| Hex | hex.statuspage.io |
| Highnote | highnote.statuspage.io |
| Hightouch | status.hightouch.io |
| Honeycomb | status.honeycomb.io |
| Hootsuite | status.hootsuite.com |
| Hopper | status.hopper.com |
| HubSpot | status.hubspot.com |
| Hypixel | status.hypixel.net |
| imgix | status.imgix.com |
| Imgur | status.imgur.com |
| incident.io | status.incident.io |
| Instructure | status.instructure.com |
| Iterable | status.iterable.com |
| JFrog | status.jfrog.io |
| Jira | jira-software.status.atlassian.com |
| Khan Academy | status.khanacademy.org |
| Kinsta | kinsta.statuspage.io |
| Klarna | status.klarna.com |
| Klaviyo | status.klaviyo.com |
| Knock | status.knock.app |
| Kong | status.konghq.com |
| Kraken.io | kraken.statuspage.io |
| Lattice | status.lattice.com |
| LaunchDarkly | status.launchdarkly.com |
| Ledger | status.ledger.com |
| Lever | lever.statuspage.io |
| Lightspeed | status.lightspeedhq.com |
| Linear | linearstatus.com |
| LinkedIn | linkedin.statuspage.io |
| Linode | status.linode.com |
| Lithic | lithic.statuspage.io |
| LiveChat | livechat.statuspage.io |
| Loops | status.loops.so |
| MailerSend | status.mailersend.com |
| Mailgun | status.mailgun.com |
| Matillion | status.matillion.com |
| Medium | medium.statuspage.io |
| Meraki | meraki.statuspage.io |
| Mercado Pago | status.mercadopago.com |
| Mercari | status.mercari.com |
| Mercury | status.mercury.com |
| Metabase | status.metabase.com |
| Miro | status.miro.com |
| Monday | status.monday.com |
| MongoDB | status.mongodb.com |
| MoonPay | moonpay.statuspage.io |
| Mozilla | mozilla.statuspage.io |
| Mux | status.mux.com |
| MyFitnessPal | status.myfitnesspal.com |
| n8n | n8n.statuspage.io |
| Namely | status.namely.com |
| Netlify | netlifystatus.com |
| New Relic | status.newrelic.com |
| ngrok | status.ngrok.com |
| Notion | notion-status.com |
| npm | npm.statuspage.io |
| Octopus | status.octopus.com |
| Olark | status.olark.com |
| OneSignal | status.onesignal.com |
| Onfido | onfido.statuspage.io |
| OpenAI | status.openai.com |
| Oura | status.ouraring.com |
| Pantheon | status.pantheon.io |
| Patreon | status.patreon.com |
| Paxos | status.paxos.com |
| Perplexity | status.perplexity.com |
| Personio | personio.statuspage.io |
| Pinecone | status.pinecone.io |
| Ping Identity | status.pingidentity.com |
| Pinterest | status.pinterest.com |
| Plaid | status.plaid.com |
| PlanetScale | planetscale.statuspage.io |
| Pleo | status.pleo.io |
| Plivo | status.plivo.com |
| Poshmark | status.poshmark.com |
| Postman | status.postman.com |
| Preset | status.preset.io |
| Productboard | productboard.statuspage.io |
| PubNub | status.pubnub.com |
| Pulumi | status.pulumi.com |
| Pusher | status.pusher.com |
| Quicknode | status.quicknode.com |
| Reddit | reddit.statuspage.io |
| Render | status.render.com |
| Resend | resend-status.com |
| Retool | status.retool.com |
| Ring | status.ring.com |
| Rippling | status.rippling.com |
| Robinhood | status.robinhood.com |
| Rollbar | status.rollbar.com |
| RudderStack | status.rudderstack.com |
| SailPoint | status.sailpoint.com |
| Sauce Labs | status.saucelabs.com |
| Scaleway | status.scaleway.com |
| Segment | status.segment.com |
| Semaphore | status.semaphoreci.com |
| Sentry | status.sentry.io |
| Sezzle | sezzle.statuspage.io |
| Sift | status.sift.com |
| Sigma Computing | sigmacomputing.statuspage.io |
| Signifyd | status.signifyd.com |
| Sinch | status.sinch.com |
| Sisense | status.sisense.com |
| Smartcar | status.smartcar.com |
| Smartlook | smartlook.statuspage.io |
| Smartsheet | status.smartsheet.com |
| SMTP2GO | status.smtp2go.com |
| Snowflake | status.snowflake.com |
| Snyk | status.snyk.io |
| SoFi | sofi.statuspage.io |
| SolarWinds Observability | status.cloud.solarwinds.com |
| Sonos | status.sonos.com |
| Splunk Observability Cloud US0 | signalfx.statuspage.io |
| SproutSocial | sproutsocial.statuspage.io |
| Squarespace | status.squarespace.com |
| Stitch | status.stitchdata.com |
| Strava | status.strava.com |
| Stream | status.getstream.io |
| Sumo Logic | status.sumologic.com |
| Sumsub | status.sumsub.com |
| Supabase | status.supabase.com |
| Sync | status.sync.com |
| Tailscale | status.tailscale.com |
| Telekom | telekom.statuspage.io |
| Travis CI | traviscistatus.com |
| Tray.ai | status.tray.io |
| Trello | trello.status.atlassian.com |
| Twilio | status.twilio.com |
| Typeform | status.typeform.com |
| Udemy | status.udemy.com |
| Upstash | status.upstash.com |
| Velocity | velocity.statuspage.io |
| Vercel | vercel-status.com |
| Vimeo | vimeo.statuspage.io |
| Wasabi | status.wasabi.com |
| Whoop | status.whoop.com |
| Wise | status.wise.com |
| Wix | status.wix.com |
| Workable | workable.statuspage.io |
| WorkOS | status.workos.com |
| Yubico | status.yubico.com |
| Zapier | status.zapier.com |
| Zoom | zoomstatus.com |
| Zuora | zuora.statuspage.io |

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

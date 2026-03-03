# Notes regarding the feature

### Phase 1:
- Scrape API based statuspages:
    - Atlassian Statuspage
    - incident.io Statuspage
    - GCP Statuspage
    - Azure Statuspage
    - AWS Statuspage

- Build a UI for the following spec:
    - Show incidents for a service
    - Have a search bar for services
    - Have a separate page to show incident updates
    - Show actual incident link in the statuspage

- Tech Checklist
    - Unified incident schema across providers
    - Retries + timeout + backoff for fetch jobs
    - Idempotent upserts (no duplicate incidents/updates)
    - Structured logs + error alerts
    - Basic health check endpoint


**Important**: These things needs to be deployed on prod. Exact deployment strategy can be figured out. 
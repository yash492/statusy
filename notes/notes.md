# Notes regarding the feature

### Phase 1:
- Scrape API based statuspages:
    - Atlassian Statuspage - Basic Scrapping Logic is done
    - incident.io Statuspage
    - GCP Statuspage
    - Azure Statuspage
    - AWS Statuspage

- Build a UI for the following spec:
    - Show incidents for a service (UI Part is done)
    - Have a search bar for services - ✅
    - Have a separate page to show incident updates - ❌ (Not doing it for this cut)
    - Show actual incident link in the statuspage ✅ (UI Part to redirect to the actual statuspage is done)
    - Subscribe to Updates (UI Part is done)
        - RSS
        - Atom
        - Slack Atom Subscribe

- Tech Checklist
    - Unified incident schema across providers ✅
    - Retries + timeout + backoff for fetch jobs
    - Idempotent upserts (no duplicate incidents/updates) ✅
    - Structured logs + error alerts
    - Basic health check endpoint
    - Implement Schedule Maintenances
    - Have good erroring system, right now not HTTP Handlers are returning proper errors.


**Important**: These things needs to be deployed on prod. Exact deployment strategy can be figured out. 

**Refactor**:
- Remove all mention of urls from services DB. Keep the services DB very lean. ID, name and slug should only be the 3 columns present.
- Remove services.yaml file and hard code the urls within app logic
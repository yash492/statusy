# Statusy: Your All-In-One Status Page Aggregator

**Stay on top of service health with ease.**

**Statusy** lets you aggregate publicly available status pages from various services into a single, centralized view, keeping you informed about potential disruptions and updates. It's built with Go for efficiency and offers flexible notification options to suit your workflow.


## Features

- **Consolidated view:** See the status of multiple services at a glance, eliminating the need to check multiple websites.
- **Public status page scraping:** Supports Plivo, Twilio, and other publicly available status pages (with plans for expansion).
- **Easy setup with Docker Compose:** Get started in minutes with a simple configuration.
- **Single binary application:** No complex installation processes, just a single executable for straightforward deployment.
- **Minimal dependencies:** Relies only on PostgreSQL, streamlining setup and maintenance.
- **Versatile notifications:** Receive alerts through:
    - **Incident management tools:** Squadcast, PagerDuty
    - **ChatOps platforms:** Slack, Microsoft Teams, Discord
    - **Custom webhooks:** Integrate with any system that supports them


## Getting Started

1. **Prerequisites:**
    - Docker and docker-compose installed

2. **Clone the repository:**
    ```bash
    git clone https://github.com/yash492/statusy.git
    ```

3. **Start the application:**
    ```bash
    cd statusy
    docker-compose up -d
    ```


**Need a status page added?**

If the status page you're looking for is not yet available in StatuSy, you can either:


1. **Create a pull request**: If you're comfortable with code, feel free to add the status page yourself and contribute through a pull request.

2. **Create an issue**: If you're not able to contribute code, create an issue to let us know about the missing status page. This will help us prioritize and work on adding it.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3). See the `LICENSE` file for details.


**Let Statusy be your eyes on service health.**

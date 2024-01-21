# Statusy: Your All-In-One Status Page Aggregator

**Statusy** lets you aggregate publicly available status pages from various services into a single, centralized view, keeping you informed about potential disruptions and updates. It's built with Go for efficiency and offers flexible notification options to suit your workflow.

## Motivation 

At Squadcast, where I contribute as a Software Developer, our dependency on various SaaS tools is critical for the success of our operations. However, when these tools encounter issues, we often find ourselves in the dark about the root causes, resorting to manual checks on multiple status pages. This not only proves to be time-consuming but also inefficient for our team.

The motivation behind Statusy lies in addressing this challenge. Statusy is designed to transform service monitoring by centralizing information from different SaaS tools, eliminating the need for manual checks on scattered status pages. The goal is to empower teams like ours with timely alerts and notifications, ensuring prompt awareness of any service disruptions or issues.

In essence, Statusy tackles the pain points of relying on disparate status pages, offering a centralized hub for monitoring and dispatching notifications that seamlessly integrate with our primary alerting tools. This not only enhances our operational efficiency but also provides a more proactive approach to addressing potential service disruptions. Statusy is not just a tool; it's a solution crafted to elevate the way teams handle and respond to issues in the dynamic landscape of SaaS dependencies.


## Features

- **Consolidated view:** See the status of multiple services at a glance, eliminating the need to check multiple websites.
- **Public status page scraping:** Supports Plivo, Twilio, and other publicly available status pages (with plans for expansion).
- **Easy setup with Docker Compose:** Get started in minutes with a simple configuration.
- **Single binary application:** No complex installation processes, just a single executable for straightforward deployment.
- **Minimal dependencies:** Relies only on PostgreSQL, streamlining setup and maintenance.
- **Versatile notifications:** Receive alerts through:
    - **Incident management tools:** Squadcast, PagerDuty
    - **ChatOps platforms:** Slack, Microsoft Teams, Discord
    - **Outgoing Webhooks:** Integrate with external systems


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


## Need a status page added?

If the status page you're looking for is not yet available in Statusy, you can either:


1. **Create a pull request**: If you're comfortable with code, feel free to add the status page yourself and contribute through a pull request.

2. **Create an issue**: If you're not able to contribute code, create an issue to let us know about the missing status page. This will help us prioritize and work on adding it.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3). See the `LICENSE` file for details.



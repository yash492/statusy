import type { PageLoad } from './$types';

export interface ServiceInfo {
	id: number;
	name: string;
	slug: string;
	provider: string;
	status: 'operational' | 'degraded' | 'major_outage' | 'maintenance';
	components: string[];
	lastChecked: string;
	lastIncident: string | null;
	history: ('operational' | 'degraded' | 'major_outage' | 'maintenance')[];
}

export interface ViewData {
	slug: string;
	name: string;
	description: string;
	services: ServiceInfo[];
}

const paymentServices: ServiceInfo[] = [
	{
		id: 0,
		name: 'RazorPay',
		slug: 'razorpay',
		provider: 'Incident.io',
		status: 'operational',
		components: ['API', 'Dashboard'],
		lastChecked: '5 mins ago',
		lastIncident: 'Degraded performance in ChatGPT response times',
		history: [...Array(10).fill('operational'), 'degraded', 'operational', 'operational', 'degraded', 'degraded']
	},
	{
		id: 99,
		name: 'Claude',
		slug: 'openai',
		provider: 'Incident.io',
		status: 'degraded',
		components: ['API', 'ChatGPT'],
		lastChecked: '5 mins ago',
		lastIncident: 'Degraded performance in ChatGPT response times',
		history: [...Array(10).fill('operational'), 'degraded', 'operational', 'operational', 'degraded', 'degraded']
	},
	{
		id: 2,
		name: 'Stripe',
		slug: 'stripe',
		provider: 'Incident.io',
		status: 'operational',
		components: ['API', 'ChatGPT'],
		lastChecked: '5 mins ago',
		lastIncident: 'Degraded performance in ChatGPT response times',
		history: [...Array(10).fill('operational'), 'degraded', 'operational', 'operational', 'degraded', 'degraded']
	},
	{
		id: 3,
		name: 'Plivo',
		slug: 'plivo',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['SMS Gateway', 'Voice Calls'],
		lastChecked: '1 min ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 4,
		name: 'Twilio',
		slug: 'twilio',
		provider: 'Atlassian Statuspage',
		status: 'major_outage',
		components: ['WhatsApp Messaging', 'Console'],
		lastChecked: '30s ago',
		lastIncident: 'Major Outage on SMS delivery in EU regions',
		history: [...Array(12).fill('operational'), 'operational', 'operational', 'major_outage']
	}
];

const infraServices: ServiceInfo[] = [
	{
		id: 5,
		name: 'GitHub',
		slug: 'github',
		provider: 'Atlassian Statuspage',
		status: 'degraded',
		components: ['Git Operations', 'GitHub Actions', 'Codespaces'],
		lastChecked: '1 min ago',
		lastIncident: 'Degraded performance on GitHub Actions runs',
		history: [...Array(11).fill('operational'), 'degraded', 'operational', 'operational', 'degraded']
	},
	{
		id: 6,
		name: 'Cloudflare',
		slug: 'cloudflare',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['DNS', 'CDN', 'Workers'],
		lastChecked: '2 mins ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 7,
		name: 'DigitalOcean',
		slug: 'digitalocean',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['Droplets', 'App Platform', 'Volumes'],
		lastChecked: '4 mins ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 8,
		name: 'Datadog',
		slug: 'datadog',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['APM', 'Log Management', 'Synthetics'],
		lastChecked: '1 min ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 9,
		name: 'SolarWinds Observability',
		slug: 'solarwindsobservability',
		provider: 'Atlassian Statuspage',
		status: 'major_outage',
		components: ['Alerting', 'Ingest Pipeline'],
		lastChecked: '1 min ago',
		lastIncident: 'Ingest latency spike and processing delays',
		history: [...Array(13).fill('operational'), 'degraded', 'major_outage']
	}
];

const devtoolsServices: ServiceInfo[] = [
	{
		id: 10,
		name: 'Cursor',
		slug: 'cursor',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['Copilot API', 'Backend Sync'],
		lastChecked: '3 mins ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 11,
		name: 'CircleCI',
		slug: 'circleci',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['Build Runners', 'API'],
		lastChecked: '5 mins ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 12,
		name: 'Claude',
		slug: 'claude',
		provider: 'Atlassian Statuspage',
		status: 'degraded',
		components: ['Claude API', 'Claude.ai Web App'],
		lastChecked: '2 mins ago',
		lastIncident: 'Claude API rate limits returning elevated 5xx errors',
		history: [...Array(14).fill('operational'), 'degraded']
	},
	{
		id: 13,
		name: 'New Relic',
		slug: 'newrelic',
		provider: 'Atlassian Statuspage',
		status: 'operational',
		components: ['Metrics Ingestion', 'User Interface'],
		lastChecked: '4 mins ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	},
	{
		id: 14,
		name: 'HashiCorp',
		slug: 'hashicorp',
		provider: 'Incident.io',
		status: 'operational',
		components: ['HCP Portal', 'Vault Secrets'],
		lastChecked: '3 mins ago',
		lastIncident: null,
		history: Array(15).fill('operational')
	}
];

function generateMockServices(slug: string): ServiceInfo[] {
	const cleanSlug = slug.replace(/[^a-zA-Z0-9]/g, ' ');
	const capitalized = cleanSlug
		.split(' ')
		.map((w) => w.charAt(0).toUpperCase() + w.slice(1))
		.join(' ');

	return [
		{
			id: 101,
			name: `${capitalized} Core API`,
			slug: `${slug}-core`,
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['Data Ingestion', 'Query Engine'],
			lastChecked: '1 min ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 102,
			name: `${capitalized} Dashboard`,
			slug: `${slug}-dashboard`,
			provider: 'Incident.io',
			status: 'degraded',
			components: ['UI Panel', 'Reports Export'],
			lastChecked: '3 mins ago',
			lastIncident: 'Minor latency in loading analytics reports',
			history: [...Array(13).fill('operational'), 'degraded', 'degraded']
		}
	];
}

export const load: PageLoad = async ({ params }) => {
	const slug = params.slug;

	let viewData: ViewData;

	if (slug === 'payment-gateways' || slug === 'payments') {
		viewData = {
			slug,
			name: 'Payment Gateways',
			description: 'Overview of external billing pathways, payment aggregators, and SMS gateways.',
			services: paymentServices
		};
	} else if (slug === 'core-infrastructure' || slug === 'infra') {
		viewData = {
			slug,
			name: 'Core Infrastructure',
			description: 'Key cloud hosting providers, deployment targets, and CDNs powering our applications.',
			services: infraServices
		};
	} else if (slug === 'developer-tools' || slug === 'dev-tools') {
		viewData = {
			slug,
			name: 'Developer Tools',
			description: 'Status of dev pipelines, observability, test runners, and coding co-pilots.',
			services: devtoolsServices
		};
	} else {
		// Dynamic generated slug
		const cleanSlug = slug.replace(/[^a-zA-Z0-9]/g, ' ');
		const name = cleanSlug
			.split(' ')
			.map((w) => w.charAt(0).toUpperCase() + w.slice(1))
			.join(' ');

		viewData = {
			slug,
			name: `${name} View`,
			description: `Auto-generated monitor board for services tagged with ${slug}.`,
			services: generateMockServices(slug)
		};
	}

	return {
		view: viewData
	};
};

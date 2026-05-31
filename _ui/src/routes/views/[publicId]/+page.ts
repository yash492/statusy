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
	public_id: string;
	name: string;
	description: string;
	services: ServiceInfo[];
}



export const load: PageLoad = async ({ params }) => {
	const publicId = params.publicId;

	return {
		view: {
			id: 1,
			name: 'Default View',
			public_id: publicId || 'default-view',
			description: 'Default status page aggregator view',
			is_default: true,
			services: [
				{
					id: 1,
					name: 'Stripe',
					slug: 'stripe',
					provider: 'Incident.io',
					status: 'operational',
					components: ['API', 'Checkout'],
					lastChecked: '5 mins ago',
					lastIncident: null,
					history: Array(15).fill('operational')
				},
				{
					id: 2,
					name: 'GitHub',
					slug: 'github',
					provider: 'Atlassian Statuspage',
					status: 'degraded',
					components: ['Git Operations', 'GitHub Actions'],
					lastChecked: '1 min ago',
					lastIncident: 'Degraded performance on GitHub Actions runs',
					history: [...Array(11).fill('operational'), 'degraded', 'operational', 'operational', 'degraded']
				}
			]
		}
	};
};

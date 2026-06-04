import { ApiClient } from '$lib/api/ky/ky';

export interface Statuspage {
	id: number;
	name: string;
	slug: string;
	url: string;
}

export interface Incident {
	id: number;
	title: string;
	status: string;
	provider_created_at: string;
	incident_url: string;
}

export interface StatuspageIncidents {
	statuspage: Statuspage;
	incidents: Incident[];
	total_count: number;
}

export interface ScheduledMaintenance {
	id: number;
	title: string;
	status: string;
	scheduled_maintenance_url: string;
	starts_at: string;
	ends_at: string;
	provider_created_at: string;
}

export interface StatuspageScheduledMaintenances {
	statuspage: Statuspage;
	scheduled_maintenances: ScheduledMaintenance[];
	total_count: number;
}

export interface Component {
	id: number;
	name: string;
	provider_id: string;
}

export interface ComponentGroup {
	id: number;
	name: string;
	provider_id: string;
	components: Component[];
}

export interface ServiceComponents {
	service_id: number;
	service_name: string;
	service_slug: string;
	grouped_components: ComponentGroup[];
	ungrouped_components: Component[];
}

export class StatuspageApi {
	private readonly basePath = 'statuspages';

	list(search?: string) {
		return ApiClient.get<Statuspage[]>(this.basePath, {
			searchParams: search ? { search } : undefined
		});
	}

	bySlug(slug: string) {
		return ApiClient.get<Statuspage>(`${this.basePath}/${encodeURIComponent(slug)}`);
	}

	incidents(slug: string, pageNumber = 1, pageSize = 10, componentIds?: number[], componentGroupIds?: number[]) {
		return ApiClient.get<StatuspageIncidents>(
			`${this.basePath}/${encodeURIComponent(slug)}/incidents`,
			{
				searchParams: {
					page_number: pageNumber,
					page_size: pageSize,
					...(componentIds && componentIds.length > 0 ? { component_ids: componentIds.join(',') } : {}),
					...(componentGroupIds && componentGroupIds.length > 0 ? { component_group_ids: componentGroupIds.join(',') } : {})
				}
			}
		);
	}

	scheduledMaintenances(slug: string, pageNumber = 1, pageSize = 10, componentIds?: number[], componentGroupIds?: number[]) {
		return ApiClient.get<StatuspageScheduledMaintenances>(
			`${this.basePath}/${encodeURIComponent(slug)}/schedule-maintenances`,
			{
				searchParams: {
					page_number: pageNumber,
					page_size: pageSize,
					...(componentIds && componentIds.length > 0 ? { component_ids: componentIds.join(',') } : {}),
					...(componentGroupIds && componentGroupIds.length > 0 ? { component_group_ids: componentGroupIds.join(',') } : {})
				}
			}
		);
	}

	getComponents(serviceSlug: string) {
		return ApiClient.get<ServiceComponents>(
			`services/${encodeURIComponent(serviceSlug)}/components`
		);
	}
}

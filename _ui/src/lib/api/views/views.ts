import { ApiClient } from '$lib/api/ky/ky';

export interface EditViewRequest {
	name: string;
	description: string;
	is_default: boolean;
}

export interface ViewServiceStatus {
	id: number;
	name: string;
	slug: string;
	status: string;
	last_incident: string;
	include_all_components: boolean;
	monitor_incidents: boolean;
	monitor_scheduled_maintenances: boolean;
	upcoming_maintenance: string;
	last_incident_link: string;
	upcoming_maintenance_link: string;
}

export interface View {
	id: number;
	name: string;
	public_id: string;
	description: string;
	is_default: boolean;
}

export interface AddViewServiceRequest {
	service_id: number;
	include_all_components: boolean;
	monitor_incidents: boolean;
	monitor_scheduled_maintenances: boolean;
	component_ids?: number[];
	component_group_ids?: number[];
}

export interface EditViewServiceRequest {
	include_all_components: boolean;
	monitor_incidents: boolean;
	monitor_scheduled_maintenances: boolean;
	component_ids?: number[];
	component_group_ids?: number[];
}

export interface ViewServiceResponse {
	id: number;
	service_id: number;
	include_all_components: boolean;
	monitor_incidents: boolean;
	monitor_scheduled_maintenances: boolean;
	component_ids?: number[];
	component_group_ids?: number[];
}

export class ViewsApi {
	private readonly basePath = 'views';

	get(publicId: string) {
		return ApiClient.get<View>(`${this.basePath}/${encodeURIComponent(publicId)}`);
	}

	edit(publicId: string, body: EditViewRequest) {
		return ApiClient.put<View>(`${this.basePath}/${encodeURIComponent(publicId)}`, {
			json: body
		});
	}

	delete(publicId: string) {
		return ApiClient.delete(`${this.basePath}/${encodeURIComponent(publicId)}`);
	}

	getUnconfiguredServices(viewPublicId: string, search?: string) {
		return ApiClient.get<
			{
				id: number;
				name: string;
				slug: string;
				url: string;
			}[]
		>(`${this.basePath}/${encodeURIComponent(viewPublicId)}/unconfigured-services`, {
			searchParams: search ? { search } : undefined
		});
	}

	addViewService(viewPublicId: string, body: AddViewServiceRequest) {
		return ApiClient.post<ViewServiceResponse>(
			`${this.basePath}/${encodeURIComponent(viewPublicId)}/services`,
			{
				json: body
			}
		);
	}

	getViewService(viewPublicId: string, serviceId: number) {
		return ApiClient.get<ViewServiceResponse>(
			`${this.basePath}/${encodeURIComponent(viewPublicId)}/services/${serviceId}`
		);
	}

	editViewService(viewPublicId: string, serviceId: number, body: EditViewServiceRequest) {
		return ApiClient.put<ViewServiceResponse>(
			`${this.basePath}/${encodeURIComponent(viewPublicId)}/services/${serviceId}`,
			{
				json: body
			}
		);
	}

	deleteViewService(viewPublicId: string, serviceId: number) {
		return ApiClient.delete(
			`${this.basePath}/${encodeURIComponent(viewPublicId)}/services/${serviceId}`
		);
	}

	getViewServices(viewPublicId: string, pageNumber?: number, pageSize?: number, search?: string) {
		return ApiClient.get<{
			services: ViewServiceStatus[];
			total_count: number;
			up_count: number;
			down_count: number;
		}>(`${this.basePath}/${encodeURIComponent(viewPublicId)}/view-services`, {
			searchParams: {
				...(pageNumber !== undefined ? { page_number: pageNumber } : {}),
				...(pageSize !== undefined ? { page_size: pageSize } : {}),
				...(search ? { search } : {})
			}
		});
	}

	createOrGetDefaultView() {
		return ApiClient.post<View>(`${this.basePath}/default`);
	}

	list(search?: string) {
		return ApiClient.get<View[]>(`${this.basePath}`, {
			searchParams: search ? { search } : undefined
		});
	}

	create(body: { name: string; description: string }) {
		return ApiClient.post<View>(`${this.basePath}`, {
			json: body
		});
	}
}

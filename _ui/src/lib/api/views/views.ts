import KyClient from '$lib/api/ky/ky';

export interface EditViewRequest {
	name: string;
	slug: string;
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
}

export interface View {
	id: number;
	name: string;
	slug: string;
	description: string;
	is_default: boolean;
	services: ViewServiceStatus[];
}

export interface AddViewServiceRequest {
	service_id: number;
	include_all_components: boolean;
	component_ids?: number[];
	component_group_ids?: number[];
}

export interface EditViewServiceRequest {
	include_all_components: boolean;
	component_ids?: number[];
	component_group_ids?: number[];
}

export interface ViewServiceResponse {
	id: number;
	service_id: number;
	include_all_components: boolean;
}

export class ViewsApi {
	private readonly basePath = 'views';

	edit(slug: string, body: EditViewRequest) {
		return KyClient.put(`${this.basePath}/${encodeURIComponent(slug)}`, {
			json: body
		}).json<View>();
	}

	delete(slug: string) {
		return KyClient.delete(`${this.basePath}/${encodeURIComponent(slug)}`);
	}

	getUnconfiguredServices(viewSlug: string, search?: string) {
		return KyClient.get(`${this.basePath}/${encodeURIComponent(viewSlug)}/unconfigured-services`, {
			searchParams: search ? { search } : undefined
		}).json<{ id: number; name: string; slug: string; url: string }[]>();
	}

	addViewService(viewSlug: string, body: AddViewServiceRequest) {
		return KyClient.post(`${this.basePath}/${encodeURIComponent(viewSlug)}/services`, {
			json: body
		}).json<ViewServiceResponse>();
	}

	editViewService(viewSlug: string, serviceId: number, body: EditViewServiceRequest) {
		return KyClient.put(`${this.basePath}/${encodeURIComponent(viewSlug)}/services/${serviceId}`, {
			json: body
		}).json<ViewServiceResponse>();
	}

	deleteViewService(viewSlug: string, serviceId: number) {
		return KyClient.delete(`${this.basePath}/${encodeURIComponent(viewSlug)}/services/${serviceId}`);
	}
}

import KyClient from '$lib/api/ky/ky';

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
}

export interface View {
	id: number;
	name: string;
	public_id: string;
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

	edit(publicId: string, body: EditViewRequest) {
		return KyClient.put(`${this.basePath}/${encodeURIComponent(publicId)}`, {
			json: body
		}).json<View>();
	}

	delete(publicId: string) {
		return KyClient.delete(`${this.basePath}/${encodeURIComponent(publicId)}`);
	}

	getUnconfiguredServices(viewPublicId: string, search?: string) {
		return KyClient.get(`${this.basePath}/${encodeURIComponent(viewPublicId)}/unconfigured-services`, {
			searchParams: search ? { search } : undefined
		}).json<{ id: number; name: string; slug: string; url: string }[]>();
	}

	addViewService(viewPublicId: string, body: AddViewServiceRequest) {
		return KyClient.post(`${this.basePath}/${encodeURIComponent(viewPublicId)}/services`, {
			json: body
		}).json<ViewServiceResponse>();
	}

	editViewService(viewPublicId: string, serviceId: number, body: EditViewServiceRequest) {
		return KyClient.put(`${this.basePath}/${encodeURIComponent(viewPublicId)}/services/${serviceId}`, {
			json: body
		}).json<ViewServiceResponse>();
	}

	deleteViewService(viewPublicId: string, serviceId: number) {
		return KyClient.delete(`${this.basePath}/${encodeURIComponent(viewPublicId)}/services/${serviceId}`);
	}
}

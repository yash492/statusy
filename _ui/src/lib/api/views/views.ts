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
	component_ids?: number[];
	component_group_ids?: number[];
}

export class ViewsApi {
	private readonly basePath = 'views';

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

	createOrGetDefaultView() {
		return ApiClient.post<View>(`${this.basePath}/default`);
	}
}

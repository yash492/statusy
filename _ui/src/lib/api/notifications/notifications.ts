import { ApiClient } from '$lib/api/ky/ky';

export interface ViewNotification {
	id: number;
	name: string;
	type: string;
	config: Record<string, unknown>;
	created_at: string;
}

export interface ViewNotificationRequest {
	name: string;
	type: 'slack' | 'discord' | 'msteams' | 'pagerduty' | 'solarwinds_incident_response' | 'webhook';
	config: Record<string, unknown>;
}

export interface PaginatedViewNotifications {
	notifications: ViewNotification[];
	total_count: number;
}

export class NotificationsApi {
	private basePath(viewPublicId: string) {
		return `views/${encodeURIComponent(viewPublicId)}/notifications`;
	}

	list(viewPublicId: string, pageNumber?: number, pageSize?: number, search?: string) {
		const params = new URLSearchParams();
		if (pageNumber !== undefined) params.set('page_number', String(pageNumber));
		if (pageSize !== undefined) params.set('page_size', String(pageSize));
		if (search !== undefined && search !== '') params.set('search', search);

		const queryString = params.toString();
		const path = this.basePath(viewPublicId) + (queryString ? `?${queryString}` : '');
		return ApiClient.get<PaginatedViewNotifications>(path);
	}

	create(viewPublicId: string, body: ViewNotificationRequest) {
		return ApiClient.post<ViewNotification>(this.basePath(viewPublicId), {
			json: body
		});
	}

	edit(viewPublicId: string, id: number, body: ViewNotificationRequest) {
		return ApiClient.put<ViewNotification>(`${this.basePath(viewPublicId)}/${id}`, {
			json: body
		});
	}

	delete(viewPublicId: string, id: number) {
		return ApiClient.delete(`${this.basePath(viewPublicId)}/${id}`);
	}
}

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

export class NotificationsApi {
	private basePath(viewPublicId: string) {
		return `views/${encodeURIComponent(viewPublicId)}/notifications`;
	}

	list(viewPublicId: string) {
		return ApiClient.get<ViewNotification[]>(this.basePath(viewPublicId));
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

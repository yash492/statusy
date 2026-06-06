import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { ViewsApi } from '$lib/api/views/views';
import { NotificationsApi } from '$lib/api/notifications/notifications';

export const load: PageLoad = async ({ params, parent }) => {
	const publicId = params.publicId;
	const defaultViewData = await parent();
	const defaultView = defaultViewData.defaultView!;

	let view;
	if (defaultView.public_id === publicId) {
		view = defaultView;
	} else {
		const viewsApi = new ViewsApi();
		const [viewResult, viewErr] = await viewsApi.get(publicId);
		if (viewErr || !viewResult) {
			throw error(404, {
				message: 'View not found'
			});
		}
		view = viewResult;
	}

	const notificationsApi = new NotificationsApi();
	const [notifications, notifErr] = await notificationsApi.list(publicId);
	if (notifErr) {
		throw error(500, {
			message: 'Failed to load notifications'
		});
	}

	return { view, notifications: notifications ?? [] };
};

import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { NotificationsApi } from '$lib/api/notifications/notifications';

export const load: PageLoad = async ({ params, url }) => {
	const publicId = params.publicId;
	const page = Number(url.searchParams.get('page') || '1');
	const pageSize = Number(url.searchParams.get('page_size') || '5');
	const search = url.searchParams.get('search') || '';

	const notificationsApi = new NotificationsApi();
	const [res, notifErr] = await notificationsApi.list(publicId, page, pageSize, search);
	if (notifErr || !res) {
		throw error(500, {
			message: 'Failed to load notifications'
		});
	}

	return {
		notifications: res.notifications,
		totalCount: res.total_count,
		page,
		pageSize,
		search
	};
};

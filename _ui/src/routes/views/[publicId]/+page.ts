import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { ViewsApi } from '$lib/api/views/views';

export const load: PageLoad = async ({ params, parent }) => {
	const publicId = params.publicId;
	const defaultViewData = await parent();
	const defaultView = defaultViewData.defaultView!;

	if (defaultView.public_id === publicId) {
		return { view: defaultView };
	}

	const viewsApi = new ViewsApi();
	const [view, err] = await viewsApi.get(publicId);
	if (err || !view) {
		throw error(404, {
			message: 'View not found'
		});
	}

	return { view };
};

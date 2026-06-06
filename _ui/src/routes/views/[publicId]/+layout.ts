import { error } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { ViewsApi } from '$lib/api/views/views';

export const load: LayoutLoad = async ({ params, parent }) => {
	const publicId = params.publicId;
	const defaultViewData = await parent();
	const defaultView = defaultViewData.defaultView;

	let view;
	if (defaultView && defaultView.public_id === publicId) {
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

	return { view };
};

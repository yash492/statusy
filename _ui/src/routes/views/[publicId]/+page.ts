import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, parent }) => {
	const publicId = params.publicId;
	const defaultViewData = await parent();
	const defaultView = defaultViewData.defaultView!;

	if (defaultView.public_id === publicId) {
		return { view: defaultView };
	}

	throw error(404, {
		message: 'View not found'
	});
};

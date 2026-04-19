import { PUBLIC_API_SERVER_ROUTE } from '$env/static/public';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params, fetch }) => {
	const { slug } = params;
	const url = `${PUBLIC_API_SERVER_ROUTE}/statuspages/${encodeURIComponent(slug)}/feed.atom`;
	
	const response = await fetch(url);
	
	return new Response(response.body, {
		status: response.status,
		headers: {
			'content-type': response.headers.get('content-type') || 'application/atom+xml; charset=utf-8'
		}
	});
};

import type { HandleFetch } from '@sveltejs/kit';
import { PUBLIC_API_SERVER_ROUTE, PUBLIC_LOCAL_API_SERVER_ROUTE } from '$env/static/public';

export const handleFetch: HandleFetch = async ({ request, fetch }) => {
	if (request.url.startsWith(PUBLIC_API_SERVER_ROUTE)) {
		// clone the original request, but change the URL
		request = new Request(
			request.url.replace(PUBLIC_API_SERVER_ROUTE, PUBLIC_LOCAL_API_SERVER_ROUTE),
			request
		);
	}

	return fetch(request);
};

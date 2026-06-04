import { ViewsApi } from '$lib/api/views/views';

export async function load() {
	const viewAPI = new ViewsApi();
	const [defaultView, defaultErr] = await viewAPI.createOrGetDefaultView();
	const [views, viewsErr] = await viewAPI.list();
	return { defaultView, views: views ?? [], err: defaultErr || viewsErr };
}

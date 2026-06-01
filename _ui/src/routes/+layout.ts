import { ViewsApi } from '$lib/api/views/views';

export async function load() {
	const viewAPI = new ViewsApi();
	const [defaultView, err] = await viewAPI.createOrGetDefaultView();
	return { defaultView, err };
}

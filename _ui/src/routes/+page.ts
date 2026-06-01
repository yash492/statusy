import { redirect } from '@sveltejs/kit';

export async function load({ parent }) {
	const data = await parent();
	if (data.err !== null) {
		return data.err;
	}

	redirect(308, `/views/${data.defaultView?.public_id}`);
}

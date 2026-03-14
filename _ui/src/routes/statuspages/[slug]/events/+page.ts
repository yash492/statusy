import { StatuspageApi } from "$lib/api/statuspage/statuspage"

export async function load({ params }) {
    const statuspageApi = new StatuspageApi()
    const incidents = await statuspageApi.incidents(params.slug)
    return {
        resp: incidents
    }
}

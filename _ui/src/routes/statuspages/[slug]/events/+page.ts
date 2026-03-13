import { PUBLIC_API_SERVER_ROUTE } from "$env/static/public"
import { StatuspageApi } from "$lib/api/statuspage/statuspage"

export async function load({ params }) {
    const statuspageApi = new StatuspageApi()
    console.log("public", PUBLIC_API_SERVER_ROUTE)
    const incidents = await statuspageApi.incidents(params.slug)
    return {
        resp: incidents
    }
}

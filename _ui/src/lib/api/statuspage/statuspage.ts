import KyClient from "$lib/api/ky/ky";

export interface Statuspage {
    id: number;
    name: string;
    slug: string;
}

export interface Incident {
    id: number;
    title: string;
    status: string;
    provider_created_at: string;
}

export interface StatuspageIncidents {
    statuspage: Statuspage;
    incidents: Incident[];
}

class StatuspageApi {
    private readonly basePath = "api/v1/statuspages";

    list(search?: string) {
        return KyClient.get(this.basePath, {
            searchParams: search ? { search } : undefined,
        }).json<Statuspage[]>();
    }

    bySlug(slug: string) {
        return KyClient.get(`${this.basePath}/${encodeURIComponent(slug)}`).json<Statuspage>();
    }

    incidents(slug: string, pageNumber = 1, pageSize = 10) {
        return KyClient
            .get(`${this.basePath}/${encodeURIComponent(slug)}/incidents`, {
                searchParams: {
                    page_number: pageNumber,
                    page_size: pageSize,
                },
            })
            .json<StatuspageIncidents[]>();
    }
}

const statuspageApi = new StatuspageApi();

export default statuspageApi;
export { StatuspageApi };

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

export class StatuspageApi {
    private readonly basePath = "statuspages";

    list(search?: string) {
        return KyClient.get(this.basePath, {
            searchParams: search ? { search } : undefined,
        }).json<Statuspage[]>();
    }

    bySlug(slug: string) {
        return KyClient.get(`${this.basePath}/${encodeURIComponent(slug)}`).json<Statuspage>();
    }

    incidents(slug: string, pageNumber = 1, pageSize = 10) {
        console.log(KyClient)
        return KyClient
            .get(`${this.basePath}/${encodeURIComponent(slug)}/incidents`, {
                searchParams: {
                    page_number: pageNumber,
                    page_size: pageSize,
                },
            })
            .json<StatuspageIncidents>();
    }
}


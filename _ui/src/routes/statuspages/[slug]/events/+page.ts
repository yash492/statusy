import { StatuspageApi } from '$lib/api/statuspage/statuspage';

const DEFAULT_PAGE = 1;
const DEFAULT_PAGE_SIZE = 10;
const DEFAULT_TYPE = 'incidents';

type TabType = 'incidents' | 'scheduled-maintenance';

function parsePositiveInt(raw: string | null, fallback: number): number {
    if (!raw) {
        return fallback;
    }

    const parsed = Number.parseInt(raw, 10);
    if (!Number.isFinite(parsed) || parsed <= 0) {
        return fallback;
    }

    return parsed;
}

function parseType(raw: string | null): TabType {
    if (raw === 'incidents' || raw === 'scheduled-maintenance') {
        return raw;
    }

    return DEFAULT_TYPE;
}

export async function load({ params, url }) {
    const page = parsePositiveInt(url.searchParams.get('page'), DEFAULT_PAGE);
    const pageSize = parsePositiveInt(url.searchParams.get('page_size'), DEFAULT_PAGE_SIZE);
    const type = parseType(url.searchParams.get('type'));

    const statuspageApi = new StatuspageApi();
    const incidents = await statuspageApi.incidents(params.slug, page, pageSize);

    return {
        resp: incidents,
        page,
        pageSize,
        type
    };
}

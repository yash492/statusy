import { StatuspageApi } from '$lib/api/statuspage/statuspage';
import { redirect } from '@sveltejs/kit';

const DEFAULT_PAGE = 1;
const DEFAULT_PAGE_SIZE = 10;
const DEFAULT_TYPE = 'incidents';

type TabType = 'incidents' | 'scheduled-maintenances';

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
    if (raw === 'incidents' || raw === 'scheduled-maintenances') {
        return raw;
    }

    return DEFAULT_TYPE;
}

export async function load({ params, url }) {
    const rawPage = url.searchParams.get('page');
    const rawPageSize = url.searchParams.get('page_size');
    const rawType = url.searchParams.get('type');

    const page = parsePositiveInt(url.searchParams.get('page'), DEFAULT_PAGE);
    const pageSize = parsePositiveInt(url.searchParams.get('page_size'), DEFAULT_PAGE_SIZE);
    const type = parseType(url.searchParams.get('type'));

    // Keep URL canonical by always including normalized pagination and tab params.
    if (rawPage !== String(page) || rawPageSize !== String(pageSize) || rawType !== type) {
        const paramsWithDefaults = new URLSearchParams(url.searchParams);
        paramsWithDefaults.set('page', String(page));
        paramsWithDefaults.set('page_size', String(pageSize));
        paramsWithDefaults.set('type', type);

        throw redirect(302, `${url.pathname}?${paramsWithDefaults.toString()}`);
    }

    const statuspageApi = new StatuspageApi();
    let resp;
    if (type === 'scheduled-maintenances') {
        resp = await statuspageApi.scheduledMaintenances(params.slug, page, pageSize);
    } else {
        resp = await statuspageApi.incidents(params.slug, page, pageSize);
    }

    return {
        resp,
        page,
        pageSize,
        type
    };
}

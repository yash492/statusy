<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import * as Tabs from '$lib/components/ui/tabs';
	import IncidentsTable, { type Incident } from '$lib/derivedcomponents/IncidentsTable.svelte';
	import ScheduledMaintenancesTable, {
		type ScheduledMaintenanceDisplay
	} from '$lib/derivedcomponents/ScheduledMaintenancesTable.svelte';
	import type { PaginationState } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
	type TabType = 'incidents' | 'scheduled-maintenance';
	type SubscribeTab = 'rss' | 'slack';
	type CopyKey = 'rss' | 'atom' | 'slack';

	let isSubscribeDialogOpen = $state(false);
	let subscribeTab = $state<SubscribeTab>('rss');
	let copiedKey = $state<CopyKey | null>(null);
	let copyResetTimer: ReturnType<typeof setTimeout> | null = null;

	const PAGE_SIZE = 10;
	const initialPagination: PaginationState = {
		pageIndex: Math.max(0, data.page - 1),
		pageSize: data.pageSize
	};

	const isIncidents = $derived('incidents' in data.resp);
	const incidentsArray = $derived(isIncidents ? (data.resp as any).incidents : []);
	const scheduledMaintenancesArray = $derived(
		!isIncidents ? (data.resp as any).scheduled_maintenances : []
	);

	function toIncidents(raw: any[]): Incident[] {
		return raw.map((incident) => ({
			created_at: incident.provider_created_at,
			id: incident.id,
			status: incident.status,
			title: incident.title,
			incident_url: incident.incident_url
		}));
	}

	function toScheduledMaintenances(raw: any[]): ScheduledMaintenanceDisplay[] {
		return raw.map((m) => ({
			id: m.id,
			title: m.title,
			status: m.status,
			starts_at: m.starts_at,
			ends_at: m.ends_at,
			scheduled_maintenance_url: m.scheduled_maintenance_url
		}));
	}

	const incidentData = $derived(toIncidents(incidentsArray));
	const scheduledMaintenanceData = $derived(toScheduledMaintenances(scheduledMaintenancesArray));
	const currentListLength = $derived(
		isIncidents ? incidentsArray.length : scheduledMaintenancesArray.length
	);

	const feedBaseUrl = $derived(
		`${$page.url.origin}/statuspages/${encodeURIComponent(data.resp.statuspage.slug)}`
	);
	const feedRssPath = $derived(`${feedBaseUrl}/feed.rss`);
	const feedAtomPath = $derived(`${feedBaseUrl}/feed.atom`);
	const slackSnippet = $derived(`/feed subscribe ${feedRssPath}`);

	const rowCount = $derived(data.resp.total_count);

	async function onPageChange(pagination: PaginationState) {
		const params = new URLSearchParams(window.location.search);
		params.set('page', String(pagination.pageIndex + 1));
		params.set('page_size', String(pagination.pageSize));
		await goto(`?${params.toString()}`, {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	}

	async function onTabChange(type: TabType) {
		const params = new URLSearchParams(window.location.search);
		params.set('type', type);
		params.set('page', '1');
		params.set('page_size', String(data.pageSize));

		await goto(`?${params.toString()}`, {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	}

	function onSubscribeDialogOpenChange(open: boolean) {
		isSubscribeDialogOpen = open;
		if (!open) {
			subscribeTab = 'rss';
			copiedKey = null;
			if (copyResetTimer) {
				clearTimeout(copyResetTimer);
				copyResetTimer = null;
			}
		}
	}

	async function copyText(value: string, key: CopyKey) {
		try {
			await navigator.clipboard.writeText(value);
			copiedKey = key;
			toast.success('Copied feed snippet to clipboard');

			if (copyResetTimer) {
				clearTimeout(copyResetTimer);
			}

			copyResetTimer = setTimeout(() => {
				copiedKey = null;
				copyResetTimer = null;
			}, 1800);
		} catch {
			toast.error('Could not copy. Please copy manually.');
		}
	}
</script>

<div class="mx-auto w-11/12">
	<div class="w-full">
		<div class="mb-6 flex justify-between md:mb-4">
			<div class="mb-4 flex w-fit flex-col gap-3 md:flex-row md:items-center">
				<p class="text-3xl font-bold">{data.resp.statuspage.name}</p>
				{#if data.resp.statuspage.url}
					<a
						href={data.resp.statuspage.url}
						target="_blank"
						rel="noopener noreferrer"
						class="flex w-fit items-center gap-1.5 rounded-full border border-zinc-800 bg-zinc-950/80 px-3 py-1 text-xs text-zinc-400 shadow-sm transition-all hover:border-zinc-700 hover:bg-zinc-900 hover:text-white"
					>
						<span>Visit Status Page</span>
					</a>
				{/if}
			</div>
		</div>

		<div>
			<div>
				<Tabs.Root
					value={data.type}
					onValueChange={(value) => onTabChange(value as TabType)}
					class="w-full"
				>
					<Tabs.List>
						<Tabs.Trigger value="incidents">Incidents</Tabs.Trigger>
						<Tabs.Trigger value="scheduled-maintenances">Scheduled Maintenances</Tabs.Trigger>
					</Tabs.List>
					<Tabs.Content value="incidents">
						<IncidentsTable
							data={incidentData}
							{rowCount}
							paginationState={initialPagination}
							{onPageChange}
						/>
					</Tabs.Content>
					<Tabs.Content value="scheduled-maintenances">
						<ScheduledMaintenancesTable
							data={scheduledMaintenanceData}
							{rowCount}
							paginationState={initialPagination}
							{onPageChange}
						/>
					</Tabs.Content>
				</Tabs.Root>
			</div>
		</div>
	</div>
</div>

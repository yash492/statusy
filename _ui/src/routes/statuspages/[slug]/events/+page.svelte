<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import IncidentsTable, { type Incident } from '$lib/derivedcomponents/IncidentsTable.svelte';
	import type { PaginationState } from '@tanstack/table-core';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
	type TabType = 'incidents' | 'scheduled-maintenance';

	const PAGE_SIZE = 10;
	const initialPagination: PaginationState = {
		pageIndex: Math.max(0, data.page - 1),
		pageSize: data.pageSize
	};

	function toIncidents(raw: typeof data.resp.incidents): Incident[] {
		return raw.map((incident) => ({
			created_at: incident.provider_created_at,
			id: incident.id,
			status: incident.status,
			title: incident.title
		}));
	}

	const incidentData = $derived(toIncidents(data.resp.incidents));

	// If first page is full, we don't know the total — use MAX to keep Next enabled.
	// Once a page returns fewer rows than pageSize, we know the exact total.
	const rowCount = $derived(
		data.resp.incidents.length < PAGE_SIZE
			? (data.page - 1) * data.pageSize + data.resp.incidents.length
			: Number.MAX_SAFE_INTEGER
	);

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
</script>

<div class="mx-auto w-4/5">
	<div class="w-full">
		<div class="mb-6 flex justify-between md:mb-4">
			<div class="mb-4 flex w-fit items-center gap-2">
				<div class="rounded-4xl border-2 bg-white">
					<img src="/provider_logo/plivo.png" width="30px" height="30px" alt="plivo-logo" />
				</div>
				<p class="text-xl font-bold">{data.resp.statuspage.name}</p>
			</div>
			<div>
				<Button class="cursor-pointer">Subscribe to Updates</Button>
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
						<Tabs.Trigger value="scheduled-maintenance">Scheduled Maintenances</Tabs.Trigger>
					</Tabs.List>
					<Tabs.Content value="incidents">
						<IncidentsTable
							data={incidentData}
							{rowCount}
							paginationState={initialPagination}
							{onPageChange}
						/>
					</Tabs.Content>
					<Tabs.Content value="scheduled-maintenance">
						<!-- <IncidentsTable /> -->
					</Tabs.Content>
				</Tabs.Root>
			</div>
		</div>
	</div>
</div>

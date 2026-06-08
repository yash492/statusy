<script lang="ts">
	import { goto } from '$app/navigation';
	import IncidentsTable, { type Incident } from '$lib/components/custom/IncidentsTable.svelte';
	import ScheduledMaintenancesTable, {
		type ScheduledMaintenanceDisplay
	} from '$lib/components/custom/ScheduledMaintenancesTable.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tabs from '$lib/components/ui/tabs';
	import Search from '@lucide/svelte/icons/search';
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

	let isFilterDialogOpen = $state(false);
	let filterSearchQuery = $state('');

	const initialPagination: PaginationState = $derived({
		pageIndex: Math.max(0, data.page - 1),
		pageSize: data.pageSize
	});

	const hasFilter = $derived(data.componentIds.length > 0 || data.componentGroupIds.length > 0);

	const activeFilters = $derived.by(() => {
		if (!data.serviceComponents) return { groups: [], ungrouped: [] };

		const groupIdsSet = new Set(data.componentGroupIds);
		const compIdsSet = new Set(data.componentIds);
		const query = filterSearchQuery.toLowerCase().trim();

		// Filter selected groups
		const groups = data.serviceComponents.grouped_components
			.map((g) => {
				const isGroupSelected = groupIdsSet.has(g.id);
				const selectedComps = g.components.filter((c) => compIdsSet.has(c.id));

				if (!isGroupSelected && selectedComps.length === 0) {
					return null;
				}

				// Check name match
				const matchesGroupQuery = g.name.toLowerCase().includes(query);
				const filteredComps = selectedComps.filter((c) => c.name.toLowerCase().includes(query));

				if (matchesGroupQuery || filteredComps.length > 0 || isGroupSelected) {
					return {
						id: g.id,
						name: g.name,
						isFullySelected: isGroupSelected,
						components: isGroupSelected
							? g.components.filter((c) => c.name.toLowerCase().includes(query))
							: filteredComps
					};
				}
				return null;
			})
			.filter((g): g is NonNullable<typeof g> => g !== null);

		// Filter selected ungrouped components
		const ungrouped = data.serviceComponents.ungrouped_components.filter((c) => {
			const isSelected = compIdsSet.has(c.id);
			const matchesQuery = c.name.toLowerCase().includes(query);
			return isSelected && matchesQuery;
		});

		return { groups, ungrouped };
	});

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

<div class="mx-auto">
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
				{#if hasFilter}
					<div class="flex items-center gap-2">
						<span
							class="flex w-fit items-center gap-1.5 rounded-full border border-zinc-800 bg-zinc-900/30 px-3 py-1 text-xs text-zinc-400"
						>
							<span class="size-1.5 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.5)]"
							></span>
							Filter Active
						</span>
						<button
							onclick={() => {
								filterSearchQuery = '';
								isFilterDialogOpen = true;
							}}
							class="flex cursor-pointer items-center gap-1 rounded-full border border-zinc-800 bg-zinc-950/80 px-3 py-1 text-xs text-zinc-400 shadow-sm transition-all hover:border-zinc-700 hover:bg-zinc-900 hover:text-white"
						>
							View Filters
						</button>
					</div>
				{:else}
					<div class="flex items-center gap-2">
						<span
							class="flex w-fit items-center gap-1.5 rounded-full border border-zinc-800 bg-zinc-900/30 px-3 py-1 text-xs text-zinc-400"
						>
							<span
								class="size-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"
							></span>
							All Components Monitored
						</span>
					</div>
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

<!-- Component Filter Dialog -->
<Dialog.Root
	open={isFilterDialogOpen}
	onOpenChange={(open) => {
		if (!open) {
			isFilterDialogOpen = false;
		}
	}}
>
	<Dialog.Content class="border-zinc-800 bg-zinc-950 text-white shadow-xl sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title class="text-lg font-bold text-white">Active Component Filters</Dialog.Title>
			<Dialog.Description class="text-sm text-zinc-400">
				Only incidents and scheduled maintenance events affecting the following components are
				displayed.
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<!-- Search Bar -->
			<div class="relative w-full">
				<input
					type="text"
					bind:value={filterSearchQuery}
					placeholder="Search active filters..."
					class="w-full rounded-lg border border-zinc-800 bg-zinc-900/40 py-2 pr-4 pl-9 text-sm text-white placeholder-zinc-500 outline-none focus:border-zinc-700"
				/>
				<div class="absolute inset-y-0 left-0 flex items-center pl-3 text-zinc-500">
					<Search class="size-4" />
				</div>
			</div>

			<!-- Filter List -->
			<div class="flex max-h-75 flex-col gap-3 overflow-y-auto pr-1">
				{#each activeFilters.groups as g (g.id)}
					<div class="rounded-lg border border-zinc-900 bg-zinc-950/40 p-3">
						<div class="flex items-center justify-between">
							<span class="text-sm font-semibold text-zinc-200">{g.name}</span>
							{#if g.isFullySelected}
								<span
									class="rounded border border-zinc-800/40 bg-zinc-900 px-1.5 py-0.5 text-[9px] font-medium tracking-wider text-zinc-500 uppercase"
								>
									Group Monitored
								</span>
							{/if}
						</div>
						{#if !g.isFullySelected && g.components.length > 0}
							<div class="mt-2 flex flex-col gap-1 border-t border-zinc-900/50 pt-2">
								{#each g.components as c (c.id)}
									<span class="flex items-center gap-1.5 pl-2 text-xs text-zinc-200">
										<span class="size-1 rounded-full bg-zinc-500"></span>
										{c.name}
									</span>
								{/each}
							</div>
						{:else if g.isFullySelected && g.components.length > 0}
							<div class="mt-2 flex flex-col gap-1 border-t border-zinc-900/50 pt-2">
								{#each g.components as c (c.id)}
									<span class="flex items-center gap-1.5 pl-2 text-xs text-zinc-200">
										<span class="size-1 rounded-full bg-zinc-500"></span>
										{c.name}
									</span>
								{/each}
							</div>
						{/if}
					</div>
				{/each}

				{#each activeFilters.ungrouped as c (c.id)}
					<div
						class="flex items-center justify-between rounded-lg border border-zinc-900 bg-zinc-950/40 p-3"
					>
						<span class="text-sm font-medium text-zinc-200">{c.name}</span>
						<span
							class="rounded border border-zinc-800/40 bg-zinc-900 px-1.5 py-0.5 text-[9px] font-medium tracking-wider text-zinc-500 uppercase"
						>
							Component Monitored
						</span>
					</div>
				{/each}

				{#if activeFilters.groups.length === 0 && activeFilters.ungrouped.length === 0}
					<div class="py-8 text-center text-sm text-zinc-500">
						No active component filters match your query.
					</div>
				{/if}
			</div>
		</div>

		<Dialog.Footer>
			<button
				class="w-full cursor-pointer rounded-lg bg-zinc-100 py-2 text-sm font-semibold text-zinc-950 transition-colors hover:bg-zinc-200"
				onclick={() => {
					isFilterDialogOpen = false;
				}}
			>
				Close
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

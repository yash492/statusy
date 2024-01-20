<script lang="ts">
	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import type { SubscriptionIncidentsTable } from '$lib/types/subscriptions';
	import { flexRender, type ColumnDef } from '@tanstack/svelte-table';
	import IncidentTable from './IncidentTable.svelte';
	import { SUBSCRIPTION_INCIDENT_LIST } from '$lib/types/query_keys';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import Status from './IncidentStatus.svelte';
	import IncidentName from './IncidentName.svelte';
	import { page } from '$app/stores';
	import { afterNavigate } from '$app/navigation';

	const _subscriptionAPI = new SubscriptionAPI();

	const queryClient = useQueryClient();

	let incidentTableData: SubscriptionIncidentsTable[] = [];
	let serviceName = '';
	let components = 'All';
	let pageNumber = 0;
	let totalCount = 0;
	const PAGE_LIMIT = 10;

	afterNavigate(() => {
		queryClient.invalidateQueries({ queryKey: [SUBSCRIPTION_INCIDENT_LIST] });
	});

	$: query = createQuery({
		queryKey: [SUBSCRIPTION_INCIDENT_LIST],
		queryFn: async () =>
			await _subscriptionAPI.GetAllIncidents($page.params.subscription_id, pageNumber, PAGE_LIMIT)
	});

	$: if ($query.isSuccess) {
		let subscriptionIncidentDetails = $query.data?.data.data;
		incidentTableData =
			subscriptionIncidentDetails.incidents?.map<SubscriptionIncidentsTable>((incident) => {
				const incidentData: SubscriptionIncidentsTable = {
					createdAt: incident.created_at,
					link: incident.link,
					name: incident.name,
					status: incident.status,
					normalisedStatus: incident.normalised_status
				};

				return incidentData;
			}) || [];

		serviceName = subscriptionIncidentDetails.service_name;
		totalCount = $query.data?.data.meta.total_count;

		if (!subscriptionIncidentDetails.is_all_components_configured) {
			components = subscriptionIncidentDetails?.components
				.map((component) => component.name)
				.join(', ');
		}
	}

	let columns: ColumnDef<SubscriptionIncidentsTable>[] = [
		{
			accessorKey: 'createdAt',
			header: () => 'Created At',
			cell: (info) => new Date(info.getValue() as Date).toLocaleString()
		},
		{
			accessorKey: 'name',
			header: () => 'Name',
			cell: (info) =>
				flexRender(IncidentName, { name: info.getValue(), link: info.row.original.link })
		},
		{
			accessorKey: 'status',
			header: () => 'Status',
			cell: (info) =>
				flexRender(Status, {
					normalisedStatus: info.row.original.normalisedStatus,
					status: info.getValue()
				})
		}
	];
</script>

<div>
	{#if $query.isLoading}
		<p>Loading...</p>
	{:else if $query.isError}
		<p>{AxiosResponseErr($query.error).error_msg}</p>
	{:else}
		<div class="w-3/4 mx-auto">
			<h1 class="text-lg font-bold w-full">
				{serviceName} Incidents
			</h1>
			<p class="text-sm text-neutral-100">
				Total Incidents: {totalCount}
			</p>
			<p
				class="mt-1 text-sm text-neutral-100 text-ellipsis w-3/4 whitespace-nowrap overflow-hidden"
			>
				Configured Components: {components}
			</p>
			<IncidentTable
				{columns}
				data={incidentTableData}
				pageLimit={PAGE_LIMIT}
				bind:pageNumber
				{totalCount}
			/>
		</div>
	{/if}
</div>

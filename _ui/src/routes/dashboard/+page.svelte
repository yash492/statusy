<script lang="ts">
	import { goto } from '$app/navigation';
	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import Button from '$lib/components/button/Button.svelte';
	import Table from './Table.svelte';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import type { DashboardTable } from '$lib/types/subscriptions';
	import { createQuery } from '@tanstack/svelte-query';
	import { flexRender, type ColumnDef } from '@tanstack/svelte-table';
	import Status from './Status.svelte';
	import ServiceColumn from './ServiceColumn.svelte';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { MagnifyingGlass } from '@steeze-ui/heroicons';
	import { DASHBOARD_SUBSCRIPTION_LIST_QUERY_KEY } from '$lib/types/query_keys';

	const _subscriptionAPI = new SubscriptionAPI();
	let serviceName = '';
	let pageNumber = 0;
	const pageLimit = 5;
	let totalCount = 0;

	$: query = createQuery({
		queryKey: [DASHBOARD_SUBSCRIPTION_LIST_QUERY_KEY],
		queryFn: async () => await _subscriptionAPI.GetAll(serviceName, pageNumber, pageLimit)
	});

	let subscriptionsListData: DashboardTable[] = [];
	let columns: ColumnDef<DashboardTable>[] = [
		{
			accessorKey: 'serviceName',
			header: () => 'Service',
			cell: (info) => {
				const original = info.row.original;
				return flexRender(ServiceColumn, {
					serviceName: original.serviceName,
					incidentLink: original.incidentLink,
					incidentName: original.incident,
					subscriptionUUID: original.subscriptionUUID
				});
			}
		},
		{
			accessorKey: 'isDown',
			header: () => 'Status',
			cell: (info) => flexRender(Status, { isDown: info.getValue() })
		}
	];

	$: if ($query.isSuccess) {
		const data = $query.data.data.data;
		const tableSubscriptionData = data?.map<DashboardTable>((subscription) => {
			const tableData: DashboardTable = {
				incident: subscription.incident_name,
				serviceName: subscription.service_name,
				isDown: subscription.is_down,
				incidentLink: subscription.incident_link,
				subscriptionUUID: subscription.subscription_uuid
			};
			return tableData;
		});
		subscriptionsListData = tableSubscriptionData;
		totalCount = $query.data.data.meta.total_count || 0;
	}
</script>

<div class="w-fit mx-auto md:w-2/3 lg:w-1/2 pt-9 md:pt-12">
	<div class="flex items-center justify-between">
		<h1 class="font-bold text-lg md:text-2xl">
			Monitoring {totalCount} services
		</h1>
		<Button on:click={() => goto('/subscriptions/add-service')}>Add Service</Button>
	</div>

	<div>
		{#if $query.isLoading}
			<p>Loading...</p>
		{:else if $query.isError}
			<p>{AxiosResponseErr($query.error).error_msg}</p>
		{:else}
			<div class="mt-6">
				<div class="flex items-center rounded-md bg-neutral-800 mb-2">
					<Icon src={MagnifyingGlass} size="20px" class="ml-2"></Icon>
					<input
						class="py-2 px-3 w-full bg-neutral-800 focus:outline-none"
						placeholder="Search Services..."
						bind:value={serviceName}
					/>
				</div>
				<Table data={subscriptionsListData} {columns} {pageLimit} bind:pageNumber {totalCount} />
			</div>
		{/if}
	</div>
</div>

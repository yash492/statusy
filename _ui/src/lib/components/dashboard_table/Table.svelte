<script lang="ts">
	import { writable } from 'svelte/store';
	import {
		createSvelteTable,
		flexRender,
		getCoreRowModel,
		type ColumnDef,
		type TableOptions
	} from '@tanstack/svelte-table';
	import Status from './Status.svelte';

	type dasboardTable = {
		serviceName: string;
		status: string;
		incident: string;
	};

	const defaultData: dasboardTable[] = [
		{
			serviceName: 'Plivo',
			status: 'Outage',
			incident: ''
		},
		{
			serviceName: 'Twilio',
			status: 'Up',
			incident: ''
		},
		{
			serviceName: 'Pagerduty',
			status: 'Up',
			incident: '45'
		}
	];

	const defaultColumns: ColumnDef<dasboardTable>[] = [
		{
			accessorKey: 'serviceName',
			header: () => 'Service',
			cell: (info) => info.getValue()
		},
		{
			accessorKey: 'status',
			header: () => 'Status',
			cell: (info) => flexRender(Status, { value: info.getValue() })
		}
		// {
		// 	accessorKey: 'incident',
		// 	header: () => 'incident'
		// }
	];

	const options = writable<TableOptions<dasboardTable>>({
		data: defaultData,
		columns: defaultColumns,
		getCoreRowModel: getCoreRowModel()
	});

	const table = createSvelteTable(options);
</script>

<div class="mt-5 overflow-x-auto py-2 px-5 border-gray-600 border rounded-sm">
	<table class="w-full text-left rounded-sm shadow-sm table-fixed">
		<thead>
			{#each $table.getHeaderGroups() as headerGroup}
				<tr>
					{#each headerGroup.headers as header}
						<th class="uppercase text-xs whitespace-nowrap last:w-auto last:text-center w-3/4 pb-2">
							{#if !header.isPlaceholder}
								<svelte:component
									this={flexRender(header.column.columnDef.header, header.getContext())}
								/>
							{/if}
						</th>
					{/each}
				</tr>
			{/each}
		</thead>
		<tbody>
			{#each $table.getRowModel().rows as row}
				<tr>
					{#each row.getVisibleCells() as cell}
						<td class="py-1">
							<svelte:component this={flexRender(cell.column.columnDef.cell, cell.getContext())} />
						</td>
					{/each}
				</tr>
			{/each}
		</tbody>
	</table>
</div>

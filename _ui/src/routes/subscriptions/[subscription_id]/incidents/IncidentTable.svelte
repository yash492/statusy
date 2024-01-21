<script lang="ts">
	import { writable } from 'svelte/store';
	import {
		createSvelteTable,
		flexRender,
		getCoreRowModel,
		type ColumnDef,
		type TableOptions,
		Pagination
	} from '@tanstack/svelte-table';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ChevronLeft, ChevronRight } from '@steeze-ui/heroicons';
	import type { SubscriptionIncidentsTable } from '$lib/types/subscriptions';

	export let data: SubscriptionIncidentsTable[] = [];
	export let columns: ColumnDef<SubscriptionIncidentsTable>[];
	export let pageNumber: number;
	export let pageLimit: number;
	export let totalCount: number;

	$: pageCount = Math.ceil(totalCount / pageLimit);

	$: options = writable<TableOptions<SubscriptionIncidentsTable>>({
		data: data,
		columns: columns,
		getCoreRowModel: getCoreRowModel(),
		pageCount: pageCount,
		state: {
			pagination: {
				pageIndex: pageNumber,
				pageSize: pageLimit
			}
		}
	});

	$: table = createSvelteTable(options);
</script>

<div class="mt-5 overflow-x-auto py-2 border border-neutral-600 rounded-sm mx-auto">
	<table class="w-full text-left rounded-sm shadow-sm table-fixed">
		<thead>
			{#each $table.getHeaderGroups() as headerGroup}
				<tr>
					{#each headerGroup.headers as header}
						<th
							class="uppercase text-xs whitespace-nowrap pb-2 last:w-1/6 first:w-1/3 md:first:w-1/4 first:pl-5"
						>
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
						<td class="py-1 first:pl-5 text-sm last:capitalize">
							<svelte:component this={flexRender(cell.column.columnDef.cell, cell.getContext())} />
						</td>
					{/each}
				</tr>
			{/each}
		</tbody>
	</table>
</div>
<div class="flex items-center justify-end h-full mt-3 gap-3">
	<button
		class="border rounded-sm p-1 border-neutral-600 hover:bg-neutral-800"
		hidden={!$table.getCanPreviousPage()}
		on:click={() => (pageNumber -= 1)}
	>
		<Icon src={ChevronLeft} size="23px" />
	</button>
	<p class="h-full text-sm">Page {pageNumber + 1} of {pageCount}</p>
	<button
		class="border rounded-sm p-1 border-neutral-600 hover:bg-neutral-800"
		on:click={() => (pageNumber += 1)}
		hidden={!$table.getCanNextPage()}
	>
		<Icon src={ChevronRight} size="23px" />
	</button>
</div>

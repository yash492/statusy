<script lang="ts">
	import { writable } from 'svelte/store';
	import {
		createSvelteTable,
		flexRender,
		getCoreRowModel,
		type ColumnDef,
		type TableOptions
	} from '@tanstack/svelte-table';
	import type { DashboardTable } from '$lib/types/subscriptions';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ChevronLeft, ChevronRight } from '@steeze-ui/heroicons';

	export let data: DashboardTable[];
	export let columns: ColumnDef<DashboardTable>[];
	export let pageNumber: number;
	export let pageLimit: number;
	export let totalCount: number;

	$: options = writable<TableOptions<DashboardTable>>({
		data: data,
		columns: columns,
		getCoreRowModel: getCoreRowModel(),
		pageCount: Math.ceil(totalCount / pageLimit),
		state: {
			pagination: {
				pageIndex: pageNumber,
				pageSize: pageLimit
			}
		}
	});

	$: table = createSvelteTable(options);
</script>

<div class="mt-2 overflow-x-auto py-2 px-5 border-neutral-600 border rounded-sm">
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
<div class="mt-3 flex gap-4 justify-end">
	<button
		class="border rounded-sm p-1 border-neutral-600 hover:bg-neutral-800"
		hidden={!$table.getCanPreviousPage()}
		on:click={() => (pageNumber -= 1)}
	>
		<Icon src={ChevronLeft} size="23px" />
	</button>
	<button
		class="border rounded-sm p-1 border-neutral-600 hover:bg-neutral-800"
		hidden={!$table.getCanNextPage()}
		on:click={() => (pageNumber += 1)}
	>
		<Icon src={ChevronRight} size="23px" />
	</button>
</div>

<script lang="ts" generics="TData extends Record<string, unknown>">
	import { Button } from '$lib/components/ui/button/index.js';
	import { FlexRender, createSvelteTable } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import {
		type ColumnDef,
		type PaginationState,
		getCoreRowModel,
		getPaginationRowModel
	} from '@tanstack/table-core';

	interface Props {
		data: TData[];
		columns: ColumnDef<TData>[];
		pagination?: boolean;
		paginationState?: PaginationState;
		/** Enable server-side pagination. Parent is responsible for fetching new data on page change. */
		manualPagination?: boolean;
		/** Total row count from server — required for correct Next/Prev button state with manualPagination. */
		rowCount?: number;

		loading?: boolean;

		onRowClick?: (row: TData) => void;

		onPageChange?: (pagination: PaginationState) => void;
	}

	let {
		data,
		columns,
		pagination: enablePagination = false,
		paginationState: initialPagination = { pageIndex: 0, pageSize: 10 },
		manualPagination = false,
		rowCount,

		loading = false,
		onRowClick,
		onPageChange
	}: Props = $props();

	let paginationState = $state<PaginationState>((() => initialPagination)());

	$effect(() => {
		paginationState = initialPagination;
	});

	const table = createSvelteTable<TData>({
		get data() {
			return data;
		},
		get columns() {
			return columns;
		},
		get rowCount() {
			return rowCount;
		},
		state: {
			get pagination() {
				return paginationState;
			}
		},
		manualPagination,
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: (() => enablePagination && !manualPagination)()
			? getPaginationRowModel()
			: undefined,
		onPaginationChange: (updater) => {
			const next = typeof updater === 'function' ? updater(paginationState) : updater;
			paginationState = next;
			onPageChange?.(next);
		}
	});

	// Skeleton rows count mirrors pageSize
	const skeletonRows = $derived(paginationState.pageSize);
</script>

<div class="rounded-md border">
	<Table.Root>
		<Table.Header>
			{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
				<Table.Row>
					{#each headerGroup.headers as header (header.id)}
						<Table.Head class="[&:has([role=checkbox])]:ps-3">
							{#if !header.isPlaceholder}
								<FlexRender
									content={header.column.columnDef.header}
									context={header.getContext()}
								/>
							{/if}
						</Table.Head>
					{/each}
				</Table.Row>
			{/each}
		</Table.Header>

		<Table.Body>
			{#if loading}
				<!-- Loading skeleton -->
				{#each { length: skeletonRows } as _, i (i)}
					<Table.Row class="pointer-events-none animate-pulse">
						{#each columns as _, j (j)}
							<Table.Cell>
								<div class="h-4 w-full rounded bg-muted"></div>
							</Table.Cell>
						{/each}
					</Table.Row>
				{/each}
			{:else}
				{#each table.getRowModel().rows as row (row.id)}
					<Table.Row
						data-state={row.getIsSelected() && 'selected'}
						class={onRowClick ? 'cursor-pointer hover:bg-muted/40' : undefined}
						tabindex={onRowClick ? 0 : undefined}
						onclick={() => onRowClick?.(row.original)}
						onkeydown={(event) => {
							if (!onRowClick) {
								return;
							}

							if (event.key === 'Enter' || event.key === ' ') {
								event.preventDefault();
								onRowClick(row.original);
							}
						}}
					>
						{#each row.getVisibleCells() as cell (cell.id)}
							<Table.Cell class="[&:has([role=checkbox])]:ps-3">
								<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
							</Table.Cell>
						{/each}
					</Table.Row>
				{:else}
					<!-- Empty state -->
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24 text-center text-muted-foreground">
							No results.
						</Table.Cell>
					</Table.Row>
				{/each}
			{/if}
		</Table.Body>
	</Table.Root>
</div>

{#if enablePagination}
	<div class="flex items-center justify-end space-x-2 pt-4">
		<div class="space-x-2">
			<Button
				variant="outline"
				size="sm"
				onclick={() => table.previousPage()}
				disabled={!table.getCanPreviousPage()}
			>
				Previous
			</Button>
			<Button
				variant="outline"
				size="sm"
				onclick={() => table.nextPage()}
				disabled={!table.getCanNextPage()}
			>
				Next
			</Button>
		</div>
	</div>
{/if}

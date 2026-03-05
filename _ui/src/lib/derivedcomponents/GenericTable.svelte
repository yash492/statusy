<script lang="ts" generics="TData extends Record<string, unknown>">
	import { Button } from '$lib/components/ui/button/index.js';
	import { FlexRender, createSvelteTable } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import {
		type ColumnDef,
		type ColumnFiltersState,
		type PaginationState,
		type RowSelectionState,
		type SortingState,
		type VisibilityState,
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel
	} from '@tanstack/table-core';

	import type { FilterState, SortState } from './types.ts';

	interface Props {
		data: TData[];
		columns: ColumnDef<TData>[];
		/** Show pagination controls. Default: false */
		pagination?: boolean;
		paginationState?: PaginationState;
		/** Enable sorting. Default: false */
		sorting?: boolean;
		sortingState?: SortState[];
		/** Enable column filtering. Default: false */
		filtering?: boolean;
		filterState?: FilterState[];
		/** Show a skeleton loader instead of rows. Default: false */
		loading?: boolean;
		/** Called when user clicks a data row */
		onRowClick?: (row: TData) => void;
		/** Called when sort state changes */
		onSort?: (sorting: SortState[]) => void;
		/** Called when page state changes */
		onPageChange?: (pagination: PaginationState) => void;
	}

	let {
		data,
		columns,
		pagination: enablePagination = false,
		paginationState: initialPagination = { pageIndex: 0, pageSize: 10 },
		sorting: enableSorting = false,
		sortingState: initialSorting = [],
		filtering: enableFiltering = false,
		filterState: initialFilters = [],
		loading = false,
		onRowClick,
		onSort,
		onPageChange
	}: Props = $props();

	let paginationState = $state<PaginationState>(initialPagination);
	let sortingState = $state<SortingState>(initialSorting);
	let columnFilters = $state<ColumnFiltersState>(initialFilters as ColumnFiltersState);
	let rowSelection = $state<RowSelectionState>({});
	let columnVisibility = $state<VisibilityState>({});

	const table = createSvelteTable<TData>({
		get data() {
			return data;
		},
		columns,
		state: {
			get pagination() {
				return paginationState;
			},
			get sorting() {
				return sortingState;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get rowSelection() {
				return rowSelection;
			},
			get columnFilters() {
				return columnFilters;
			}
		},
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: enablePagination ? getPaginationRowModel() : undefined,
		getSortedRowModel: enableSorting ? getSortedRowModel() : undefined,
		getFilteredRowModel: enableFiltering ? getFilteredRowModel() : undefined,
		onPaginationChange: (updater) => {
			const next = typeof updater === 'function' ? updater(paginationState) : updater;
			paginationState = next;
			onPageChange?.(next);
		},
		onSortingChange: (updater) => {
			const next = typeof updater === 'function' ? updater(sortingState) : updater;
			sortingState = next;
			onSort?.(next as SortState[]);
		},
		onColumnFiltersChange: (updater) => {
			columnFilters = typeof updater === 'function' ? updater(columnFilters) : updater;
		},
		onColumnVisibilityChange: (updater) => {
			columnVisibility = typeof updater === 'function' ? updater(columnVisibility) : updater;
		},
		onRowSelectionChange: (updater) => {
			rowSelection = typeof updater === 'function' ? updater(rowSelection) : updater;
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
					<!-- svelte-ignore a11y_click_events_have_key_events -->
					<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
					<Table.Row
						data-state={row.getIsSelected() && 'selected'}
						class={onRowClick ? 'cursor-pointer' : ''}
						onclick={() => onRowClick?.(row.original)}
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
		<div class="flex-1 text-sm text-muted-foreground">
			{table.getFilteredSelectedRowModel().rows.length} of
			{table.getFilteredRowModel().rows.length} row(s) selected.
		</div>
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

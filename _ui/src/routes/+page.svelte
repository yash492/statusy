<script lang="ts">
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
	import { createRawSnippet } from 'svelte';

	import { Button } from '$lib/components/ui/button/index.js';
	import {
		FlexRender,
		createSvelteTable,
		renderSnippet
	} from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';

	type Payment = {
		id: string;
		amount: number;
		status: 'Pending' | 'Processing' | 'Success' | 'Failed';
		email: string;
	};

	const data: Payment[] = [
		{
			id: 'm5gr84i9',
			amount: 316,
			status: 'Success',
			email: 'ken99@yahoo.com'
		},
		{
			id: '3u1reuv4',
			amount: 242,
			status: 'Success',
			email: 'Abe45@gmail.com'
		},
		{
			id: 'derv1ws0',
			amount: 837,
			status: 'Processing',
			email: 'Monserrat44@gmail.com'
		},
		{
			id: '5kma53ae',
			amount: 874,
			status: 'Success',
			email: 'Silas22@gmail.com'
		},
		{
			id: 'bhqecj4p',
			amount: 721,
			status: 'Failed',
			email: 'carmella@hotmail.com'
		}
	];

	const columns: ColumnDef<Payment>[] = [
		{
			id: 'select',
			header: ({ table }) => {},
			cell: ({ row }) => {},
			enableSorting: false,
			enableHiding: false
		},
		{
			accessorKey: 'status',
			header: 'Title',
			cell: ({ row }) => {
				const statusSnippet = createRawSnippet<[{ status: string }]>((getStatus) => {
					const { status } = getStatus();
					return {
						render: () => `<div class="capitalize">${status}</div>`
					};
				});
				return renderSnippet(statusSnippet, {
					status: row.original.status
				});
			}
		},
		{
			accessorKey: 'email',
			header: 'Status',
			cell: ({ row }) => {
				const emailSnippet = createRawSnippet<[{ email: string }]>((getEmail) => {
					const { email } = getEmail();
					return {
						render: () => `<div class="lowercase">${email}</div>`
					};
				});

				return renderSnippet(emailSnippet, {
					email: row.original.email
				});
			}
		},
		{
			accessorKey: 'amount',
			header: () => {
				const amountHeaderSnippet = createRawSnippet(() => {
					return {
						render: () => `<div class="text-end">Amount</div>`
					};
				});
				return renderSnippet(amountHeaderSnippet);
			},
			cell: ({ row }) => {
				const formatter = new Intl.NumberFormat('en-US', {
					style: 'currency',
					currency: 'USD'
				});

				const amountCellSnippet = createRawSnippet<[{ amount: number }]>((getAmount) => {
					const { amount } = getAmount();
					const formatted = formatter.format(amount);
					return {
						render: () => `<div class="text-end font-medium">${formatted}</div>`
					};
				});
				return renderSnippet(amountCellSnippet, {
					amount: row.original.amount
				});
			}
		},
		{
			id: 'actions',
			enableHiding: false,
			cell: ({ row }) => {}
		}
	];

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let rowSelection = $state<RowSelectionState>({});
	let columnVisibility = $state<VisibilityState>({});

	const table = createSvelteTable({
		get data() {
			return data;
		},
		columns,
		state: {
			get pagination() {
				return pagination;
			},
			get sorting() {
				return sorting;
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
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
		},
		onColumnFiltersChange: (updater) => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}
		},
		onColumnVisibilityChange: (updater) => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}
		},
		onRowSelectionChange: (updater) => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}
		}
	});
</script>

<div class="mx-auto w-4/5">
	<div class="-mb-8 w-full">
		<div class="flex items-center py-4">
			<p class="text-xl font-bold">Plivo</p>
		</div>
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
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row data-state={row.getIsSelected() && 'selected'}>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell class="[&:has([role=checkbox])]:ps-3">
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
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
	</div>
</div>

<script lang="ts">
	import { renderSnippet } from '$lib/components/ui/data-table';
	import { Temporal } from '@js-temporal/polyfill';
	import type { ColumnDef } from '@tanstack/table-core';
	import { createRawSnippet } from 'svelte';
	import GenericTable from './GenericTable.svelte';

	export type Incident = {
		id: number;
		title: string;
		status: string;
		created_at: string;
		incident_url: string;
	};

	import type { PaginationState } from '@tanstack/table-core';

	let { data, rowCount, paginationState, onPageChange } = $props<{
		data: Incident[];
		rowCount?: number;
		paginationState?: PaginationState;
		onPageChange?: (p: PaginationState) => void;
	}>();

	const columns: ColumnDef<Incident>[] = [
		{
			accessorKey: 'created_at',
			header: 'Created At',
			cell: ({ row }) => {
				const createdAt = row.getValue<string>('created_at');
				const instant = Temporal.Instant.from(createdAt);
				return instant.toLocaleString();
			}
		},
		{
			accessorKey: 'title',
			header: 'Title',
			cell: ({ row }) => {
				const titleSnippet = createRawSnippet<[{ title: string }]>((getTitle) => {
					const { title } = getTitle();

					return {
						render: () => `
						<a href="${row.original.incident_url}" target="_blank">
							${title}
						</a>`
					};
				});
				return renderSnippet(titleSnippet, {
					title: row.original.title
				});
			}
		},
		{
			accessorKey: 'status',
			header: 'Status',
			cell: ({ row }) => {
				const statusSnippet = createRawSnippet<[{ status: string }]>((getStatus) => {
					const { status } = getStatus();
					let color = 'bg-red-500';
					switch (status) {
						case 'investigating':
							color = 'bg-yellow-500';
							break;
						case 'resolved':
							color = 'bg-green-500';
							break;
						case 'postmortem':
							color = 'bg-green-500';
					}
					return {
						render: () => `
						<div class="capitalize border rounded-md text-center py-1 flex w-24 justify-center items-center gap-2">
							<div class="rounded-full w-2 h-2 ${color}">
							</div>
							${status}
						</div>`
					};
				});
				return renderSnippet(statusSnippet, {
					status: row.original.status
				});
			}
		}
	];
</script>

<GenericTable
	{data}
	{columns}
	pagination
	{paginationState}
	manualPagination={true}
	{rowCount}
	{onPageChange}
	loading={false}
/>

<script lang="ts">
	import { renderSnippet } from '$lib/components/ui/data-table';
	import { Temporal } from '@js-temporal/polyfill';
	import type { ColumnDef, PaginationState } from '@tanstack/table-core';
	import { createRawSnippet } from 'svelte';
	import GenericTable from './GenericTable.svelte';

	export type ScheduledMaintenanceDisplay = {
		id: number;
		title: string;
		status: string;
		starts_at: string;
		ends_at: string;
		scheduled_maintenance_url: string;
	};

	let { data, rowCount, paginationState, onPageChange } = $props<{
		data: ScheduledMaintenanceDisplay[];
		rowCount?: number;
		paginationState?: PaginationState;
		onPageChange?: (p: PaginationState) => void;
	}>();

	function formatTime(isoString: string) {
		const instant = Temporal.Instant.from(isoString);
		return new Intl.DateTimeFormat('en-US', {
			year: 'numeric',
			month: 'short',
			day: '2-digit',
			hour: 'numeric',
			minute: '2-digit',
			hour12: true
		}).format(instant.epochMilliseconds);
	}

	const columns: ColumnDef<ScheduledMaintenanceDisplay>[] = [
		{
			accessorKey: 'title',
			header: 'Title',
			cell: ({ row }) => {
				const titleSnippet = createRawSnippet<[{ title: string }]>((getTitle) => {
					return {
						render: () => {
							const t = getTitle().title;
							const safeT = t.replace(/"/g, '&quot;');
							return '<div class="w-[50vw] md:w-[50%] max-w-[50vw] md:max-w-[50%] truncate font-medium" title="' + safeT + '">' + t + '</div>';
						}
					};
				});
				return renderSnippet(titleSnippet, { title: row.original.title });
			}
		},
		{
			accessorKey: 'starts_at',
			header: 'Starts At',
			cell: ({ row }) => formatTime(row.getValue<string>('starts_at'))
		},
		{
			accessorKey: 'ends_at',
			header: 'Ends At',
			cell: ({ row }) => formatTime(row.getValue<string>('ends_at'))
		},
		{
			accessorKey: 'status',
			header: 'Status',
			cell: ({ row }) => {
				const statusSnippet = createRawSnippet<[{ status: string }]>((getStatus) => {
					const { status } = getStatus();
					let color = 'bg-blue-500';
					switch (status) {
						case 'scheduled':
							color = 'bg-blue-500';
							break;
						case 'in_progress':
							color = 'bg-yellow-500';
							break;
						case 'verifying':
							color = 'bg-purple-500';
							break;
						case 'completed':
							color = 'bg-green-500';
							break;
						default:
							color = 'bg-gray-500';
					}
					return {
						render: () => `
						<div class="capitalize rounded-md text-center py-1 flex w-fit justify-center items-center gap-2 ">
							<div class="rounded-full w-2 h-2 ${color}"></div>
							${status.replace('_', ' ')}
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
	onRowClick={(row) => {
		window.open(row.scheduled_maintenance_url, '_blank', 'noopener,noreferrer');
	}}
	pagination
	{paginationState}
	manualPagination={true}
	{rowCount}
	{onPageChange}
	loading={false}
/>

<script lang="ts">
	import type { ColumnDef } from '@tanstack/table-core';
	import GenericTable from './GenericTable.svelte';

	type Incident = {
		id: string;
		title: string;
		status: string;
		created_at: Date;
	};

	const data: Incident[] = [];

	const columns: ColumnDef<Incident>[] = [
		{
			accessorKey: 'title',
			header: 'Title',
			cell: ({ row }) => {}
		},
		{
			accessorKey: 'status',
			header: 'Status',
			cell: ({ row }) => {}
		},
		{
			accessorKey: 'created_at',
			header: 'Created At',
			cell: ({ row }) => {}
		}
	];
</script>

<GenericTable {data} {columns} pagination loading={false} />

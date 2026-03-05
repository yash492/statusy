import type { ColumnDef } from '@tanstack/table-core';

export type { ColumnDef as GenericTableColumn };

export type SortDirection = 'asc' | 'desc' | false;

export interface SortState {
	id: string;
	desc: boolean;
}

export interface PaginationState {
	pageIndex: number;
	pageSize: number;
}

export interface FilterState {
	id: string;
	value: unknown;
}

export interface GenericTableProps<TData> {
	/** Row data */
	data: TData[];
	/** TanStack column definitions */
	columns: ColumnDef<TData>[];
	/** Enable pagination controls */
	pagination?: boolean;
	/** Initial / controlled pagination state */
	paginationState?: PaginationState;
	/** Enable column sorting */
	sorting?: boolean;
	/** Initial / controlled sorting state */
	sortingState?: SortState[];
	/** Enable column filters */
	filtering?: boolean;
	/** Initial / controlled column filter state */
	filterState?: FilterState[];
	/** Show a loading skeleton instead of data */
	loading?: boolean;
	/** Callback fired when a row is clicked */
	onRowClick?: (row: TData) => void;
	/** Callback fired when sort changes */
	onSort?: (sorting: SortState[]) => void;
	/** Callback fired when page changes */
	onPageChange?: (pagination: PaginationState) => void;
}

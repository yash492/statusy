<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { ViewsApi } from '$lib/api/views/views';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Table from '$lib/components/ui/table';
	import ViewForm from '$lib/components/statusy/ViewForm.svelte';
	import AlertTriangle from '@lucide/svelte/icons/alert-triangle';
	import Bell from '@lucide/svelte/icons/bell';
	import Calendar from '@lucide/svelte/icons/calendar';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Search from '@lucide/svelte/icons/search';
	import Settings from '@lucide/svelte/icons/settings';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const viewsApi = new ViewsApi();

	// View Settings and Meta State
	let viewName = $state('');
	let viewDescription = $state('');
	let isDefaultView = $state(false);

	// Local reactive services state
	let localServices = $state<any[]>([]);
	let totalServicesCount = $state(0);
	let upServicesCount = $state(0);
	let downServicesCount = $state(0);
	let currentPage = $state(1);
	let searchQuery = $state('');
	const itemsPerPage = 5;

	$effect(() => {
		// Sync local state when view data from page load changes
		viewName = data.view.name;
		viewDescription = data.view.description;
		isDefaultView = data.view.is_default;

		currentPage = 1;
	});

	async function loadServices() {
		const [res, err] = await viewsApi.getViewServices(
			data.view.public_id,
			currentPage,
			itemsPerPage,
			searchQuery
		);
		if (err) {
			toast.error(err.message || 'Failed to load services');
			return;
		}
		if (res) {
			localServices = res.services;
			totalServicesCount = res.total_count;
			upServicesCount = res.up_count;
			downServicesCount = res.down_count;
		}
	}

	$effect(() => {
		if (data.view.public_id) {
			const _ = currentPage;
			const __ = searchQuery;
			void loadServices();
		}
	});

	$effect(() => {
		const _ = searchQuery;
		currentPage = 1;
	});

	// Compute metrics dynamically from localState
	const totalCount = $derived(totalServicesCount);
	const upCount = $derived(upServicesCount);
	const downCount = $derived(downServicesCount);

	const totalPages = $derived(Math.ceil(totalServicesCount / itemsPerPage) || 1);

	// Deletion State
	let isDeleteConfirmOpen = $state(false);
	let serviceToDelete = $state<any | null>(null);

	function openDeleteConfirm(service: any) {
		serviceToDelete = service;
		isDeleteConfirmOpen = true;
	}

	async function confirmRemove() {
		if (serviceToDelete) {
			const [, err] = await viewsApi.deleteViewService(data.view.public_id, serviceToDelete.id);
			if (err) {
				toast.error(err.message || 'Failed to remove service');
				return;
			}
			toast.success(`Successfully removed ${serviceToDelete.name}`);
			isDeleteConfirmOpen = false;
			serviceToDelete = null;

			if (currentPage > 1 && localServices.length === 1) {
				currentPage -= 1;
			} else {
				void loadServices();
			}
		}
	}

	function navigateToAdd() {
		void goto(`/views/${data.view.public_id}/add-service`);
	}

	function navigateToEdit(serviceSlug: string) {
		void goto(`/views/${data.view.public_id}/edit-service/${serviceSlug}`);
	}

	// Edit View Dialog State
	let isEditViewOpen = $state(false);
	let editViewName = $state('');
	let editViewDescription = $state('');
	let editViewIsDefault = $state(false);

	function openEditViewDialog() {
		editViewName = viewName;
		editViewDescription = viewDescription;
		editViewIsDefault = isDefaultView;
		isEditViewOpen = true;
	}

	async function saveViewMeta() {
		const [, err] = await viewsApi.edit(data.view.public_id, {
			name: editViewName,
			description: editViewDescription,
			is_default: editViewIsDefault
		});

		if (err) {
			toast.error(err.message || 'Failed to update view details');
			return;
		}

		toast.success('View details updated successfully');
		isEditViewOpen = false;
		await invalidateAll();
	}

	// Delete View Dialog State
	let isDeleteViewOpen = $state(false);

	function openDeleteViewDialog() {
		isDeleteViewOpen = true;
	}

	async function confirmDeleteView() {
		const [, err] = await viewsApi.delete(data.view.public_id);

		if (err) {
			toast.error(err.message || 'Failed to delete view');
			return;
		}

		toast.success('View deleted successfully');
		isDeleteViewOpen = false;
		await goto('/');
		await invalidateAll();
	}

	function getEventsUrl(service: any) {
		const params = new URLSearchParams();
		if (!service.include_all_components) {
			if (service.component_ids && service.component_ids.length > 0) {
				for (const id of service.component_ids) {
					params.append('component_ids', String(id));
				}
			}
			if (service.component_group_ids && service.component_group_ids.length > 0) {
				for (const id of service.component_group_ids) {
					params.append('component_group_ids', String(id));
				}
			}
		}
		const queryString = params.toString();
		return (
			`/views/${data.view.public_id}/${service.slug}/events` +
			(queryString ? `?${queryString}` : '')
		);
	}
</script>

<svelte:head>
	<title>{viewName} | Statusy View</title>
	<meta name="description" content={viewDescription} />
</svelte:head>

<div class="mx-auto">
	<!-- Header and Subtitle -->
	<div class="mb-6">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<h1 class="text-3xl font-extrabold tracking-tight text-white sm:text-4xl">
					{viewName}
				</h1>

				{#if isDefaultView}
					<span
						class="rounded border border-zinc-700/50 bg-zinc-800/80 px-2 py-0.5 text-[10px] font-medium tracking-wider text-zinc-400 uppercase"
					>
						Default
					</span>
				{/if}
			</div>

			<div class="flex items-center gap-2.5">
				<a
					href="/views/{data.view.public_id}/notifications"
					class="inline-flex h-8 cursor-pointer items-center justify-center gap-1.5 rounded-lg border border-zinc-800 bg-zinc-900/50 px-3 text-xs font-medium text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-white"
					title="Configure Notifications"
					aria-label="Configure notifications"
				>
					<Bell class="size-3.5" />
					<span>Configure Notifications</span>
				</a>

				<DropdownMenu.Root>
					<DropdownMenu.Trigger
						class="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900/50 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-white"
						title="View Actions"
						aria-label="View actions"
					>
						<Settings class="size-4" />
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="end" class="min-w-45 border-zinc-800 bg-zinc-950 text-white">
						<DropdownMenu.Item
							onclick={openEditViewDialog}
							class="cursor-pointer hover:bg-zinc-900/50"
						>
							<Pencil class="mr-2 size-3.5" />
							<span>Edit View</span>
						</DropdownMenu.Item>
						<DropdownMenu.Separator class="bg-zinc-900" />
						<DropdownMenu.Item
							onclick={openDeleteViewDialog}
							variant="destructive"
							class="cursor-pointer hover:bg-red-950/20"
						>
							<Trash2 class="mr-2 size-3.5 text-red-500" />
							<span>Delete View</span>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</div>
		<p class="mt-2 max-w-2xl text-zinc-400">
			{viewDescription}
		</p>

		<!-- Concise inline stats -->
		<div class="mt-3.5 flex flex-wrap items-center gap-x-4 gap-y-2 text-xs text-zinc-400">
			<span
				class="flex items-center gap-1.5 rounded-full border border-zinc-800 bg-zinc-900/30 px-2.5 py-1"
			>
				<strong>{totalCount}</strong> services
			</span>
			{#if upCount > 0}
				<span
					class="flex items-center gap-1.5 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2.5 py-1 font-medium text-emerald-400"
				>
					<span class="size-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"
					></span>
					<strong>{upCount}</strong> Up
				</span>
			{/if}
			{#if downCount > 0}
				<span
					class="flex items-center gap-1.5 rounded-full border border-red-500/20 bg-red-500/10 px-2.5 py-1 font-medium text-red-400"
				>
					<span
						class="size-1.5 animate-pulse rounded-full bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.5)]"
					></span>
					<strong>{downCount}</strong> Down
				</span>
			{/if}
		</div>
	</div>

	<!-- Search & Add Bar -->
	<div class="mb-4 flex items-center justify-between gap-3">
		<div class="relative w-full max-w-sm">
			<input
				type="text"
				placeholder="Search configured services..."
				bind:value={searchQuery}
				class="w-full rounded-lg border border-zinc-800 bg-zinc-900/40 py-2 pr-4 pl-9 text-sm text-white placeholder-zinc-500 outline-none focus:border-zinc-700 focus:ring-1 focus:ring-zinc-700"
			/>
			<div class="absolute inset-y-0 left-0 flex items-center pl-3 text-zinc-500">
				<Search class="size-4" />
			</div>
		</div>

		<button
			onclick={navigateToAdd}
			class="flex cursor-pointer items-center gap-1.5 rounded-lg border border-zinc-800 bg-zinc-900/50 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-zinc-800"
		>
			<Plus class="size-4" /> Add Service
		</button>
	</div>

	<!-- Configured Services Table -->
	<div class="overflow-hidden rounded-xl border border-zinc-800 bg-zinc-950/20 shadow-md">
		<Table.Root class="w-full table-auto">
			<Table.Header class="bg-zinc-900/40">
				<Table.Row class="border-zinc-800">
					<Table.Head class="w-[65%] py-3 font-semibold text-zinc-300">Service</Table.Head>
					<Table.Head class="w-[20%] py-3 font-semibold text-zinc-300">Status</Table.Head>
					<Table.Head class="w-[15%] py-3 text-right font-semibold text-zinc-300"
						>Actions</Table.Head
					>
				</Table.Row>
			</Table.Header>

			<Table.Body>
				{#each localServices as service (service.id)}
					<Table.Row
						class="group cursor-pointer border-zinc-800 transition-all duration-200 hover:bg-zinc-900/30"
						onclick={(e) => {
							const target = e.target as HTMLElement;
							if (target.closest('a') || target.closest('button')) {
								return;
							}
							void goto(getEventsUrl(service));
						}}
					>
						<!-- Service & Clickable Status Details -->
						<Table.Cell class="py-3">
							<div class="flex flex-col gap-1.5">
								<div class="flex flex-col">
									<a
										href={getEventsUrl(service)}
										class="font-bold text-white transition-colors hover:text-zinc-300 hover:underline"
									>
										{service.name}
									</a>
								</div>

								<div class="flex flex-col gap-1">
									<!-- Incident/Status info -->
									{#if service.last_incident}
										<a
											href={service.last_incident_link || '#'}
											target="_blank"
											rel="noopener noreferrer"
											class="group/link flex w-fit cursor-pointer items-center gap-1.5 text-xs text-zinc-400 transition-all hover:text-white"
										>
											<span
												class="flex items-center gap-1 font-medium text-red-400 group-hover/link:underline"
												title={service.last_incident}
											>
												<AlertTriangle class="size-3.5 shrink-0 text-red-500" />
												{service.last_incident.length > 50
													? service.last_incident.slice(0, 50) + '...'
													: service.last_incident}
											</span>
										</a>
									{:else}
										<span class="text-xs text-emerald-400">No recent incidents</span>
									{/if}

									<!-- Upcoming Maintenance info -->
									{#if service.monitor_scheduled_maintenances && service.upcoming_maintenance}
										<a
											href={service.upcoming_maintenance_link || '#'}
											target="_blank"
											rel="noopener noreferrer"
											class="group/link flex w-fit cursor-pointer items-center gap-1 text-xs text-blue-400 transition-all hover:text-white"
										>
											<span
												class="flex items-center gap-1 font-medium group-hover/link:underline"
												title={service.upcoming_maintenance}
											>
												<Calendar class="size-3.5 shrink-0 text-blue-500" />
												Maintenance: {service.upcoming_maintenance.length > 50
													? service.upcoming_maintenance.slice(0, 50) + '...'
													: service.upcoming_maintenance}
											</span>
										</a>
									{/if}
								</div>
							</div>
						</Table.Cell>

						<!-- Status Badge -->
						<Table.Cell class="py-3">
							<div class="flex items-center">
								{#if service.status === 'up'}
									<span
										class="flex items-center gap-1.5 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2.5 py-1 text-xs font-semibold text-emerald-400"
									>
										<span
											class="size-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"
										></span>
										Up
									</span>
								{:else}
									<span
										class="flex items-center gap-1.5 rounded-full border border-red-500/20 bg-red-500/10 px-2.5 py-1 text-xs font-semibold text-red-400"
									>
										<span
											class="size-1.5 animate-pulse rounded-full bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.5)]"
										></span>
										Down
									</span>
								{/if}
							</div>
						</Table.Cell>

						<!-- Actions Column -->
						<Table.Cell class="py-3 text-right">
							<div class="flex items-center justify-end gap-1.5">
								<button
									onclick={() => navigateToEdit(service.slug)}
									class="hover:bg-zinc-855 inline-flex size-7 cursor-pointer items-center justify-center rounded-md border border-zinc-800 bg-zinc-900/50 text-zinc-400 transition-all hover:text-white"
									title="Edit Subscribed Components"
								>
									<Pencil class="size-3.5" />
								</button>
								<button
									onclick={() => openDeleteConfirm(service)}
									class="inline-flex size-7 cursor-pointer items-center justify-center rounded-md border border-red-500/20 bg-red-950/20 text-red-400 transition-all hover:bg-red-900/40 hover:text-red-300"
									title="Remove Service"
								>
									<Trash2 class="size-3.5" />
								</button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={3} class="h-28 text-center text-zinc-500">
							No services configured in this view.
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>

	<!-- Pagination Controls (Always visible) -->
	<div
		class="mt-4 flex items-center justify-between border-t border-zinc-900 pt-4 text-sm text-zinc-400"
	>
		<div>
			Showing <span class="font-medium text-white"
				>{totalServicesCount === 0 ? 0 : (currentPage - 1) * itemsPerPage + 1}</span
			>
			to
			<span class="font-medium text-white"
				>{Math.min(currentPage * itemsPerPage, totalServicesCount)}</span
			>
			of
			<span class="font-medium text-white">{totalServicesCount}</span> services
		</div>
		<div class="flex gap-2">
			<Button
				variant="outline"
				size="sm"
				class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
				disabled={currentPage === 1 || totalPages <= 1}
				onclick={() => (currentPage -= 1)}
			>
				Previous
			</Button>
			<Button
				variant="outline"
				size="sm"
				class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
				disabled={currentPage === totalPages || totalPages <= 1}
				onclick={() => (currentPage += 1)}
			>
				Next
			</Button>
		</div>
	</div>
</div>

<!-- Delete Confirmation Dialog -->
<Dialog.Root
	open={isDeleteConfirmOpen}
	onOpenChange={(open) => {
		if (!open) {
			isDeleteConfirmOpen = false;
			serviceToDelete = null;
		}
	}}
>
	<Dialog.Content class="border-zinc-800 bg-zinc-950 text-white shadow-xl sm:max-w-[400px]">
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-2 text-lg font-bold text-white">
				<AlertTriangle class="size-5 text-red-500" />
				Remove Service?
			</Dialog.Title>
			<Dialog.Description class="mt-2 text-zinc-400">
				Are you sure you want to remove <span class="font-bold text-white"
					>{serviceToDelete?.name}</span
				> from this view? You can add it back later.
			</Dialog.Description>
		</Dialog.Header>

		<Dialog.Footer class="mt-4 gap-2">
			<Button
				variant="outline"
				class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
				onclick={() => {
					isDeleteConfirmOpen = false;
					serviceToDelete = null;
				}}
			>
				Cancel
			</Button>
			<Button class="cursor-pointer bg-red-600 text-white hover:bg-red-500" onclick={confirmRemove}>
				Remove
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Edit View Dialog -->
<Dialog.Root
	open={isEditViewOpen}
	onOpenChange={(open) => {
		if (!open) {
			isEditViewOpen = false;
		}
	}}
>
	<Dialog.Content class="border-zinc-800 bg-zinc-950 text-white shadow-xl sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title class="text-lg font-bold text-white">Edit View</Dialog.Title>
			<Dialog.Description class="text-zinc-400">
				Make changes to the view title, description, and settings. Click save when you're done.
			</Dialog.Description>
		</Dialog.Header>

		<ViewForm
			bind:name={editViewName}
			bind:description={editViewDescription}
			bind:isDefault={editViewIsDefault}
			showDefaultCheckbox={true}
			submitText="Save Changes"
			cancelText="Cancel"
			namePlaceholder="Payment Gateways"
			descriptionPlaceholder="Track statuses of Payment APIs and web portals"
			onsubmit={saveViewMeta}
			oncancel={() => {
				isEditViewOpen = false;
			}}
		/>
	</Dialog.Content>
</Dialog.Root>

<!-- Delete View Dialog -->
<Dialog.Root
	open={isDeleteViewOpen}
	onOpenChange={(open) => {
		if (!open) {
			isDeleteViewOpen = false;
		}
	}}
>
	<Dialog.Content class="border-zinc-800 bg-zinc-950 text-white shadow-xl sm:max-w-[400px]">
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-2 text-lg font-bold text-white">
				<AlertTriangle class="size-5 text-red-500" />
				Delete View?
			</Dialog.Title>
			<Dialog.Description class="mt-2 text-zinc-400">
				Are you sure you want to delete the view <span class="font-bold text-white">{viewName}</span
				>? This action will remove all configured services for this view and cannot be undone.
			</Dialog.Description>
		</Dialog.Header>

		<Dialog.Footer class="mt-4 gap-2">
			<Button
				variant="outline"
				class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
				onclick={() => {
					isDeleteViewOpen = false;
				}}
			>
				Cancel
			</Button>
			<Button
				class="cursor-pointer bg-red-600 text-white hover:bg-red-500"
				onclick={confirmDeleteView}
			>
				Delete View
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

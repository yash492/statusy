<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Table from '$lib/components/ui/table';
	import { Textarea } from '$lib/components/ui/textarea';
	import AlertTriangle from '@lucide/svelte/icons/alert-triangle';
	import CheckCircle from '@lucide/svelte/icons/check-circle';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { onMount } from 'svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const view = $derived(data.view);

	// View Settings and Meta State
	let viewName = $state(data.view.name);
	let viewDescription = $state(data.view.description);
	let isDefaultView = $state(false);

	// Local reactive services state synced with localStorage
	let localServices = $state<any[]>([]);

	onMount(() => {
		const stored = localStorage.getItem(`statusy_view_${data.view.slug}`);
		if (stored) {
			localServices = JSON.parse(stored);
		} else {
			localServices = [...data.view.services];
			localStorage.setItem(`statusy_view_${data.view.slug}`, JSON.stringify(localServices));
		}

		const storedMeta = localStorage.getItem(`statusy_view_meta_${data.view.slug}`);
		if (storedMeta) {
			const meta = JSON.parse(storedMeta);
			viewName = meta.name;
			viewDescription = meta.description;
		}

		const defaultSlug = localStorage.getItem('statusy_default_view_slug');
		isDefaultView = (defaultSlug === data.view.slug);
	});

	$effect(() => {
		const stored = localStorage.getItem(`statusy_view_${data.view.slug}`);
		if (stored) {
			localServices = JSON.parse(stored);
		} else {
			localServices = [...data.view.services];
			localStorage.setItem(`statusy_view_${data.view.slug}`, JSON.stringify(localServices));
		}

		const storedMeta = localStorage.getItem(`statusy_view_meta_${data.view.slug}`);
		if (storedMeta) {
			const meta = JSON.parse(storedMeta);
			viewName = meta.name;
			viewDescription = meta.description;
		} else {
			viewName = data.view.name;
			viewDescription = data.view.description;
		}

		const defaultSlug = localStorage.getItem('statusy_default_view_slug');
		isDefaultView = (defaultSlug === data.view.slug);
	});

	// Compute metrics dynamically from localState
	const totalCount = $derived(localServices.length);
	const upCount = $derived(localServices.filter((s) => s.status === 'operational').length);
	const downCount = $derived(localServices.filter((s) => s.status !== 'operational').length);

	// Search and Pagination State
	let searchQuery = $state('');
	let currentPage = $state(1);
	const itemsPerPage = 5;

	// Filtered services based on search query
	const filteredServices = $derived(
		localServices.filter(
			(service) =>
				service.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				service.slug.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	// Paginated services based on filtered list
	const totalPages = $derived(Math.ceil(filteredServices.length / itemsPerPage) || 1);
	const paginatedServices = $derived(
		filteredServices.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
	);

	// Reset to page 1 if search query changes
	$effect(() => {
		if (searchQuery) {
			currentPage = 1;
		}
	});

	// Deletion State
	let isDeleteConfirmOpen = $state(false);
	let serviceToDelete = $state<any | null>(null);

	function openDeleteConfirm(service: any) {
		serviceToDelete = service;
		isDeleteConfirmOpen = true;
	}

	function confirmRemove() {
		if (serviceToDelete) {
			const updated = localServices.filter((s) => s.id !== serviceToDelete.id);
			localServices = updated;
			localStorage.setItem(`statusy_view_${data.view.slug}`, JSON.stringify(updated));
			isDeleteConfirmOpen = false;
			serviceToDelete = null;
		}
	}

	function navigateToAdd() {
		void goto(`/views/${data.view.slug}/add-service`);
	}

	function navigateToEdit(serviceSlug: string) {
		void goto(`/views/${data.view.slug}/edit-service/${serviceSlug}`);
	}

	// Edit View Dialog State
	let isEditViewOpen = $state(false);
	let editViewName = $state('');
	let editViewDescription = $state('');
	let editViewIsDefault = $state(false);
	let isMenuOpen = $state(false);

	function openEditViewDialog() {
		editViewName = viewName;
		editViewDescription = viewDescription;
		editViewIsDefault = isDefaultView;
		isEditViewOpen = true;
	}

	function saveViewMeta() {
		localStorage.setItem(`statusy_view_meta_${data.view.slug}`, JSON.stringify({
			name: editViewName,
			description: editViewDescription
		}));
		viewName = editViewName;
		viewDescription = editViewDescription;

		if (editViewIsDefault) {
			localStorage.setItem('statusy_default_view_slug', data.view.slug);
			isDefaultView = true;
		} else {
			if (localStorage.getItem('statusy_default_view_slug') === data.view.slug) {
				localStorage.removeItem('statusy_default_view_slug');
			}
			isDefaultView = false;
		}
		isEditViewOpen = false;
	}

	// Delete View Dialog State
	let isDeleteViewOpen = $state(false);

	function openDeleteViewDialog() {
		isDeleteViewOpen = true;
	}

	function confirmDeleteView() {
		localStorage.removeItem(`statusy_view_${data.view.slug}`);
		localStorage.removeItem(`statusy_view_meta_${data.view.slug}`);
		if (localStorage.getItem('statusy_default_view_slug') === data.view.slug) {
			localStorage.removeItem('statusy_default_view_slug');
		}
		isDeleteViewOpen = false;
		void goto('/');
	}
</script>

<svelte:head>
	<title>{viewName} | Statusy View</title>
	<meta name="description" content={viewDescription} />
</svelte:head>

<div class="mx-auto w-3/5 pt-4">
	<!-- Header and Subtitle -->
	<div class="mb-6">
		<div class="flex items-center gap-3">
			<h1 class="text-3xl font-extrabold tracking-tight text-white sm:text-4xl">
				{viewName}
			</h1>

			{#if isDefaultView}
				<span class="rounded border border-zinc-700/50 bg-zinc-800/80 px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider text-zinc-400">
					Default
				</span>
			{/if}

			<button
				onclick={openEditViewDialog}
				class="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900/50 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-white"
				title="Edit View Settings"
				aria-label="Edit view settings"
			>
				<Pencil class="size-4" />
			</button>

			<button
				onclick={openDeleteViewDialog}
				class="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-red-500/20 bg-red-950/20 text-red-400 transition-colors hover:bg-red-900/40 hover:text-red-300"
				title="Delete View"
				aria-label="Delete view"
			>
				<Trash2 class="size-4" />
			</button>
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
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="size-4"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.602 10.602Z"
					/>
				</svg>
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
				{#each paginatedServices as service (service.id)}
					<Table.Row class="group border-zinc-800 transition-all duration-200 hover:bg-zinc-900/30">
						<!-- Service & Clickable Incident Detail -->
						<Table.Cell class="py-3">
							<div class="flex flex-col gap-1.5">
								<div class="flex flex-col">
									<span class="font-bold text-white transition-colors group-hover:text-zinc-200">
										{service.name}
									</span>
								</div>

								<!-- Clickable Recent Incident info directly under Service Info -->
								<a
									href={`/statuspages/${service.slug}/events`}
									class="group/link flex w-fit cursor-pointer items-center gap-2 text-xs text-zinc-400 transition-all hover:text-white"
								>
									{#if service.lastIncident}
										<span
											class="flex items-center gap-1 font-medium text-amber-400 group-hover/link:underline"
											title={service.lastIncident}
										>
											<AlertTriangle class="size-3.5 shrink-0 text-amber-500" />
											{service.lastIncident}
										</span>
									{:else}
										<span class="flex items-center gap-1 text-zinc-500 group-hover/link:underline">
											<CheckCircle class="size-3.5 shrink-0 text-emerald-500/70" />
											No recent incidents
										</span>
									{/if}
								</a>
							</div>
						</Table.Cell>

						<!-- Status Badge -->
						<Table.Cell class="py-3">
							<div class="flex items-center">
								{#if service.status === 'operational'}
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
				>{filteredServices.length === 0 ? 0 : (currentPage - 1) * itemsPerPage + 1}</span
			>
			to
			<span class="font-medium text-white"
				>{Math.min(currentPage * itemsPerPage, filteredServices.length)}</span
			>
			of
			<span class="font-medium text-white">{filteredServices.length}</span> services
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
			<Dialog.Title class="text-lg font-bold text-white">Edit View Settings</Dialog.Title>
			<Dialog.Description class="text-zinc-400">
				Make changes to the view title, description, and settings. Click save when you're done.
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="view-name" class="text-sm font-semibold text-zinc-300">Name</Label>
				<Input
					id="view-name"
					bind:value={editViewName}
					placeholder="Payment Gateways"
					class="border-zinc-800 bg-zinc-900/50 text-white placeholder-zinc-500 focus-visible:ring-zinc-700"
				/>
			</div>
			<div class="grid gap-2">
				<Label for="view-description" class="text-sm font-semibold text-zinc-300">Description</Label>
				<Textarea
					id="view-description"
					bind:value={editViewDescription}
					placeholder="Track statuses of Payment APIs and web portals"
					class="border-zinc-800 bg-zinc-900/50 text-white placeholder-zinc-500 focus-visible:ring-zinc-700"
				/>
			</div>
			
			<label class="mt-2 flex cursor-pointer items-center gap-2 select-none">
				<input
					type="checkbox"
					bind:checked={editViewIsDefault}
					class="size-4 rounded border-zinc-800 bg-zinc-900 text-white accent-emerald-500 transition-colors focus:ring-0"
				/>
				<span class="text-sm text-zinc-300">Make this the default view</span>
			</label>
		</div>

		<Dialog.Footer class="gap-2">
			<Button
				variant="outline"
				class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
				onclick={() => {
					isEditViewOpen = false;
				}}
			>
				Cancel
			</Button>
			<Button
				class="cursor-pointer bg-zinc-100 text-zinc-950 hover:bg-zinc-200"
				onclick={saveViewMeta}
			>
				Save Changes
			</Button>
		</Dialog.Footer>
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
				Are you sure you want to delete the view <span class="font-bold text-white">{viewName}</span>? This action will remove all configured services for this view and cannot be undone.
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

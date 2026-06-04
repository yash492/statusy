<script lang="ts">
	import { page } from '$app/state';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import ViewForm from '$lib/components/ViewForm.svelte';
	import { ChevronsUpDown, LayoutDashboard, Plus } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	const defaultView = $derived(page.data.defaultView);
	const viewUrl = $derived(defaultView ? `/views/${defaultView.public_id}` : '/');
	const viewName = $derived(defaultView?.name ?? 'Default View');
	const isWorkspaceActive = $derived(
		(page.url.pathname as string).startsWith('/views/') ||
			(page.url.pathname as string).startsWith('/statuspages/')
	);

	// Views management modal state
	let isViewsModalOpen = $state(false);
	let searchQuery = $state('');

	// Add view modal state
	let isAddingView = $state(false);
	let newName = $state('');
	let newDescription = $state('');

	let mockViews = $state([
		{
			id: 1,
			public_id: 'default-view',
			name: 'Default View',
			description: 'Main view tracking all services'
		},
		{
			id: 2,
			public_id: 'payment-gateways',
			name: 'Payment Gateways',
			description: 'Track payment APIs and web portals'
		},
		{
			id: 3,
			public_id: 'internal-infrastructure',
			name: 'Internal Infra',
			description: 'Internal databases, cache, and servers'
		}
	]);

	// Filtered and limited views (shows top 5 results)
	const filteredViews = $derived(
		mockViews.filter(
			(v) =>
				v.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				v.description.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	const limitViews = $derived(filteredViews.slice(0, 5));

	function addView() {
		if (!newName.trim()) {
			toast.error('View name is required');
			return;
		}
		const newId = mockViews.length > 0 ? Math.max(...mockViews.map((v) => v.id)) + 1 : 1;
		const newSlug = newName
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/(^-|-$)/g, '');
		mockViews.push({
			id: newId,
			public_id: newSlug,
			name: newName,
			description: newDescription
		});
		toast.success(`Created view "${newName}"`);
		newName = '';
		newDescription = '';
		isAddingView = false;
	}
</script>

<Sidebar.Root>
	<Sidebar.Header class="border-b border-zinc-900/50 px-5 py-6">
		<a href={viewUrl} class="flex items-center gap-3 select-none">
			<div
				class="flex size-9 items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900 shadow-[0_0_12px_rgba(255,255,255,0.05)]"
			>
				<img src="/logos/statusy.svg" class="size-5" alt="Statusy Logo" />
			</div>
			<span
				class="bg-linear-to-r from-white via-zinc-200 to-zinc-400 bg-clip-text text-xl font-bold tracking-tight text-transparent"
				>Statusy</span
			>
		</a>
	</Sidebar.Header>

	<Sidebar.Content class="flex flex-col gap-4 px-3 py-4">
		<Sidebar.Group>
			<Sidebar.GroupLabel
				class="px-2 pt-4 text-[10px] font-bold tracking-wider text-zinc-500 uppercase select-none"
				>Views</Sidebar.GroupLabel
			>
			<Sidebar.GroupContent class="mt-1.5">
				<Sidebar.Menu>
					<Sidebar.MenuItem class="pt-3">
						<!-- Main View button -->
						<Sidebar.MenuButton
							isActive={isWorkspaceActive}
							class="rounded-lg px-3 py-2 text-zinc-400 transition-all duration-150 hover:bg-zinc-900/50 hover:text-white"
							onclick={() => {
								isViewsModalOpen = true;
							}}
						>
							{#snippet child({ props })}
								<button
									class="group flex w-full cursor-pointer items-center justify-between border-0 bg-transparent p-0 text-left text-zinc-400 hover:text-white"
									{...props}
								>
									<div class="flex items-center gap-2.5 text-sm font-semibold">
										<LayoutDashboard class="size-4" />
										<span>{viewName}</span>
									</div>
									<ChevronsUpDown
										class="size-3.5 text-zinc-500 transition-colors group-hover:text-zinc-300"
									/>
								</button>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>

	<!-- Footer with Github Link and Sleek Layout -->
	<Sidebar.Footer class="flex items-center justify-between border-t border-zinc-900/50 px-5 py-4">
		<a
			href="https://github.com"
			target="_blank"
			rel="noopener noreferrer"
			class="flex items-center gap-2 text-xs font-semibold text-zinc-500 transition-colors hover:text-zinc-300"
		>
			<svg
				class="size-4 shrink-0"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				><path
					d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"
				/></svg
			>
			<span>GitHub</span>
		</a>
		<span class="text-[10px] font-medium text-zinc-600 select-none">v2.0.0</span>
	</Sidebar.Footer>
</Sidebar.Root>

<!-- Views list dialog -->
<Dialog.Root
	open={isViewsModalOpen}
	onOpenChange={(open) => {
		isViewsModalOpen = open;
		if (!open) {
			isAddingView = false;
			newName = '';
			newDescription = '';
		}
	}}
>
	<Dialog.Content
		class="rounded-xl border-zinc-800 bg-zinc-950 p-6 text-white shadow-xl sm:max-w-120"
	>
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-3 text-lg font-bold text-white">
				<span>Views</span>
				{#if !isAddingView}
					<button
						onclick={() => (isAddingView = true)}
						class="flex cursor-pointer items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900 p-1.5 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-white"
						title="Add View"
					>
						<Plus class="size-3.5" />
					</button>
				{/if}
			</Dialog.Title>
		</Dialog.Header>

		{#if isAddingView}
			<!-- Add View Form -->
			<div class="mt-4 space-y-3.5 rounded-lg border border-zinc-800 bg-zinc-900/20 p-4">
				<h4 class="text-xs font-bold tracking-wider text-zinc-300 uppercase">New View</h4>
				<ViewForm
					bind:name={newName}
					bind:description={newDescription}
					namePlaceholder="e.g. API Services"
					descriptionPlaceholder="Describe the services in this view"
					submitText="Create"
					cancelText="Cancel"
					onsubmit={addView}
					oncancel={() => (isAddingView = false)}
				/>
			</div>
		{:else}
			<!-- Search Bar -->
			<div class="relative mt-4">
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search views..."
					class="w-full rounded-md border border-zinc-800 bg-zinc-900/50 px-3 py-2 text-xs text-white placeholder-zinc-500 outline-none focus:border-zinc-700"
				/>
			</div>

			<!-- Views list (max 5 results) -->
			<div class="mt-4 space-y-2">
				{#each limitViews as view (view.id)}
					<a
						href={`/views/${view.public_id}`}
						onclick={() => {
							isViewsModalOpen = false;
						}}
						class="flex items-start justify-between gap-3 rounded-lg border border-zinc-900 bg-zinc-950/40 p-3 text-left transition-colors hover:border-zinc-800 hover:bg-zinc-900/40"
					>
						<div class="min-w-0 flex-1">
							<div class="flex items-center gap-1.5">
								<h4 class="truncate text-xs font-bold text-white">{view.name}</h4>
								{#if view.public_id === defaultView?.public_id}
									<span
										class="py-0.2 rounded border border-zinc-700/55 bg-zinc-800 px-1 text-[8px] font-medium tracking-wider text-zinc-400 uppercase select-none"
										>Default</span
									>
								{/if}
							</div>
							<p class="text-zinc-550 mt-0.5 truncate text-[10px]">
								{view.description || 'No description'}
							</p>
						</div>
					</a>
				{:else}
					<div
						class="py-8 text-center text-xs text-zinc-500 border border-zinc-900 bg-zinc-950/20 rounded-lg"
					>
						No views match your search.
					</div>
				{/each}
			</div>
		{/if}

		<Dialog.Footer class="mt-6 flex justify-end border-t border-zinc-900 pt-4">
			<button
				onclick={() => {
					isViewsModalOpen = false;
				}}
				class="border-zinc-850 cursor-pointer rounded-lg border bg-zinc-900 px-4 py-2 text-xs font-semibold text-white transition-colors hover:bg-zinc-800"
			>
				Close
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

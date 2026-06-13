<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { ViewsApi, type View } from '$lib/api/views/views';
	import ViewForm from '$lib/components/statusy/ViewForm.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { ChevronsUpDown, Plus } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	const viewsApi = new ViewsApi();

	const defaultView = $derived(page.data.defaultView);
	const viewUrl = $derived(defaultView ? `/views/${defaultView.public_id}` : '/');
	const activeViewName = $derived(page.data.view?.name ?? defaultView?.name ?? 'Default View');

	// Views management modal state
	let isViewsModalOpen = $state(false);

	// Add view modal state
	let isAddViewDialogOpen = $state(false);
	let newName = $state('');
	let newDescription = $state('');
	let searchQuery = $state('');

	let viewsList = $state<View[]>([]);
	let totalCount = $state(0);
	const limitViews = $derived(viewsList);

	$effect(() => {
		// Reactive dependency on page.data.views and searchQuery
		const _trigger = page.data.views;
		const q = searchQuery;

		viewsApi.list(q).then(([res, err]) => {
			if (!err && res) {
				viewsList = res.views;
				totalCount = res.total_count;
			}
		});
	});

	$effect(() => {
		if (!isViewsModalOpen) {
			searchQuery = '';
		}
	});

	$effect(() => {
		if (!isAddViewDialogOpen) {
			newName = '';
			newDescription = '';
		}
	});

	async function addView() {
		if (!newName.trim()) {
			toast.error('View name is required');
			return;
		}

		const [view, err] = await viewsApi.create({
			name: newName,
			description: newDescription
		});

		if (err) {
			toast.error(err.message || 'Failed to create view');
			return;
		}

		toast.success(`Created view "${newName}"`);
		newName = '';
		newDescription = '';
		isAddViewDialogOpen = false;
		await invalidateAll();
		if (view?.public_id) {
			await goto(`/views/${view.public_id}`);
		}
	}
</script>

<header class="sticky top-0 z-40 w-full border-b border-zinc-900 bg-zinc-950/80 backdrop-blur-md">
	<div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
		<!-- Left: Logo & Selector -->
		<div class="flex items-center gap-6">
			<!-- Logo -->
			<a href={viewUrl} class="flex items-center gap-3 select-none">
				<div
					class="flex size-10 items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900 shadow-[0_0_12px_rgba(255,255,255,0.05)]"
				>
					<img src="/logos/statusy.svg" class="size-6" alt="Statusy Logo" />
				</div>
				<span class="text-2xl font-bold tracking-tight text-white">Statusy</span>
			</a>

			<!-- Separator -->
			<div class="h-5 w-px bg-zinc-800"></div>

			<!-- View Selector Button -->
			<button
				class="group flex cursor-pointer items-center gap-2 rounded-lg border border-zinc-900 bg-zinc-950/40 px-3 py-1.5 text-left text-sm font-semibold text-zinc-400 transition-all hover:border-zinc-800 hover:text-white"
				onclick={() => {
					isViewsModalOpen = true;
				}}
			>
				<span>{activeViewName}</span>
				<ChevronsUpDown class="size-4 text-zinc-500 transition-colors group-hover:text-zinc-300" />
			</button>
		</div>

		<!-- Right: GitHub Link -->
		<div>
			<a
				href="https://github.com/yash492/statusy"
				target="_blank"
				rel="noopener noreferrer"
				class="flex items-center gap-2 text-xs font-semibold text-zinc-500 transition-colors hover:text-zinc-300"
			>
				<img src="/logos/github.svg" class="size-4 shrink-0" alt="GitHub" />
				<span>GitHub</span>
			</a>
		</div>
	</div>
</header>

<!-- Views list dialog -->
<Dialog.Root bind:open={isViewsModalOpen}>
	<Dialog.Content
		class="rounded-xl border-zinc-800 bg-zinc-950 p-5 text-white shadow-xl sm:max-w-96"
	>
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-3 text-lg font-bold text-white">
				<span>Views ({totalCount})</span>
				<button
					onclick={() => {
						isViewsModalOpen = false;
						isAddViewDialogOpen = true;
					}}
					class="flex cursor-pointer items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900 p-1 text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-white"
					title="Add View"
				>
					<Plus class="size-3" />
				</button>
			</Dialog.Title>
		</Dialog.Header>

		<!-- Search Bar -->
		<div class="relative mt-3">
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search views..."
				class="placeholder-zinc-555 w-full rounded-md border border-zinc-800 bg-zinc-900/50 px-3 py-1.5 text-base text-white outline-none focus:border-zinc-700"
			/>
		</div>

		<!-- Views list -->
		<div class="mt-3 space-y-1.5">
			{#each limitViews as view (view.public_id)}
				<a
					href={`/views/${view.public_id}`}
					onclick={() => {
						isViewsModalOpen = false;
					}}
					class="flex items-start justify-between gap-3 rounded-lg border border-zinc-900 bg-zinc-950/40 p-2 text-left transition-colors hover:border-zinc-800 hover:bg-zinc-900/40"
				>
					<div class="min-w-0 flex-1">
						<div class="flex items-center gap-1.5">
							<h4 class="truncate text-base font-bold text-white">{view.name}</h4>
							{#if view.is_default}
								<span
									class="py-0.2 rounded border border-zinc-700/55 bg-zinc-800 px-1.5 text-[10px] font-medium tracking-wider text-zinc-400 uppercase select-none"
									>Default</span
								>
							{/if}
						</div>
						<p class="mt-0.5 truncate text-xs text-zinc-500">
							{view.description || 'No description'}
						</p>
					</div>
				</a>
			{:else}
				<div
					class="py-6 text-center text-base text-zinc-500 border border-zinc-900 bg-zinc-950/20 rounded-lg"
				>
					No views match your search.
				</div>
			{/each}

			{#if totalCount > viewsList.length}
				<p class="mt-2 text-center text-[11px] text-zinc-500 select-none">
					Showing {viewsList.length} of {totalCount} views. Search to filter.
				</p>
			{/if}
		</div>

		<Dialog.Footer class="mt-5 flex justify-end border-t border-zinc-900 pt-3">
			<button
				onclick={() => {
					isViewsModalOpen = false;
				}}
				class="border-zinc-850 cursor-pointer rounded-lg border bg-zinc-900 px-3 py-1.5 text-base font-semibold text-white transition-colors hover:bg-zinc-800"
			>
				Close
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Add view dialog popup -->
<Dialog.Root bind:open={isAddViewDialogOpen}>
	<Dialog.Content class="border-zinc-800 bg-zinc-950 text-white shadow-xl sm:max-w-120">
		<Dialog.Header>
			<Dialog.Title class="text-xl font-bold text-white">New View</Dialog.Title>
		</Dialog.Header>

		<ViewForm
			bind:name={newName}
			bind:description={newDescription}
			namePlaceholder="e.g. API Services"
			descriptionPlaceholder="Describe the services in this view"
			submitText="Create"
			cancelText="Cancel"
			onsubmit={addView}
			oncancel={() => (isAddViewDialogOpen = false)}
		/>
	</Dialog.Content>
</Dialog.Root>

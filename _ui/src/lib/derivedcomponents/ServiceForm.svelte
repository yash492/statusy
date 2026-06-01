<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';

	let {
		mode,
		publicId,
		serviceSlug
	}: { mode: 'add' | 'edit'; publicId: string; serviceSlug?: string } = $props();

	// Mock available services for select dropdown
	let availableServices = $state([
		{ id: 1, name: 'Stripe', slug: 'stripe' },
		{ id: 2, name: 'GitHub', slug: 'github' },
		{ id: 3, name: 'Cloudflare', slug: 'cloudflare' },
		{ id: 4, name: 'PostgreSQL', slug: 'postgresql' }
	]);

	let selectedServiceIdStr = $state('');
	let selectedServiceId = $derived(selectedServiceIdStr ? Number(selectedServiceIdStr) : null);
	let selectOpen = $state(false);
	let serviceSearchQuery = $state('');
	let componentMode = $state<'all' | 'custom'>('all');
	let errorMessage = $state<string | null>(null);

	// Mock component structure for the selected service
	let configuringService = $state({
		service_name: 'Stripe',
		grouped_components: [
			{
				id: 10,
				name: 'API',
				components: [
					{ id: 101, name: 'v3 API' },
					{ id: 102, name: 'Checkout' }
				]
			},
			{
				id: 11,
				name: 'Dashboard',
				components: [
					{ id: 103, name: 'Billing Panel' },
					{ id: 104, name: 'Developer Logs' }
				]
			}
		],
		ungrouped_components: [
			{ id: 105, name: 'Webhooks' },
			{ id: 106, name: 'Support Portal' }
		]
	});

	let selectedComponentIds = $state<number[]>([101, 102, 103, 104, 105, 106]);
	let selectedComponentGroupIds = $state<number[]>([10, 11]);

	let activeService = $state({
		name: 'Stripe'
	});

	const selectedServiceLabel = $derived.by(() => {
		if (mode === 'edit') return activeService.name;
		const id = Number(selectedServiceIdStr);
		const service = availableServices.find((s) => s.id === id);
		return service ? service.name : '';
	});

	function autofocus(node: HTMLElement) {
		requestAnimationFrame(() => {
			node.focus();
		});
	}

	function selectAll() {
		selectedComponentGroupIds = configuringService.grouped_components.map((g) => g.id);
		const ids = configuringService.grouped_components.flatMap((g) => g.components.map((c) => c.id));
		ids.push(...configuringService.ungrouped_components.map((c) => c.id));
		selectedComponentIds = ids;
	}

	function toggleComponent(componentId: number, groupId?: number) {
		if (selectedComponentIds.includes(componentId)) {
			selectedComponentIds = selectedComponentIds.filter((id) => id !== componentId);
			if (groupId) {
				selectedComponentGroupIds = selectedComponentGroupIds.filter((id) => id !== groupId);
			}
		} else {
			selectedComponentIds = [...selectedComponentIds, componentId];
			if (groupId) {
				const group = configuringService.grouped_components.find((g) => g.id === groupId);
				if (group) {
					const allChecked = group.components.every((c) => selectedComponentIds.includes(c.id));
					if (allChecked) {
						selectedComponentGroupIds = [...selectedComponentGroupIds, groupId];
					}
				}
			}
		}
	}

	function toggleGroup(group: any) {
		const allChecked = group.components.every((c: any) => selectedComponentIds.includes(c.id));
		if (allChecked) {
			selectedComponentGroupIds = selectedComponentGroupIds.filter((id) => id !== group.id);
			selectedComponentIds = selectedComponentIds.filter(
				(id) => !group.components.some((c: any) => c.id === id)
			);
		} else {
			selectedComponentGroupIds = [...selectedComponentGroupIds, group.id];
			const newIds = group.components
				.map((c: any) => c.id)
				.filter((id: number) => !selectedComponentIds.includes(id));
			selectedComponentIds = [...selectedComponentIds, ...newIds];
		}
	}

	function saveService() {
		void goto(`/views/${publicId}`);
	}

	function handleCancel() {
		void goto(`/views/${publicId}`);
	}
</script>

<svelte:head>
	<title>{mode === 'add' ? 'Add Service' : 'Edit Service'} | Statusy</title>
</svelte:head>

<div class="mx-auto w-3/5 pt-4 text-white">
	<!-- Header -->
	<div class="mb-6 max-w-xl">
		<h1 class="text-3xl font-extrabold tracking-tight text-white sm:text-4xl">
			{mode === 'add' ? 'Add Service' : 'Edit Service'}
		</h1>
		<p class="mt-2 text-sm text-zinc-400">
			{mode === 'add'
				? 'Choose a service and select which components to subscribe to for this dashboard.'
				: 'Customize the component subscriptions for this service dashboard view.'}
		</p>
	</div>

	<!-- Form Card (Borderless & Background Merged) -->
	<div class="max-w-xl">
		{#if errorMessage}
			<div class="mb-4 rounded-lg border border-red-500/20 bg-red-950/20 p-3 text-sm text-red-400">
				{errorMessage}
			</div>
		{/if}

		{#if mode === 'add' || activeService}
			<div class="grid gap-6">
				<!-- Service Name Input/Dropdown Selector -->
				<div class="grid gap-2">
					<Label for="service-select" class="text-sm font-semibold text-zinc-300"
						>Service Name</Label
					>
					{#if mode === 'add'}
						<Select.Root type="single" bind:value={selectedServiceIdStr} bind:open={selectOpen}>
							<Select.Trigger
								id="service-select"
								class="flex h-12 w-full items-center justify-between rounded-lg border border-zinc-800 bg-zinc-900/50 px-4 text-sm text-white placeholder-zinc-500 transition-colors outline-none hover:bg-zinc-900/80 focus:bg-zinc-900/80 focus:ring-1 focus:ring-zinc-700"
							>
								{selectedServiceLabel || 'Select a service to add...'}
							</Select.Trigger>
							<Select.Content class="border border-zinc-800 bg-zinc-950 p-2 text-white">
								<div class="px-2 py-1.5">
									<input
										type="text"
										placeholder="Search service..."
										use:autofocus
										bind:value={serviceSearchQuery}
										class="w-full rounded-md border border-zinc-800 bg-zinc-900/50 px-3 py-2 text-xs text-white placeholder-zinc-500 outline-none focus:border-zinc-700 focus:ring-1 focus:ring-zinc-700"
										onkeydown={(e) => {
											if (e.key === 'Escape' || e.key === 'Tab') return;
											e.stopPropagation();
										}}
									/>
								</div>
								<div class="my-1 h-px bg-zinc-800/60"></div>
								<div class="max-h-[200px] overflow-y-auto">
									{#if availableServices.length === 0}
										<div class="px-4 py-3 text-xs text-zinc-500">No services found.</div>
									{/if}
									{#each availableServices as service}
										<Select.Item
											value={String(service.id)}
											label={service.name}
											class="cursor-pointer px-4 py-2.5 text-zinc-300 focus:bg-zinc-900 focus:text-white"
										>
											{service.name}
										</Select.Item>
									{/each}
								</div>
							</Select.Content>
						</Select.Root>
					{:else if activeService}
						<div class="rounded-lg bg-zinc-900/50 px-4 py-3 text-sm font-medium text-zinc-200">
							{activeService.name}
						</div>
					{/if}
				</div>

				<!-- Component Checklist with Hierarchical groups -->
				{#if configuringService}
					<div class="grid gap-4 pt-5">
						<Label class="text-sm font-semibold text-zinc-300">Monitored Components</Label>

						<!-- Custom Radio Buttons and Component checklist in border container -->
						<div
							class="flex flex-col gap-4 rounded-lg border border-zinc-800/40 bg-zinc-900/20 p-4"
						>
							<label class="group flex cursor-pointer items-start gap-3">
								<input
									type="radio"
									name="component-mode"
									value="all"
									checked={componentMode === 'all'}
									onchange={() => {
										componentMode = 'all';
										selectAll();
									}}
									class="sr-only"
								/>
								<div
									class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full border border-zinc-700 transition-all group-hover:border-zinc-500 {componentMode ===
									'all'
										? 'border-emerald-500 bg-emerald-500/10'
										: ''}"
								>
									{#if componentMode === 'all'}
										<div class="size-2 rounded-full bg-emerald-500"></div>
									{/if}
								</div>
								<div class="flex flex-col">
									<span
										class="text-sm font-medium text-zinc-200 transition-colors group-hover:text-white"
									>
										Monitor all components
									</span>
									<span class="mt-0.5 text-xs text-zinc-500">
										Automatically monitor all current and future components for this service
									</span>
								</div>
							</label>

							<label class="group flex cursor-pointer items-start gap-3">
								<input
									type="radio"
									name="component-mode"
									value="custom"
									checked={componentMode === 'custom'}
									onchange={() => {
										componentMode = 'custom';
									}}
									class="sr-only"
								/>
								<div
									class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full border border-zinc-700 transition-all group-hover:border-zinc-500 {componentMode ===
									'custom'
										? 'border-emerald-500 bg-emerald-500/10'
										: ''}"
								>
									{#if componentMode === 'custom'}
										<div class="size-2 rounded-full bg-emerald-500"></div>
									{/if}
								</div>
								<div class="flex flex-col">
									<span
										class="text-sm font-medium text-zinc-200 transition-colors group-hover:text-white"
									>
										Customize component selection
									</span>
									<span class="mt-0.5 text-xs text-zinc-500">
										Select specific component groups or individual components to monitor
									</span>
								</div>
							</label>

							{#if componentMode === 'custom'}
								<div class="ml-7 grid max-h-72 grid-cols-1 gap-3 overflow-y-auto pr-1">
									{#each configuringService.grouped_components as group}
										{@const groupChecked = group.components.every((c) =>
											selectedComponentIds.includes(c.id)
										)}
										{@const groupSomeChecked =
											group.components.some((c) => selectedComponentIds.includes(c.id)) &&
											!groupChecked}

										<div class="py-2">
											<!-- Group Header -->
											<button
												type="button"
												onclick={() => toggleGroup(group)}
												class="flex w-full cursor-pointer items-center gap-3 text-left transition-all hover:text-white"
											>
												<div
													class="flex size-4 shrink-0 items-center justify-center rounded transition-all {groupChecked
														? 'bg-emerald-500 text-zinc-950'
														: groupSomeChecked
															? 'bg-emerald-500/20 text-emerald-400'
															: 'bg-zinc-900'}"
												>
													{#if groupChecked}
														<svg
															xmlns="http://www.w3.org/2000/svg"
															fill="none"
															viewBox="0 0 24 24"
															stroke-width="3"
															stroke="currentColor"
															class="size-2.5"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																d="m4.5 12.75 6 6 9-13.5"
															/>
														</svg>
													{:else if groupSomeChecked}
														<div class="size-1.5 rounded-sm bg-emerald-400"></div>
													{/if}
												</div>
												<span class="text-sm font-bold text-zinc-200">{group.name} Group</span>
											</button>

											<!-- Child Components -->
											<div class="mt-2 ml-7 grid gap-2 pl-3">
												{#each group.components as component}
													{@const componentChecked = selectedComponentIds.includes(component.id)}
													<button
														type="button"
														onclick={() => toggleComponent(component.id, group.id)}
														class="flex cursor-pointer items-center gap-3 text-left transition-all hover:text-zinc-200 {componentChecked
															? 'text-zinc-300'
															: 'text-zinc-500'}"
													>
														<div
															class="flex size-3.5 shrink-0 items-center justify-center rounded transition-all {componentChecked
																? 'bg-emerald-500/80 text-zinc-950'
																: 'bg-zinc-900'}"
														>
															{#if componentChecked}
																<svg
																	xmlns="http://www.w3.org/2000/svg"
																	fill="none"
																	viewBox="0 0 24 24"
																	stroke-width="3.5"
																	stroke="currentColor"
																	class="size-2"
																>
																	<path
																		stroke-linecap="round"
																		stroke-linejoin="round"
																		d="m4.5 12.75 6 6 9-13.5"
																	/>
																</svg>
															{/if}
														</div>
														<span class="text-xs font-medium">{component.name}</span>
													</button>
												{/each}
											</div>
										</div>
									{/each}

									{#if configuringService.ungrouped_components.length > 0}
										<div class="py-2">
											<span class="text-sm font-bold text-zinc-200">General Components</span>
											<div class="mt-2 ml-7 grid gap-2 pl-3">
												{#each configuringService.ungrouped_components as component}
													{@const componentChecked = selectedComponentIds.includes(component.id)}
													<button
														type="button"
														onclick={() => toggleComponent(component.id)}
														class="flex cursor-pointer items-center gap-3 text-left transition-all hover:text-zinc-200 {componentChecked
															? 'text-zinc-300'
															: 'text-zinc-500'}"
													>
														<div
															class="flex size-3.5 shrink-0 items-center justify-center rounded transition-all {componentChecked
																? 'bg-emerald-500/80 text-zinc-950'
																: 'bg-zinc-900'}"
														>
															{#if componentChecked}
																<svg
																	xmlns="http://www.w3.org/2000/svg"
																	fill="none"
																	viewBox="0 0 24 24"
																	stroke-width="3.5"
																	stroke="currentColor"
																	class="size-2"
																>
																	<path
																		stroke-linecap="round"
																		stroke-linejoin="round"
																		d="m4.5 12.75 6 6 9-13.5"
																	/>
																</svg>
															{/if}
														</div>
														<span class="text-xs font-medium">{component.name}</span>
													</button>
												{/each}
											</div>
										</div>
									{/if}
								</div>
							{/if}
						</div>
					</div>
				{/if}
			</div>

			<!-- Footer Buttons -->
			<div class="mt-8 flex items-center justify-end gap-3 pt-5">
				<Button
					variant="ghost"
					class="cursor-pointer px-5 hover:bg-zinc-900 hover:text-white"
					onclick={handleCancel}
				>
					Cancel
				</Button>
				<Button
					class="cursor-pointer px-6"
					disabled={(mode === 'add' && !selectedServiceId) ||
						(componentMode === 'custom' && selectedComponentIds.length === 0)}
					onclick={saveService}
				>
					{mode === 'add' ? 'Add Service' : 'Save Changes'}
				</Button>
			</div>
		{:else}
			<div class="py-6 text-center text-zinc-500">Loading service details...</div>
		{/if}
	</div>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { StatuspageApi } from '$lib/api/statuspage/statuspage';
	import { ViewsApi } from '$lib/api/views/views';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import Check from '@lucide/svelte/icons/check';
	import Search from '@lucide/svelte/icons/search';
	import Select from 'svelte-select';
	import { toast } from 'svelte-sonner';

	interface ServiceComponent {
		id: number;
		name: string;
	}

	interface ServiceComponentGroup {
		id: number;
		name: string;
		components: ServiceComponent[];
	}

	let {
		mode,
		publicId,
		serviceSlug = ''
	}: { mode: 'add' | 'edit'; publicId: string; serviceSlug?: string } = $props();

	const viewsApi = new ViewsApi();
	const statuspageApi = new StatuspageApi();

	let selectedService = $state<{ id: number; name: string; slug: string; url: string } | null>(
		null
	);
	let selectedServiceIdStr = $state('');
	let selectedServiceId = $derived(selectedServiceIdStr ? Number(selectedServiceIdStr) : null);
	let componentMode = $state<'all' | 'custom'>('all');
	let errorMessage = $state<string | null>(null);
	let saving = $state(false);

	let configuringService = $state<{
		service_name: string;
		grouped_components: ServiceComponentGroup[];
		ungrouped_components: ServiceComponent[];
	} | null>(null);

	let selectedComponentIds = $state<number[]>([]);
	let selectedComponentGroupIds = $state<number[]>([]);

	let activeService = $state<{ id: number; name: string } | null>(null);

	let componentSearchQuery = $state('');

	let filteredGroupedComponents = $derived.by(() => {
		if (!configuringService) return [];
		const query = componentSearchQuery.toLowerCase().trim();
		if (!query) return configuringService.grouped_components;

		return configuringService.grouped_components
			.map((group) => {
				const matchesGroup = group.name.toLowerCase().includes(query);
				const matchingComponents = group.components.filter((c) =>
					c.name.toLowerCase().includes(query)
				);

				if (matchesGroup || matchingComponents.length > 0) {
					return {
						...group,
						components: matchesGroup ? group.components : matchingComponents
					};
				}
				return null;
			})
			.filter((g): g is NonNullable<typeof g> => g !== null);
	});

	let filteredUngroupedComponents = $derived.by(() => {
		if (!configuringService) return [];
		const query = componentSearchQuery.toLowerCase().trim();
		if (!query) return configuringService.ungrouped_components;

		return configuringService.ungrouped_components.filter((c) =>
			c.name.toLowerCase().includes(query)
		);
	});

	async function fetchServiceComponents(slug: string) {
		const [res, err] = await statuspageApi.getComponents(slug);
		if (err) {
			errorMessage = err.message || 'Failed to fetch components for service';
			toast.error(errorMessage);
			return null;
		}
		return res;
	}

	async function fetchViewServiceConfig(serviceId: number) {
		const [res, err] = await viewsApi.getViewService(publicId, serviceId);
		if (err) {
			errorMessage = err.message || 'Failed to fetch existing configuration';
			toast.error(errorMessage);
			return;
		}
		if (res) {
			if (!res.include_all_components) {
				selectedComponentIds = res.component_ids ?? [];
				selectedComponentGroupIds = res.component_group_ids ?? [];
			}
		}
	}

	$effect(() => {
		if (mode === 'edit' && serviceSlug) {
			errorMessage = null;
			void fetchServiceComponents(serviceSlug).then((res) => {
				if (res) {
					configuringService = {
						service_name: res.service_name,
						grouped_components: res.grouped_components,
						ungrouped_components: res.ungrouped_components
					};
					activeService = {
						id: res.service_id,
						name: res.service_name
					};
					void fetchViewServiceConfig(res.service_id);
				}
			});
		}
	});

	$effect(() => {
		if (selectedService) {
			selectedServiceIdStr = String(selectedService.id);
			errorMessage = null;
			void fetchServiceComponents(selectedService.slug).then((res) => {
				if (res) {
					configuringService = {
						service_name: res.service_name,
						grouped_components: res.grouped_components,
						ungrouped_components: res.ungrouped_components
					};
				}
			});
		} else {
			selectedServiceIdStr = '';
			configuringService = null;
		}
	});

	async function loadUnconfiguredServices(filterText: string) {
		const [res, err] = await viewsApi.getUnconfiguredServices(publicId, filterText || undefined);
		if (err) {
			toast.error(err.message || 'Failed to fetch unconfigured services');
			return [];
		}
		return res ?? [];
	}

	function toggleComponent(componentId: number, groupId?: number) {
		if (!configuringService) return;
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

	function toggleGroup(group: ServiceComponentGroup) {
		if (!configuringService) return;
		const allChecked = group.components.every((c: ServiceComponent) =>
			selectedComponentIds.includes(c.id)
		);
		if (allChecked) {
			selectedComponentGroupIds = selectedComponentGroupIds.filter((id) => id !== group.id);
			selectedComponentIds = selectedComponentIds.filter(
				(id) => !group.components.some((c: ServiceComponent) => c.id === id)
			);
		} else {
			selectedComponentGroupIds = [...selectedComponentGroupIds, group.id];
			const newIds = group.components
				.map((c: ServiceComponent) => c.id)
				.filter((id: number) => !selectedComponentIds.includes(id));
			selectedComponentIds = [...selectedComponentIds, ...newIds];
		}
	}

	async function saveService() {
		if (saving) return;

		errorMessage = null;
		saving = true;

		const componentIds = componentMode === 'all' ? [] : selectedComponentIds;
		const componentGroupIds = componentMode === 'all' ? [] : selectedComponentGroupIds;
		const includeAllComponents = componentMode === 'all';

		if (mode === 'add') {
			if (!selectedServiceId) {
				errorMessage = 'Please select a service';
				saving = false;
				return;
			}

			const [_, err] = await viewsApi.addViewService(publicId, {
				service_id: selectedServiceId,
				include_all_components: includeAllComponents,
				component_ids: componentIds,
				component_group_ids: componentGroupIds
			});

			saving = false;

			if (err) {
				errorMessage = err.message || 'Failed to add service to view';
				toast.error(errorMessage);
				return;
			}

			toast.success('Service added successfully');
			void goto(resolve(`/views/${publicId}`));
		} else {
			if (!activeService) {
				errorMessage = 'Service details not loaded';
				saving = false;
				return;
			}

			const [_, err] = await viewsApi.editViewService(publicId, activeService.id, {
				include_all_components: includeAllComponents,
				component_ids: componentIds,
				component_group_ids: componentGroupIds
			});

			saving = false;

			if (err) {
				errorMessage = err.message || 'Failed to save service configuration';
				toast.error(errorMessage);
				return;
			}

			toast.success('Service configuration saved successfully');
			void goto(resolve(`/views/${publicId}`));
		}
	}

	function handleCancel() {
		void goto(resolve(`/views/${publicId}`));
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
					<Label for="service-search" class="text-lg font-bold text-zinc-200">Service Name</Label>
					{#if mode === 'add'}
						<div id="service-search-container" class="svelte-select-container relative">
							<Select
								loadOptions={loadUnconfiguredServices}
								itemId="id"
								label="name"
								bind:value={selectedService}
								placeholder="Search services to add..."
								class="w-full text-sm text-white outline-none"
							>
								<div
									slot="prepend"
									class="pointer-events-none absolute top-1/2 left-3 z-10 flex -translate-y-1/2 items-center text-zinc-500"
								>
									<Search class="size-4" />
								</div>
								<div slot="item" let:item class="flex w-full items-center justify-between py-1">
									<span>{item.name}</span>
									{#if selectedService && selectedService.id === item.id}
										<span class="text-emerald-400">
											<Check class="size-4" />
										</span>
									{/if}
								</div>
							</Select>
						</div>
					{:else if activeService}
						<div class="rounded-lg bg-zinc-900/50 px-4 py-3 text-sm font-medium text-zinc-200">
							{activeService.name}
						</div>
					{/if}
				</div>

				<!-- Component Checklist with Hierarchical groups -->
				{#if configuringService}
					<div class="grid gap-4 pt-5">
						<Label class="text-lg font-bold text-zinc-200">Monitored Components</Label>

						<!-- Custom Radio Buttons and Component checklist in borderless container -->
						<div class="flex flex-col gap-4 rounded-lg bg-zinc-900/20 p-4">
							<label class="group flex cursor-pointer items-start gap-3">
								<input
									type="radio"
									name="component-mode"
									value="all"
									checked={componentMode === 'all'}
									onchange={() => {
										componentMode = 'all';
									}}
									class="sr-only"
								/>
								<div
									class="mt-1 flex size-5 shrink-0 items-center justify-center rounded-full border border-zinc-700 transition-all group-hover:border-zinc-500 {componentMode ===
									'all'
										? 'border-emerald-500 bg-emerald-500/10'
										: ''}"
								>
									{#if componentMode === 'all'}
										<div class="size-2.5 rounded-full bg-emerald-500"></div>
									{/if}
								</div>
								<div class="flex flex-col">
									<span
										class="text-lg font-semibold text-zinc-200 transition-colors group-hover:text-white"
									>
										Monitor all components
									</span>
									<span class="mt-0.5 text-sm text-zinc-500">
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
									class="mt-1 flex size-5 shrink-0 items-center justify-center rounded-full border border-zinc-700 transition-all group-hover:border-zinc-500 {componentMode ===
									'custom'
										? 'border-emerald-500 bg-emerald-500/10'
										: ''}"
								>
									{#if componentMode === 'custom'}
										<div class="size-2.5 rounded-full bg-emerald-500"></div>
									{/if}
								</div>
								<div class="flex flex-col">
									<span
										class="text-lg font-semibold text-zinc-200 transition-colors group-hover:text-white"
									>
										Customize component selection
									</span>
									<span class="mt-0.5 text-sm text-zinc-500">
										Select specific component groups or individual components to monitor
									</span>
								</div>
							</label>

							{#if componentMode === 'custom'}
								<div class="ml-7 flex flex-col gap-4">
									<!-- Search Input -->
									<div class="relative w-full">
										<div
											class="pointer-events-none absolute inset-y-0 left-3 flex items-center text-zinc-500"
										>
											<Search class="size-4" />
										</div>
										<input
											type="text"
											bind:value={componentSearchQuery}
											placeholder="Filter components by name..."
											class="w-full rounded-lg border border-zinc-800 bg-zinc-950/40 py-2.5 pr-4 pl-10 text-sm text-white placeholder-zinc-500 transition-colors outline-none focus:border-zinc-700"
										/>
									</div>

									<div class="grid max-h-[36rem] grid-cols-1 gap-3 overflow-y-auto pr-1">
										{#each filteredGroupedComponents as group (group.id)}
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
														class="flex size-5.5 shrink-0 items-center justify-center rounded transition-all {groupChecked
															? 'bg-emerald-500 text-zinc-950'
															: groupSomeChecked
																? 'bg-emerald-500/20 text-emerald-400'
																: 'bg-zinc-900'}"
													>
														{#if groupChecked}
															<Check class="size-3.5" />
														{:else if groupSomeChecked}
															<div class="size-2 rounded-sm bg-emerald-400"></div>
														{/if}
													</div>
													<span class="text-lg font-bold text-zinc-200">{group.name} Group</span>
												</button>

												<!-- Child Components -->
												<div class="mt-2 ml-7 grid gap-2.5 pl-3">
													{#each group.components as component (component.id)}
														{@const componentChecked = selectedComponentIds.includes(component.id)}
														<button
															type="button"
															onclick={() => toggleComponent(component.id, group.id)}
															class="flex cursor-pointer items-center gap-3 text-left transition-all hover:text-zinc-200 {componentChecked
																? 'text-zinc-300'
																: 'text-zinc-500'}"
														>
															<div
																class="flex size-5 shrink-0 items-center justify-center rounded transition-all {componentChecked
																	? 'bg-emerald-500/80 text-zinc-950'
																	: 'bg-zinc-900'}"
															>
																{#if componentChecked}
																	<Check class="size-3" />
																{/if}
															</div>
															<span class="text-base font-medium">{component.name}</span>
														</button>
													{/each}
												</div>
											</div>
										{/each}

										{#if filteredUngroupedComponents.length > 0}
											<div class="py-2">
												<div class="mt-2 ml-7 grid gap-2.5 pl-3">
													{#each filteredUngroupedComponents as component (component.id)}
														{@const componentChecked = selectedComponentIds.includes(component.id)}
														<button
															type="button"
															onclick={() => toggleComponent(component.id)}
															class="flex cursor-pointer items-center gap-3 text-left transition-all hover:text-zinc-200 {componentChecked
																? 'text-zinc-300'
																: 'text-zinc-500'}"
														>
															<div
																class="flex size-5 shrink-0 items-center justify-center rounded transition-all {componentChecked
																	? 'bg-emerald-500/80 text-zinc-950'
																	: 'bg-zinc-900'}"
															>
																{#if componentChecked}
																	<Check class="size-3" />
																{/if}
															</div>
															<span class="text-base font-medium">{component.name}</span>
														</button>
													{/each}
												</div>
											</div>
										{/if}

										{#if filteredGroupedComponents.length === 0 && filteredUngroupedComponents.length === 0}
											<div class="py-4 text-center text-sm text-zinc-500">
												No components match your search filter.
											</div>
										{/if}
									</div>
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
					disabled={saving}
					onclick={handleCancel}
				>
					Cancel
				</Button>
				<Button
					class="cursor-pointer px-6"
					disabled={saving ||
						(mode === 'add' && !selectedServiceId) ||
						(componentMode === 'custom' && selectedComponentIds.length === 0)}
					onclick={saveService}
				>
					{saving ? 'Saving...' : mode === 'add' ? 'Add Service' : 'Save Changes'}
				</Button>
			</div>
		{:else}
			<div class="py-6 text-center text-zinc-500">Loading service details...</div>
		{/if}
	</div>
</div>

<style>
	:global(.svelte-select-container) {
		--background: rgba(24, 24, 27, 0.5) !important;
		--border: 1px solid #27272a !important;
		--border-hover: 1px solid #3f3f46 !important;
		--border-focused: 1px solid #3f3f46 !important;
		--list-background: #09090b !important;
		--item-color: #a1a1aa !important;
		--item-hover-bg: rgba(24, 24, 27, 0.5) !important;
		--item-hover-color: #ffffff !important;
		--item-active-background: #27272a !important;
		--item-is-active-bg: #27272a !important;
		--input-color: #ffffff !important;
		--placeholder-color: #71717a !important;
		--chevron-color: #a1a1aa !important;
		--clear-icon-color: #a1a1aa !important;
		--list-empty-color: #71717a !important;
		--height: 48px !important;
		--input-padding: 0 0 0 28px !important;
		--value-container-padding: 0 0 0 28px !important;
		border-radius: 8px !important;
	}
</style>

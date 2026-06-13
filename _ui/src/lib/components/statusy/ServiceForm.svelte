<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { StatuspageApi } from '$lib/api/statuspage/statuspage';
	import { ViewsApi } from '$lib/api/views/views';
	import * as Accordion from '$lib/components/ui/accordion';
	import { ControlledCheckbox } from '$lib/components/ui/checkbox';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import * as RadioGroup from '$lib/components/ui/radio-group';
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
	const componentModeOptions = [
		{
			value: 'all' as const,
			id: 'mode-all',
			label: 'Monitor all components',
			description: 'Automatically monitor all current and future components for this service'
		},
		{
			value: 'custom' as const,
			id: 'mode-custom',
			label: 'Customize component selection',
			description: 'Select specific component groups or individual components to monitor'
		}
	];
	let errorMessage = $state<string | null>(null);
	let saving = $state(false);

	let monitorIncidents = $state(true);
	let monitorMaintenances = $state(true);

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
			componentMode = res.include_all_components ? 'all' : 'custom';
			monitorIncidents = res.monitor_incidents;
			monitorMaintenances = res.monitor_scheduled_maintenances;
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

		if (groupId) {
			const group = configuringService.grouped_components.find((g) => g.id === groupId);
			if (group) {
				const isGroupSelected = selectedComponentGroupIds.includes(groupId);
				if (isGroupSelected) {
					// Group was fully selected. Unchecking this component means the group is no longer fully selected.
					// We remove the group ID, and add all OTHER components of this group to selectedComponentIds.
					selectedComponentGroupIds = selectedComponentGroupIds.filter((id) => id !== groupId);
					const otherComponentIds = group.components
						.map((c) => c.id)
						.filter((id) => id !== componentId);
					selectedComponentIds = [...selectedComponentIds, ...otherComponentIds];
				} else {
					// Group was not fully selected.
					const isCompSelected = selectedComponentIds.includes(componentId);
					if (isCompSelected) {
						// Uncheck it
						selectedComponentIds = selectedComponentIds.filter((id) => id !== componentId);
					} else {
						// Check it
						const newSelected = [...selectedComponentIds, componentId];
						// Check if all components of the group are now selected
						const allChecked = group.components.every((c) => newSelected.includes(c.id));
						if (allChecked) {
							// Upgrade to group selection: add group, remove individual components of this group
							selectedComponentGroupIds = [...selectedComponentGroupIds, groupId];
							selectedComponentIds = selectedComponentIds.filter(
								(id) => !group.components.some((c) => c.id === id)
							);
						} else {
							selectedComponentIds = newSelected;
						}
					}
				}
				return;
			}
		}

		// Ungrouped component logic
		if (selectedComponentIds.includes(componentId)) {
			selectedComponentIds = selectedComponentIds.filter((id) => id !== componentId);
		} else {
			selectedComponentIds = [...selectedComponentIds, componentId];
		}
	}

	function toggleGroup(group: ServiceComponentGroup) {
		if (!configuringService) return;
		const isGroupSelected = selectedComponentGroupIds.includes(group.id);
		if (isGroupSelected) {
			// Uncheck the group: remove group ID and make sure none of its components are selected
			selectedComponentGroupIds = selectedComponentGroupIds.filter((id) => id !== group.id);
			selectedComponentIds = selectedComponentIds.filter(
				(id) => !group.components.some((c) => c.id === id)
			);
		} else {
			// Check the group: add group ID and remove any individual components of this group from selectedComponentIds
			selectedComponentGroupIds = [...selectedComponentGroupIds, group.id];
			selectedComponentIds = selectedComponentIds.filter(
				(id) => !group.components.some((c) => c.id === id)
			);
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
				monitor_incidents: monitorIncidents,
				monitor_scheduled_maintenances: monitorMaintenances,
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
			await goto(resolve(`/views/${publicId}`));
			await invalidateAll();
		} else {
			if (!activeService) {
				errorMessage = 'Service details not loaded';
				saving = false;
				return;
			}

			const [_, err] = await viewsApi.editViewService(publicId, activeService.id, {
				include_all_components: includeAllComponents,
				monitor_incidents: monitorIncidents,
				monitor_scheduled_maintenances: monitorMaintenances,
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
			await goto(resolve(`/views/${publicId}`));
			await invalidateAll();
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
				? 'Choose a service and select which components to subscribe to for this view.'
				: 'Customize the component subscriptions for this service.'}
		</p>
	</div>

	<!-- Form Card (Borderless & Background Merged) -->
	<div class="max-w-xl">
		{#if mode === 'add'}
			<div class="mb-6 rounded-lg border border-zinc-800 bg-zinc-900/30 p-4 text-sm text-zinc-300">
				<p class="font-medium text-white">Can't find a status page?</p>
				<p class="mt-1 text-xs text-zinc-400">
					If the service you want to track is not listed, you can
					<a
						href="https://github.com/yash492/statusy/issues/new"
						target="_blank"
						rel="noopener noreferrer"
						class="text-zinc-200 underline hover:text-white"
						>request it by creating an issue here</a
					>.
				</p>
			</div>
		{/if}

		{#if errorMessage}
			<div class="mb-4 rounded-lg border border-red-500/20 bg-red-950/20 p-3 text-sm text-red-400">
				{errorMessage}
			</div>
		{/if}

		{#if mode === 'add' || activeService}
			<div class="grid gap-6">
				<!-- Service Name Input/Dropdown Selector -->
				<div class="grid gap-2">
					<Label for="service-search" class="text-base font-bold text-zinc-200">Service Name</Label>
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
										<span class="text-zinc-200">
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

				<!-- Alert Types to Monitor -->
				{#if configuringService}
					<div class="grid gap-4 pt-5">
						<Label class="text-base font-bold text-zinc-200">Alert Types to Monitor</Label>
						<div class="flex flex-col gap-3 rounded-lg bg-zinc-900/20 p-4">
							<label class="flex cursor-pointer items-center gap-3 select-none">
								<ControlledCheckbox
									id="monitor-incidents"
									checked={monitorIncidents}
									onchange={() => monitorIncidents = !monitorIncidents}
									class="size-5 shrink-0 cursor-pointer"
								/>
								<span class="text-sm font-semibold text-zinc-200">Monitor Incidents</span>
							</label>
							<label class="flex cursor-pointer items-center gap-3 select-none">
								<ControlledCheckbox
									id="monitor-maintenances"
									checked={monitorMaintenances}
									onchange={() => monitorMaintenances = !monitorMaintenances}
									class="size-5 shrink-0 cursor-pointer"
								/>
								<span class="text-sm font-semibold text-zinc-200">Monitor Scheduled Maintenances</span>
							</label>
						</div>
					</div>
				{/if}

				<!-- Component Checklist with Hierarchical groups -->
				{#if configuringService}
					<div class="grid gap-4 pt-5">
						<Label class="text-base font-bold text-zinc-200">Monitored Components</Label>

						<!-- Radio Buttons and Component checklist in borderless container -->
						<div class="flex flex-col gap-4 rounded-lg bg-zinc-900/20 p-4">
							<RadioGroup.Root bind:value={componentMode} class="flex flex-col gap-5">
								{#each componentModeOptions as opt}
									<div class="flex flex-col gap-1">
										<div class="flex items-center gap-3">
											<RadioGroup.Item
												value={opt.value}
												id={opt.id}
												class="size-5 shrink-0 cursor-pointer"
											/>
											<Label
												for={opt.id}
												class="cursor-pointer text-base font-semibold text-zinc-200 transition-colors select-none hover:text-white"
											>
												{opt.label}
											</Label>
										</div>
										<span class="pl-8 text-xs font-normal text-zinc-500 select-none">
											{opt.description}
										</span>
									</div>
								{/each}
							</RadioGroup.Root>

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

									<div class="grid max-h-[36rem] grid-cols-1 gap-1 overflow-y-auto pr-3">
										<Accordion.Root type="multiple" class="w-full">
											{#each filteredGroupedComponents as group (group.id)}
												{@const groupChecked = selectedComponentGroupIds.includes(group.id)}
												{@const groupSomeChecked =
													!groupChecked &&
													group.components.some((c) => selectedComponentIds.includes(c.id))}

												<Accordion.Item value={String(group.id)} class="border-none">
													<!-- Group Header Row -->
													<div
														class="flex w-full items-center justify-between gap-3 border-b border-zinc-800/40 py-2"
													>
														<div class="flex flex-1 items-center gap-3">
															<ControlledCheckbox
																id={`group-${group.id}`}
																checked={groupChecked}
																indeterminate={groupSomeChecked}
																onchange={() => toggleGroup(group)}
																class="size-5 shrink-0 cursor-pointer"
															/>
															<Label
																for={`group-${group.id}`}
																class="flex-1 cursor-pointer py-1 text-sm font-medium transition-colors select-none {groupChecked
																	? 'text-zinc-100'
																	: 'text-zinc-300'}"
															>
																{group.name}
																<span class="ml-1 text-xs font-normal text-zinc-500">
																	({group.components.length})
																</span>
															</Label>
														</div>
														<Accordion.Trigger
															class="hover:bg-zinc-850 flex h-8 w-8 shrink-0 items-center justify-center rounded-md p-1.5 text-zinc-400 transition-colors hover:text-white"
														/>
													</div>

													<Accordion.Content class="pt-0 pb-0">
														<!-- Child Components -->
														<div class="flex w-full flex-col">
															{#each group.components as component (component.id)}
																{@const componentChecked =
																	groupChecked || selectedComponentIds.includes(component.id)}
																<div
																	class="flex w-full items-center gap-3 border-b border-zinc-800/40 py-2 pl-8 last:border-none"
																>
																	<ControlledCheckbox
																		id={`comp-${component.id}`}
																		checked={componentChecked}
																		onchange={() => toggleComponent(component.id, group.id)}
																		class="size-5 shrink-0 cursor-pointer"
																	/>
																	<Label
																		for={`comp-${component.id}`}
																		class="flex-1 cursor-pointer py-1 text-sm font-medium transition-colors select-none {componentChecked
																			? 'text-zinc-100'
																			: 'text-zinc-300'}"
																	>
																		{component.name}
																	</Label>
																</div>
															{/each}
														</div>
													</Accordion.Content>
												</Accordion.Item>
											{/each}
										</Accordion.Root>

										{#if filteredUngroupedComponents.length > 0}
											<div class="py-1">
												<div class="flex flex-col gap-1">
													{#each filteredUngroupedComponents as component (component.id)}
														{@const componentChecked = selectedComponentIds.includes(component.id)}
														<div
															class="flex w-full items-center gap-3 border-b border-zinc-800/40 py-2 last:border-none"
														>
															<ControlledCheckbox
																id={`comp-${component.id}`}
																checked={componentChecked}
																onchange={() => toggleComponent(component.id)}
																class="size-5 shrink-0 cursor-pointer"
															/>
															<Label
																for={`comp-${component.id}`}
																class="flex-1 cursor-pointer py-1 text-sm font-medium transition-colors select-none {componentChecked
																	? 'text-zinc-100'
																	: 'text-zinc-300'}"
															>
																{component.name}
															</Label>
														</div>
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
						(!monitorIncidents && !monitorMaintenances) ||
						(componentMode === 'custom' &&
							selectedComponentIds.length === 0 &&
							selectedComponentGroupIds.length === 0)}
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

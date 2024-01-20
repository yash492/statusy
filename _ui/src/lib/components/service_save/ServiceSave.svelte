<script lang="ts">
	import Select from 'svelte-select';
	import Button from '../button/Button.svelte';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { MagnifyingGlass } from '@steeze-ui/heroicons';
	import type { ComponentsForService } from '$lib/types/components';
	import type { SaveSubscription } from '$lib/types/subscriptions';
	import { afterNavigate } from '$app/navigation';
	import { SUBSCRIPTION_BY_ID_QUERY_KEY } from '$lib/types/query_keys';
	import { useQueryClient } from '@tanstack/svelte-query';

	const ALL_COMPONENTS_RADIO_BUTTON = 'all-components';
	const CUSTOM_COMPONENTS_RADIO_BUTTON = 'custom-components';

	const queryClient = useQueryClient();
	afterNavigate(() => {
		queryClient.invalidateQueries({ queryKey: [SUBSCRIPTION_BY_ID_QUERY_KEY] });
	});

	export let onSaveService: (subscription: SaveSubscription) => void;
	export let fetchServices:
		| ((filterText: string) => Promise<
				{
					value: number;
					label: string;
				}[]
		  >)
		| null = null;
	export let selectedService: number | undefined;
	export let isAllComponents: boolean;

	export let components: ComponentsForService[];
	export let customComponentCheckbox: number[];
	export let editMode: boolean;
	export let serviceName: string = '';

	$: isAllComponents = isAllComponents;
	$: components = components;
	$: customComponentCheckbox = customComponentCheckbox;

	let searchComponentsValue = '';
	let searchingComponents = components;
	let selectedComponentChoice = isAllComponents
		? ALL_COMPONENTS_RADIO_BUTTON
		: CUSTOM_COMPONENTS_RADIO_BUTTON;

	$: {
		searchingComponents = components.filter((component) => {
			return component.name.toLowerCase().includes(searchComponentsValue.toLowerCase());
		});
	}

	$: {
		isAllComponents = selectedComponentChoice === ALL_COMPONENTS_RADIO_BUTTON ? true : false;
		if (isAllComponents) {
			customComponentCheckbox = [];
		}
	}
</script>

<form
	on:submit|preventDefault={() =>
		onSaveService({
			custom_components: customComponentCheckbox,
			is_all_components: Boolean(isAllComponents),
			service_id: selectedService || 0
		})}
>
	<h1 class="font-bold text-xl pt-5">
		{editMode ? `Edit Service - ${serviceName} ` : 'Select a Service'}
	</h1>
	<div class="md:w-[40rem] mt-4">
		<div>
			<label for="service-select" class="py-1 block font-medium"
				>{editMode ? 'Your Service' : 'Select a Service'}</label
			>
			<div class="block">
				{#if !editMode}
					<Select
						id="service-select"
						--background="rgb(38,38,38)"
						--list-background="rgb(38,38,38)"
						--border="rgb(115,115,115)"
						--border-hover="rgb(115,115,115)"
						--placeholder-color="rgb(120,120,120)"
						--item-hover-bg="rgb(70,70,70)"
						placeholder="Choose a service"
						loadOptions={fetchServices}
						bind:justValue={selectedService}
						required
					></Select>
				{:else}
					<div>
						<p class="w-full border py-2 px-4 rounded-md bg-neutral-800 border-neutral-700">
							{serviceName}
						</p>
					</div>
				{/if}
			</div>
		</div>
		{#if selectedService}
			<div>
				<label for="components-select" class="mt-3 block py-1 font-medium">Select Components</label>
				<div>
					<div>
						<input
							bind:group={selectedComponentChoice}
							id="all-components-radio-button"
							name="all-components-radio-button"
							type="radio"
							value={ALL_COMPONENTS_RADIO_BUTTON}
							class="mr-1"
						/>
						<label for="all-components-radio-button">All Components</label>
					</div>
					<div class="pb-2">
						<input
							bind:group={selectedComponentChoice}
							id="custom-components-radio-button"
							name="custom-components-radio-button"
							type="radio"
							value={CUSTOM_COMPONENTS_RADIO_BUTTON}
							class="mr-1"
						/>
						<label for="custom-components-radio-button">Choose Components</label>
					</div>
					{#if selectedComponentChoice === CUSTOM_COMPONENTS_RADIO_BUTTON}
						<div class="border bg-neutral-800 rounded-md border-neutral-500 py-3 px-3 text-sm">
							<div class="flex items-center rounded-md bg-neutral-700 mb-2">
								<Icon src={MagnifyingGlass} size="23px" class="ml-2"></Icon>
								<input
									class="py-2 px-3 w-full bg-neutral-700 focus:outline-none"
									placeholder="Search Components..."
									bind:value={searchComponentsValue}
								/>
							</div>
							<div class="grid md:grid-cols-2">
								{#each searchingComponents as component}
									<ul class="py-2">
										<input
											type="checkbox"
											class="border-none mr-1"
											bind:group={customComponentCheckbox}
											value={component.id}
										/>
										{component.name}
									</ul>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			</div>
			<div class="mt-5">
				<Button type="submit">{editMode ? 'Update Service' : 'Add Service'}</Button>
			</div>
		{/if}
	</div>
</form>

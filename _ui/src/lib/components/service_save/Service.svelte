<script lang="ts">
	import Select from 'svelte-select';
	import Button from '../button/Button.svelte';
	import type { ComponentsForService } from '$lib/types/components';
	import { createMutation, createQuery } from '@tanstack/svelte-query';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { MagnifyingGlass } from '@steeze-ui/heroicons';
	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import { ComponentsAPI } from '$lib/apis/components';
	import { SERVICE_LIST_FOR_SUBSCRIPTION_QUERY_KEY } from '$lib/types/query_keys';
	import { goto } from '$app/navigation';
	import type { AddSubscription } from '$lib/types/subscriptions';

	const _subscriptionAPI = new SubscriptionAPI();
	const _componentsAPI = new ComponentsAPI();

	const ALL_COMPONENTS_RADIO_BUTTON = 'all-components';
	const CUSTOM_COMPONENTS_RADIO_BUTTON = 'custom-components';

	let selectedComponentChoice = ALL_COMPONENTS_RADIO_BUTTON;
	let searchComponentsValue = '';

	let components: ComponentsForService[] = [];
	let selectedService: number | undefined;
	let customComponentCheckbox: number[] = [];

	$: query = createQuery({
		queryKey: [SERVICE_LIST_FOR_SUBSCRIPTION_QUERY_KEY],
		queryFn: () => _componentsAPI.ComponentsForService(selectedService || 0)
	});

	$: {
		components = $query.data?.data.data || [];
		components = components.filter((component) => {
			return component.name.toLowerCase().includes(searchComponentsValue.toLowerCase());
		});
	}

	async function fetchServices(filterText: string) {
		const services = await _subscriptionAPI.ServicesList(filterText);

		const list = services.data.data.map((item) => {
			return {
				value: item.id,
				label: item.name
			};
		});

		return list;
	}

	const addServiceMutation = createMutation<{}, {}, AddSubscription>({
		mutationFn: (subscription) => _subscriptionAPI.Add(subscription),
		onSuccess: () => {
			goto('/dashboard');
		},
		onError: () => {}
	});

	async function onSaveService() {
		const subscription = {
			custom_components: customComponentCheckbox,
			is_all_components: selectedComponentChoice === ALL_COMPONENTS_RADIO_BUTTON,
			service_id: selectedService || 0
		};

		$addServiceMutation.mutate(subscription);
	}
</script>

<form on:submit|preventDefault={onSaveService}>
	<h1 class="font-bold text-xl pt-5">Add Service</h1>
	<div class="md:w-[40rem] mt-4">
		<div>
			<label for="service-select" class="py-1 block font-medium">Select a Service</label>
			<div class="block">
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
							value="all-components"
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
							value="custom-components"
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
								{#if $query.isLoading}
									<p>Loading Components ...</p>
								{:else if $query.isError}
									<p>Error while loading components {$query.error}</p>
								{:else}
									{#each components as component}
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
								{/if}
							</div>
						</div>
					{/if}
				</div>
			</div>
			<div class="mt-5">
				<Button type="submit">Add Service</Button>
			</div>
		{/if}
	</div>
</form>

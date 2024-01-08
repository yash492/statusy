<script lang="ts">
	import { createMutation, createQuery } from '@tanstack/svelte-query';
	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import { SUBSCRIPTION_BY_ID_QUERY_KEY } from '$lib/types/query_keys';
	import ServiceSave from '$lib/components/service_save/ServiceSave.svelte';
	import type { ComponentsForService } from '$lib/types/components';
	import type { SaveSubscription } from '$lib/types/subscriptions';
	import type { PageData } from './$types';
	import { Toast } from '$lib/toast/toast';
	import { goto } from '$app/navigation';

	const _subscriptionAPI = new SubscriptionAPI();

	export let data: PageData;
	const _toast = new Toast();

	$: query = createQuery({
		queryKey: [SUBSCRIPTION_BY_ID_QUERY_KEY],
		queryFn: () => _subscriptionAPI.GetByID(data.subscriptionID)
	});

	$: subscriptionData = $query.data?.data.data;

	let components: ComponentsForService[] = [];
	let customComponentCheckbox: number[] = [];
	let isAllComponentsFromAPI = true;
	let selectedServiceID = 0;
	let serviceName = '';

	$: {
		selectedServiceID = subscriptionData?.service_id || 0;
		serviceName = subscriptionData?.service_name || '';
		isAllComponentsFromAPI = !!subscriptionData?.is_all_components;
	}

	$: isAllComponents = isAllComponentsFromAPI;
	$: console.log(isAllComponents, isAllComponentsFromAPI, 'is all components from api');

	$: {
		if (components.length < 1) {
			for (const component of subscriptionData?.components || []) {
				components.push({
					id: component.id,
					name: component.name
				});

				if (component.is_configured) {
					customComponentCheckbox.push(component.id);
				}
			}
			components = components;
			customComponentCheckbox = customComponentCheckbox;
		}
	}

	const mutation = createMutation<{}, {}, SaveSubscription>({
		mutationFn: (subscription) => _subscriptionAPI.Update(subscription, data.subscriptionID),
		onSuccess: () => {
			_toast.success(`${serviceName} service successfully updated!`);
			goto('/dashboard');
		},
		onError: () => {}
	});

	function onSaveService() {
		const subscription: SaveSubscription = {
			custom_components: customComponentCheckbox,
			is_all_components: isAllComponents,
			service_id: selectedServiceID || 0
		};

		if (!subscription.is_all_components && subscription.custom_components.length === 0) {
			_toast.error('Please choose components for this option');
			return;
		}

		$mutation.mutate(subscription);
	}
</script>

<div class="mx-auto w-fit pb-11">
	{#if $query.isLoading}
		<p>Loading...</p>
	{:else if $query.isError}
		<p>Error {$query.error}</p>
	{:else}
		<ServiceSave
			editMode={true}
			{onSaveService}
			bind:selectedService={selectedServiceID}
			{components}
			bind:customComponentCheckbox
			bind:isAllComponents
			{serviceName}
		/>
	{/if}
</div>

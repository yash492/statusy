<script lang="ts">
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import { SUBSCRIPTION_BY_ID_QUERY_KEY } from '$lib/types/query_keys';
	import ServiceSave from '$lib/components/service_save_form/ServiceSaveForm.svelte';
	import type { ComponentsForService } from '$lib/types/components';
	import type { SaveSubscription } from '$lib/types/subscriptions';
	import { Toast } from '$lib/toast/toast';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	const _subscriptionAPI = new SubscriptionAPI();
	const _toast = new Toast();

	$: subscriptionIDFromParams = $page.params.subscription_id;
	$: query = createQuery({
		queryKey: [SUBSCRIPTION_BY_ID_QUERY_KEY],
		queryFn: async () => {
			return await _subscriptionAPI.GetByID(subscriptionIDFromParams);
		}
	});

	$: subscriptionData = $query.data?.data.data;

	let components: ComponentsForService[] = [];
	let customComponentCheckbox: number[] = [];
	$: subscriptionComponents = subscriptionData?.components || [];

	$: if (subscriptionComponents.length > 0) {
		// To prevent appending of old subscription components
		components = [];
		customComponentCheckbox = [];
		for (const component of subscriptionComponents) {
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

	const mutation = createMutation<{}, {}, SaveSubscription>({
		mutationFn: async (subscription) =>
			await _subscriptionAPI.Update(subscription, subscriptionIDFromParams),
		onSuccess: () => {
			_toast.success(`${subscriptionData?.service_name} service successfully updated!`);
			goto('/dashboard');
		},
		onError: () => {}
	});

	function onSaveService(subscription: SaveSubscription) {
		if (!subscription.is_all_components && subscription.custom_components.length === 0) {
			_toast.error('Please choose components for this option');
			return;
		}

		$mutation.mutate(subscription);
	}
</script>

<div class="mx-auto w-full pb-11">
	{#if $query.isLoading}
		<p>Loading...</p>
	{:else if $query.isError}
		<p>Error {$query.error}</p>
	{:else}
		<ServiceSave
			editMode={true}
			{onSaveService}
			selectedService={subscriptionData?.service_id}
			{components}
			{customComponentCheckbox}
			isAllComponents={Boolean(subscriptionData?.is_all_components)}
			serviceName={subscriptionData?.service_name}
		/>
	{/if}
</div>

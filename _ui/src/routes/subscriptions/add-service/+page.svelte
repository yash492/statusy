<script lang="ts">
	import type { ComponentsForService } from '$lib/types/components';
	import { createMutation, createQuery } from '@tanstack/svelte-query';

	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import { ComponentsAPI } from '$lib/apis/components';
	import { SERVICE_LIST_FOR_SUBSCRIPTION_QUERY_KEY } from '$lib/types/query_keys';
	import { goto } from '$app/navigation';
	import type { SaveSubscription } from '$lib/types/subscriptions';
	import ServiceSave from '$lib/components/service_save/ServiceSave.svelte';
	import toast from 'svelte-french-toast';
	import { Toast } from '$lib/toast/toast';

	const _subscriptionAPI = new SubscriptionAPI();
	const _componentsAPI = new ComponentsAPI();
	const _toast = new Toast();

	let components: ComponentsForService[] = [];
	let selectedServiceID: number | undefined;
	let customComponentCheckbox: number[] = [];
	let isAllComponents = true;

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

	$: query = createQuery({
		queryKey: [SERVICE_LIST_FOR_SUBSCRIPTION_QUERY_KEY],
		queryFn: () => _componentsAPI.ComponentsForService(selectedServiceID || 0)
	});

	$: components = $query.data?.data.data || [];

	const mutation = createMutation<{}, {}, SaveSubscription>({
		mutationFn: (subscription) => _subscriptionAPI.Add(subscription),
		onSuccess: () => {
			_toast.success(`Service successfully added!`);
			goto('/dashboard');
		},
		onError: () => {}
	});

	async function onSaveService(subscription: SaveSubscription) {
		if (!subscription.is_all_components && subscription.custom_components.length === 0) {
			_toast.error('Please choose components for custom components option');
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
			{fetchServices}
			editMode={false}
			{onSaveService}
			selectedService={selectedServiceID}
			{components}
			{customComponentCheckbox}
			isAllComponents
		/>
	{/if}
</div>

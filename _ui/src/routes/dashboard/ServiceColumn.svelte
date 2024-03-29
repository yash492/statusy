<script lang="ts">
	import { goto } from '$app/navigation';
	import { SubscriptionAPI } from '$lib/apis/subscriptions';
	import Button from '$lib/components/button/Button.svelte';
	import Modal from '$lib/components/modal/Modal.svelte';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import { Toast } from '$lib/toast/toast';
	import { DASHBOARD_SUBSCRIPTION_LIST_QUERY_KEY } from '$lib/types/query_keys';
	import { AdjustmentsHorizontal, XCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';

	export let serviceName: string;
	export let incidentName: string;
	export let incidentLink: string;
	export let subscriptionUUID: string;

	const queryClient = useQueryClient();

	let showModal = false;

	const _subscriptionAPI = new SubscriptionAPI();
	const toast = new Toast();

	const mutation = createMutation({
		mutationFn: async () => _subscriptionAPI.Delete(subscriptionUUID),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: [DASHBOARD_SUBSCRIPTION_LIST_QUERY_KEY] });
			toast.success('Subscription deleted successfully!');
			showModal = false;
		},
		onError: (err) => {
			toast.error(AxiosResponseErr(err).error_msg);
		}
	});
</script>

<div class="flex items-center justify-between mr-2">
	<div class="w-4/5">
		<p>
			<a class="hover:text-blue-300" href={`/subscriptions/${subscriptionUUID}/incidents`}>
				{serviceName}</a
			>
		</p>
		<p class="text-xs text-neutral-300">
			<a target="_blank" class="hover:text-blue-300" href={incidentLink}>{incidentName}</a>
		</p>
	</div>

	<div class="flex items-center gap-1">
		<button on:click={() => goto(`subscriptions/${subscriptionUUID}/edit`)}>
			<Icon src={AdjustmentsHorizontal} size={'23px'} />
		</button>
		<button on:click={() => (showModal = true)}>
			<Icon src={XCircle} size={'23px'} />
		</button>
	</div>
</div>

<Modal bind:showModal modalTitle="Delete Subscription">
	<div class="flex flex-col">
		<p>Are you sure you want to delete this subscription?</p>
		<div class="flex gap-3 mt-5">
			<Button on:click={() => $mutation.mutate()}>Yes</Button>
			<Button on:click={() => (showModal = false)}>No</Button>
		</div>
	</div>
</Modal>

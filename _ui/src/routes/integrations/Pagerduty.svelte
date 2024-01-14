<script lang="ts">
	import { IntegrationsAPI } from '$lib/apis/integrations';
	import Button from '$lib/components/button/Button.svelte';
	import IntegrationModalWithList from '$lib/components/integration_component_list/IntegrationModalWithList.svelte';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import { Toast } from '$lib/toast/toast';
	import type { SavePagerduty } from '$lib/types/integrations';
	import { INCIDENT_MANAGEMENT_GET_QUERY_KEY } from '$lib/types/query_keys';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	export let routingKey = '';
	export let uuid: string;
	export let isConfigured: boolean;

	const _integrationAPI = new IntegrationsAPI();
	const _toast = new Toast();
	const queryClient = useQueryClient();
	let showModal: boolean;

	const mutation = createMutation({
		mutationFn: (data: SavePagerduty) => {
			return _integrationAPI.SavePagerduty(data);
		},
		onSuccess: () => {
			showModal = false;
			_toast.success('Pagerduty Integration is successfully saved');
			queryClient.invalidateQueries({ queryKey: [INCIDENT_MANAGEMENT_GET_QUERY_KEY] });
		}
	});

	function onSave() {
		$mutation.mutate({
			routing_key: routingKey,
			uuid: isConfigured ? uuid : undefined
		});
	}
</script>

<IntegrationModalWithList
	isIntegrated={isConfigured}
	name="Pagerduty"
	modalTitle={'Pagerduty Integration'}
	bind:showModal
>
	<div class="mt-10 mb-8">
		<label class="font-medium" for="routing-key"
			>Your Routing Key <span class="ml-1"
				>(<a
					class="text-sm hover:text-blue-400 text-neutral-300 hover:underline"
					href="https://support.pagerduty.com/docs/services-and-integrations#generate-a-new-integration-key"
					target="_blank">Where to find it?</a
				>)
			</span>
		</label>
		<input
			class="w-full font-mono bg-neutral-600 outline-none py-2 px-3 mt-2 rounded-md"
			id="routing-key"
			bind:value={routingKey}
		/>
	</div>
	{#if $mutation.isError}
		<p class="mb-5">{AxiosResponseErr($mutation.error)?.error_msg}</p>
	{/if}
	<div class="mb-3 flex gap-4">
		<Button on:click={onSave}>Add</Button>
		<Button>Delete</Button>
	</div>
</IntegrationModalWithList>

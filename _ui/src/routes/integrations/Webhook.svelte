<script lang="ts">
	import { IntegrationsAPI } from '$lib/apis/integrations';
	import Button from '$lib/components/button/Button.svelte';
	import IntegrationModalWithList from '$lib/components/integration_component_list/IntegrationModalWithList.svelte';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import { Toast } from '$lib/toast/toast';
	import type { SaveWebhook } from '$lib/types/integrations';
	import { WEBHOOKS_GET_QUERY_KEY } from '$lib/types/query_keys';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';

	const _integrationAPI = new IntegrationsAPI();
	let showModal = false;
	const queryClient = useQueryClient();
	const _toast = new Toast();

	export let webhookURL = '';
	export let webhookSecret = '';
	export let uuid: string;
	export let isConfigured = false;

	$: mutation = createMutation({
		mutationFn: (data: SaveWebhook) => {
			return _integrationAPI.SaveWebhook(data);
		},
		onSuccess() {
			showModal = false;
			_toast.success('Webhook Integration is sucessfully saved');
			queryClient.invalidateQueries({ queryKey: [WEBHOOKS_GET_QUERY_KEY] });
		},
		onSettled() {
			webhookURL = '';
			webhookSecret = '';
		}
	});

	function onSave() {
		$mutation.mutate({
			webhook_url: webhookURL || '',
			uuid: isConfigured ? uuid : undefined,
			secret: webhookSecret || ''
		});
	}
</script>

<IntegrationModalWithList
	isIntegrated={isConfigured}
	name="Webhook"
	modalTitle={'Webhook Integration'}
	bind:showModal
>
	<div class="mt-10 mb-4">
		<label class="font-medium" for="webhook">Your External Webhook URL </label>
		<div
			class="flex items-center w-full font-mono bg-neutral-600 outline-none py-2 px-3 mt-2 rounded-md"
		>
			<p class="pr-2 font-semibold">POST</p>
			<input class="w-full bg-neutral-600 outline-none" id="webhook" bind:value={webhookURL} />
		</div>
	</div>
	<div class="mb-7">
		<label for="secret" class="font-medium">Secret</label>
		<input
			class="w-full bg-neutral-600 outline-none py-2 mt-2 px-3 rounded-md font-mono"
			id="secret"
			bind:value={webhookSecret}
		/>
	</div>
	{#if $mutation.isError}
		<p class="mb-5">{AxiosResponseErr($mutation.error)?.error_msg}</p>
	{/if}
	<div class="mb-3 flex gap-4">
		<Button on:click={onSave}>Save</Button>
		<Button>Delete</Button>
	</div>
</IntegrationModalWithList>

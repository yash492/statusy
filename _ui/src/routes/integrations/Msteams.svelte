<script lang="ts">
	import { IntegrationsAPI } from '$lib/apis/integrations';
	import Button from '$lib/components/button/Button.svelte';
	import IntegrationModalWithList from '$lib/components/integration_component_list/IntegrationModalWithList.svelte';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import { Toast } from '$lib/toast/toast';
	import type { DeleteChatopsData, SaveChatOps } from '$lib/types/integrations';
	import { CHATOPS_GET_ALL_QUERY_KEY } from '$lib/types/query_keys';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	export let webhookURL = '';
	export let uuid: string | undefined;
	export let isConfigured: boolean;

	const _integrationAPI = new IntegrationsAPI();
	const _toast = new Toast();
	const queryClient = useQueryClient();

	let showModal: boolean;

	$: mutation = createMutation({
		mutationFn: (data: SaveChatOps) => _integrationAPI.SaveChatOps(data),
		onSuccess() {
			queryClient.invalidateQueries({ queryKey: [CHATOPS_GET_ALL_QUERY_KEY] });
			showModal = false;
			_toast.success('M Integration is sucessfully saved');
		},
		onSettled() {
			webhookURL = '';
		}
	});

	function onSave() {
		$mutation.mutate({
			type: 'msteams',
			webhook_url: webhookURL.trim(),
			uuid: isConfigured ? uuid : undefined
		});
	}

	$: deleteMutation = createMutation({
		mutationFn: (data: DeleteChatopsData) => _integrationAPI.DeleteChatOps(data.uuid, data.type),
		onSuccess() {
			showModal = false;
			queryClient.invalidateQueries({ queryKey: [CHATOPS_GET_ALL_QUERY_KEY] });
			_toast.success('MS Teams Integration is sucessfully deleted');
			webhookURL = '';
		}
	});

	function onDelete() {
		$deleteMutation.mutate({
			type: 'msteams',
			uuid: uuid || ''
		});
	}
</script>

<IntegrationModalWithList
	isIntegrated={isConfigured}
	name="MS Teams"
	modalTitle="MS Teams Integration"
	bind:showModal
>
	<div class="mt-10 mb-8">
		<label class="font-medium" for="webhook"
			>Your Webhook URL <span class="ml-1"
				>(<a
					class="text-sm hover:text-blue-400 text-neutral-300 hover:underline"
					href="https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook?tabs=dotnet#create-incoming-webhooks-1"
					target="_blank">Where to find it?</a
				>)
			</span>
		</label>
		<input
			class="w-full font-mono bg-neutral-600 outline-none py-2 px-3 mt-2 rounded-md"
			id="webhook"
			bind:value={webhookURL}
		/>
	</div>
	{#if $mutation.isError}
		<p class="mb-5">{AxiosResponseErr($mutation.error)?.error_msg}</p>
	{/if}
	{#if $deleteMutation.isError}
		<p class="mb-5">{AxiosResponseErr($deleteMutation.error)?.error_msg}</p>
	{/if}
	<div class="mb-3 flex gap-4">
		<Button on:click={onSave}>Save</Button>
		{#if isConfigured}
			<Button on:click={onDelete}>Delete</Button>
		{/if}
	</div>
</IntegrationModalWithList>

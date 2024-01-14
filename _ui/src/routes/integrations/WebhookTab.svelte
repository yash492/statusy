<script lang="ts">
	import { IntegrationsAPI } from '$lib/apis/integrations';
	import { AxiosResponseErr } from '$lib/helpers/errors';
	import { WEBHOOKS_GET_QUERY_KEY } from '$lib/types/query_keys';
	import { createQuery } from '@tanstack/svelte-query';
	import Webhook from './Webhook.svelte';

	const _integrationAPI = new IntegrationsAPI();

	const query = createQuery({
		queryKey: [WEBHOOKS_GET_QUERY_KEY],
		queryFn: () => _integrationAPI.GetWebhook()
	});

	$: webhook = $query.data?.data.data;
</script>

{#if $query.isLoading}
	<p>Loading</p>
{:else if $query.isError}
	<p>{AxiosResponseErr($query.error)?.error_msg}</p>
{:else}
	<Webhook
		isConfigured={webhook?.is_configured}
		uuid={webhook?.uuid || ''}
		webhookSecret={webhook?.secret}
		webhookURL={webhook?.webhook_url}
	/>
{/if}

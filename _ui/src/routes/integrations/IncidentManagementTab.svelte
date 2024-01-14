<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import Pagerduty from './Pagerduty.svelte';
	import Squadcast from './Squadcast.svelte';
	import { IntegrationsAPI } from '$lib/apis/integrations';
	import { INCIDENT_MANAGEMENT_GET_QUERY_KEY } from '$lib/types/query_keys';
	import { AxiosResponseErr } from '$lib/helpers/errors';

	const _integrationAPI = new IntegrationsAPI();

	const query = createQuery({
		queryKey: [INCIDENT_MANAGEMENT_GET_QUERY_KEY],
		queryFn: () => _integrationAPI.GetIncidentManagement()
	});

	$: squadcast = $query.data?.data.data.squadcast;
	$: pagerduty = $query.data?.data.data.pagerduty;
</script>

{#if $query.isLoading}
	<p>Loading...</p>
{:else if $query.isError}
	<p>{AxiosResponseErr($query.error)?.error_msg}</p>
{:else}
	<Squadcast
		webhookURL={squadcast?.webhook_url}
		isConfigured={squadcast?.is_configured || false}
		uuid={squadcast?.uuid || ''}
	/>
	<Pagerduty
		routingKey={pagerduty?.routing_key}
		isConfigured={pagerduty?.is_configured || false}
		uuid={pagerduty?.uuid || ''}
	/>
{/if}

<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import Discord from './Discord.svelte';
	import Msteams from './Msteams.svelte';
	import Slack from './Slack.svelte';
	import { IntegrationsAPI } from '$lib/apis/integrations';
	import { CHATOPS_GET_ALL_QUERY_KEY } from '$lib/types/query_keys';
	import { AxiosResponseErr } from '$lib/helpers/errors';

	const _integrationAPI = new IntegrationsAPI();

	$: query = createQuery({
		queryFn: () => _integrationAPI.GetChatOps(),
		queryKey: [CHATOPS_GET_ALL_QUERY_KEY]
	});

	$: slackResp = $query.data?.data.data.slack;
	$: msteamsResp = $query.data?.data.data.msteams;
	$: discordResp = $query.data?.data.data.discord;
</script>

{#if $query.isLoading}
	<p>Loading...</p>
{:else if $query.isError}
	<p>{AxiosResponseErr($query.error)?.error_msg}</p>
{:else}
	<Slack
		isConfigured={slackResp?.is_configured || false}
		uuid={slackResp?.uuid}
		webhookURL={slackResp?.webhook_url}
	/>
	<Discord
		isConfigured={discordResp?.is_configured || false}
		uuid={discordResp?.uuid}
		webhookURL={discordResp?.webhook_url}
	/>
	<Msteams
		isConfigured={msteamsResp?.is_configured || false}
		uuid={msteamsResp?.uuid}
		webhookURL={msteamsResp?.webhook_url}
	/>
{/if}

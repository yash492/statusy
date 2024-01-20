<script>
	import Tab from '$lib/components/tabs/Tab.svelte';
	import Tabs from '$lib/components/tabs/Tabs.svelte';
	import ChatopsTab from './ChatopsTab.svelte';
	import IncidentManagementTab from './IncidentManagementTab.svelte';
	import WebhookTab from './WebhookTab.svelte';
</script>

<div class="pt-12 md:ml-6">
	<Tabs>
		<Tab title="Incident Management" id="1">
			<IncidentManagementTab />
		</Tab>
		<Tab title="Chat Ops" id="2">
			<ChatopsTab />
		</Tab>
		<Tab title="Webhook" id="3">
			<WebhookTab />
		</Tab>
	</Tabs>
</div>

<script>
	import Tab from '$lib/components/tabs/Tab.svelte';
	import Tabs from '$lib/components/tabs/Tabs.svelte';
	import Chatops from './Chatops.svelte';
	import IncidentManagement from './IncidentManagement.svelte';
	import Webhook from './Webhook.svelte';
</script>

<div class="pt-12 ml-6">
	<Tabs>
		<Tab title="Incident Management" id="1">
			<IncidentManagement />
		</Tab>
		<Tab title="Chat Ops" id="2">
			<Chatops />
		</Tab>
		<Tab title="Webhook" id="3">
			<Webhook />
		</Tab>
	</Tabs>
</div>

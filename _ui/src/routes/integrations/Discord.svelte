<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import IntegrationModalWithList from '$lib/components/integration_component_list/IntegrationModalWithList.svelte';
	let webhookURL = '';
	$: console.log(webhookURL);
</script>

<IntegrationModalWithList isIntegrated={false} name="Discord" modalTitle="Discord Integration">
	<div class="mt-10 mb-8">
		<label class="font-medium" for="webhook"
			>Your Webhook URL <span class="ml-1"
				>(<a
					class="text-sm hover:text-blue-400 text-neutral-300 hover:underline"
					href="https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks"
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
	<div class="mb-3 flex gap-4">
		<Button>Add</Button>
		<Button>Delete</Button>
	</div>
</IntegrationModalWithList>

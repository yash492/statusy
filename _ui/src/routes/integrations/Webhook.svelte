<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import IntegrationModalWithList from '$lib/components/integration_component_list/IntegrationModalWithList.svelte';
	let webhookURL = '';
	let webhookSecret = '';
</script>

<IntegrationModalWithList isIntegrated={false} name="Webhook" modalTitle={'Webhook Integration'}>
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
			class="w-full bg-neutral-600 outline-none py-2 mt-2 px-3 rounded-md"
			id="secret"
			bind:value={webhookSecret}
		/>
	</div>

	<div class="mb-3 flex gap-4">
		<Button>Add</Button>
		<Button>Delete</Button>
	</div>
</IntegrationModalWithList>

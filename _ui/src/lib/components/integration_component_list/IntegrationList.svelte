<script lang="ts">
	import Button from '../button/Button.svelte';
	export let showModal: Boolean;
	export let name: string;
	export let isIntegrated: boolean;
	let dotColor = 'text-red-500';
	let integratedDisplay = 'Not Integrated';
	$: if (isIntegrated) {
		dotColor = 'text-green-500';
		integratedDisplay = 'Integrated';
	}
</script>

<div class="w-full md:w-1/2 py-1 flex justify-between items-center rounded-lg border-neutral-700">
	<p class="font-bold">{name}</p>
	<div class="flex items-center gap-4">
		<Button on:click={() => (showModal = true)}>Configure</Button>
		<div class="flex items-center gap-1">
			<p class={`text-3xl ${dotColor}`}>&#x2022</p>
			<p class="text-sm font-semibold">{integratedDisplay}</p>
		</div>
	</div>
</div>

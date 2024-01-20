<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import Button from '../button/Button.svelte';
	import { AdjustmentsHorizontal } from '@steeze-ui/heroicons';
	export let showModal: Boolean;
	export let name: string;
	export let isIntegrated: boolean;
	let dotColor = 'text-red-500';
	let integratedDisplay = 'Not Integrated';
	$: if (isIntegrated) {
		dotColor = 'text-green-500';
		integratedDisplay = 'Integrated';
	}

	$: if (!isIntegrated) {
		dotColor = 'text-red-500';
		integratedDisplay = 'Not Integrated';
	}
</script>

<div class="w-full md:w-1/2 py-1 flex justify-between items-center rounded-lg border-neutral-700">
	<p class="font-bold text-sm">{name}</p>
	<div class="flex items-center md:gap-4 gap-2 justify-center">
		<div class="hidden md:block">
			<Button class="text-sm" on:click={() => (showModal = true)}>Configure</Button>
		</div>
		<div class="md:hidden">
			<button on:click={() => (showModal = true)}>
				<Icon src={AdjustmentsHorizontal} size="16px" />
			</button>
		</div>
		<div class="flex items-center gap-1 w-32">
			<p class={`text-3xl ${dotColor}`}>&#x2022</p>
			<p class="text-xs font-semibold md:text-sm">{integratedDisplay}</p>
		</div>
	</div>
</div>

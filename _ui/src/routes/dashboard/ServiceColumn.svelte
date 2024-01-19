<script lang="ts">
	import { goto } from '$app/navigation';
	import { AdjustmentsHorizontal, XCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	export let serviceName: string;
	export let incidentName: string;
	export let incidentLink: string;
	export let subscriptionUUID: string;
</script>

<div class="flex items-center justify-between mr-2">
	<div class="w-4/5">
		<p>
			<a class="hover:text-blue-300" href={`/subscriptions/${subscriptionUUID}/incidents`}>
				{serviceName}</a
			>
		</p>
		<p class="text-xs text-neutral-300">
			<a target="_blank" class="hover:text-blue-300" href={incidentLink}>{incidentName}</a>
		</p>
	</div>

	<div class="flex items-center gap-1">
		<button on:click={() => goto(`subscriptions/${subscriptionUUID}/edit`)}>
			<Icon src={AdjustmentsHorizontal} size={'23px'} />
		</button>
		<button>
			<Icon src={XCircle} size={'23px'} />
		</button>
	</div>
</div>

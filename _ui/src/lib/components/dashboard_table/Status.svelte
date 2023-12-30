<script lang="ts">
	export let value: string;

	let cssClasses = 'border-green-600 bg-green-950';
	if (value === 'Outage') {
		cssClasses = 'border-red-600 bg-red-950';
	}
</script>

<div class={`uppercase font-bold text-sm text-center py-1 rounded-md border-2 ${cssClasses}`}>
	{value}
</div>

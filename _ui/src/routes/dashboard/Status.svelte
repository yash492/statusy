<script lang="ts">
	export let isDown: boolean;

	let cssClasses = 'border-green-600 bg-green-950';
	if (isDown) {
		cssClasses = 'border-red-600 bg-red-950';
	}
</script>

<div class={`uppercase font-bold text-sm text-center py-1 rounded-md border-2 ${cssClasses}`}>
	{isDown ? 'down' : 'up'}
</div>

<script lang="ts">
	export let type: 'button' | 'reset' | 'submit' = 'button';
</script>

<button
	{type}
	class="w-fit text-sm font-bold px-3 py-1 border-2 rounded-md bg-purple-950 border-purple-600 hover:bg-purple-800"
	on:click><slot /></button
>

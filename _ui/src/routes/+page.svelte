<script lang="ts">
	const { data } = $props();
</script>

<div class="text-centershadow-lg mx-auto mt-16 max-w-md items-center rounded-xl border p-6">
	<h2 class="font-bold text-white">Error -</h2>
	<p class="font-mono">
		{data.err?.message}
	</p>
</div>

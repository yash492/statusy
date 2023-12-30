<script>
	import { goto } from '$app/navigation';
	import Button from '$lib/components/button/Button.svelte';
	import Table from '$lib/components/dashboard_table/Table.svelte';
</script>

<div class="w-fit mx-auto md:w-2/3 lg:w-2/5 pt-9 md:pt-12">
	<div class="flex items-center justify-between">
		<h1 class="font-bold text-lg md:text-2xl">Montoring 3 services</h1>
		<Button on:click={() => goto('/dashboard/add-service')}>Add Service</Button>
	</div>
	<Table />
</div>

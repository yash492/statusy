<script lang="ts">
	import Select from 'svelte-select';
	import Button from '../button/Button.svelte';
</script>

<form>
	<h1 class="font-bold text-xl pt-5">Add Service</h1>
	<div class="w-[30rem] mt-4">
		<div>
			<label for="service-select" class="py-1 block">Select a Service</label>
			<div class="block">
				<Select
					id="service-select"
					--background="rgb(38,38,38)"
					--list-background="rgb(38,38,38)"
					--border="rgb(115,115,115)"
					--border-hover="rgb(115,115,115)"
					--placeholder-color="rgb(200,200,200)"
					placeholder="Squadcast"
				></Select>
			</div>
		</div>
		<div>
			<label for="components-select" class="mt-5 block py-1">Select Components</label>
			<div class="border h-60 bg-neutral-800 rounded-md border-neutral-500"></div>
		</div>
		<div class="mt-5">
			<Button type="submit">Add Service</Button>
		</div>
	</div>
</form>

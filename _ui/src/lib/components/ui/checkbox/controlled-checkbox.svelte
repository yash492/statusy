<script lang="ts">
	import Checkbox from "./checkbox.svelte";

	let {
		id,
		checked = false,
		indeterminate = false,
		disabled = false,
		class: className,
		onchange
	} = $props<{
		id: string;
		checked?: boolean;
		indeterminate?: boolean;
		disabled?: boolean;
		class?: string;
		onchange?: (checked: boolean | "indeterminate") => void;
	}>();

	let localChecked = $state(false);
	let localIndeterminate = $state(false);

	$effect(() => {
		localChecked = checked;
	});

	$effect(() => {
		localIndeterminate = indeterminate;
	});
</script>

<Checkbox
	{id}
	bind:checked={localChecked}
	bind:indeterminate={localIndeterminate}
	{disabled}
	onCheckedChange={(val) => {
		onchange?.(val);
	}}
	class={className}
/>

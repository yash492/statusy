<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';

	let {
		name = $bindable(''),
		description = $bindable(''),
		isDefault = $bindable(false),
		showDefaultCheckbox = false,
		submitText = 'Save Changes',
		cancelText = 'Cancel',
		namePlaceholder = '',
		descriptionPlaceholder = '',
		onsubmit,
		oncancel
	}: {
		name?: string;
		description?: string;
		isDefault?: boolean;
		showDefaultCheckbox?: boolean;
		submitText?: string;
		cancelText?: string;
		namePlaceholder?: string;
		descriptionPlaceholder?: string;
		onsubmit: () => void;
		oncancel: () => void;
	} = $props();
</script>

<div class="grid gap-4 py-4">
	<div class="grid gap-2">
		<Label for="view-name" class="text-sm font-semibold text-zinc-300">Name</Label>
		<Input
			id="view-name"
			bind:value={name}
			placeholder={namePlaceholder}
			class="border-zinc-800 bg-zinc-900/50 text-white placeholder-zinc-500 focus-visible:ring-zinc-700"
		/>
	</div>
	<div class="grid gap-2">
		<Label for="view-description" class="text-sm font-semibold text-zinc-300">Description</Label>
		<Textarea
			id="view-description"
			bind:value={description}
			placeholder={descriptionPlaceholder}
			class="min-h-[60px] border-zinc-800 bg-zinc-900/50 text-white placeholder-zinc-500 focus-visible:ring-zinc-700"
		/>
	</div>

	{#if showDefaultCheckbox}
		<label class="mt-2 flex cursor-pointer items-center gap-2 select-none">
			<input
				type="checkbox"
				bind:checked={isDefault}
				class="size-4 rounded border-zinc-800 bg-zinc-900 text-white accent-zinc-300 transition-colors focus:ring-0"
			/>
			<span class="text-sm text-zinc-300">Make this the default view</span>
		</label>
	{/if}

	<div class="flex justify-end gap-2 pt-2">
		<button
			class="border-zinc-850 cursor-pointer rounded-lg border bg-zinc-900 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-zinc-800"
			onclick={oncancel}
		>
			{cancelText}
		</button>
		<button
			class="cursor-pointer rounded-lg bg-white px-4 py-2 text-base font-bold text-zinc-950 transition-colors hover:bg-zinc-200"
			onclick={onsubmit}
		>
			{submitText}
		</button>
	</div>
</div>

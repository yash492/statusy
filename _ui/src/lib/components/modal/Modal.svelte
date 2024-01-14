<script lang="ts">
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	let dialog: HTMLDialogElement;
	export let showModal: boolean;
	export let modalTitle: string;

	$: if (dialog && showModal) {
		dialog.showModal();
	}

	$: if (dialog && !showModal) {
		dialog.close();
	}
</script>

<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->

<dialog
	on:click|self={() => {
		showModal = false;
	}}
	on:introstart
	on:outroend
	bind:this={dialog}
	class="w-4/5 md:w-1/2 rounded-md shadow-md backdrop:bg-neutral-800 backdrop:opacity-60 bg-neutral-800 text-white border-2 border-neutral-700 mt-20"
>
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div on:click|stopPropagation class="w-full h-full px-3 py-2">
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-semibold">{modalTitle}</h2>
			<button
				class="hover:bg-neutral-600 hover:rounded-sm"
				on:click={() => {
					showModal = false;
				}}
			>
				<Icon src={XMark} size="23px"></Icon>
			</button>
		</div>
		<slot />
	</div>
</dialog>

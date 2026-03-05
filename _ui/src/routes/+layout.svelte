<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import Input from '$lib/components/ui/input/input.svelte';
	import * as Kbd from '$lib/components/ui/kbd';
	import { Label } from '$lib/components/ui/label';
	import { SearchIcon } from '@lucide/svelte';
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';
	let { children } = $props();

	let modifierKeyPrefix = $state('');

	modifierKeyPrefix =
		navigator.userAgent.toLowerCase().includes('mac') ||
		navigator.userAgent.toLowerCase().includes('iphone')
			? '⌘ + K' // command key
			: 'Ctrl + K'; // control key
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<nav class="h-11 w-full">
	<div class="mx-auto mt-3 flex w-4/5 justify-between">
		<div class="flex items-center gap-2">
			<img src="/logos/statusy.svg" height="30rem" width="30rem" alt="statusy_logo" />
			<p class="text-2xl font-bold">Statusy</p>
		</div>
		<div class="flex gap-2">
			<Dialog.Root>
				<form class="">
					<Dialog.Trigger
						type="button"
						class={buttonVariants({
							variant: 'outline',
							class: 'cursor-pointer text-gray-200',
							size: 'lg'
						})}
					>
						<SearchIcon />
						<span class="hidden md:block"
							>Search Status Pages <Kbd.Root class="ml-10">{modifierKeyPrefix}</Kbd.Root></span
						>
					</Dialog.Trigger>
					<Dialog.Content class="sm:max-w-xl">
						<div class="mb-28 grid gap-4">
							<div class="grid gap-3">
								<Label for="search-statuspage" class="text-xl">Search Status Page</Label>
								<Input id="search-statuspage" name="search-statuspage" />
							</div>
						</div>
					</Dialog.Content>
				</form>
			</Dialog.Root>
		</div>
	</div>
</nav>

{@render children()}

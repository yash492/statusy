<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as Kbd from '$lib/components/ui/kbd';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import { SearchIcon } from '@lucide/svelte';

	let modifierKeyPrefix = $state('');
	let isSearchDialogOpen = $state(false);

	modifierKeyPrefix =
		navigator.userAgent.toLowerCase().includes('mac') ||
		navigator.userAgent.toLowerCase().includes('iphone')
			? '⌘ + K' // command key
			: 'Ctrl + K'; // control key

	function handleWindowKeydown(event: KeyboardEvent) {
		if (event.key === 'k' && (event.ctrlKey || event.metaKey)) {
			event.preventDefault(); // Prevent default browser save action
			isSearchDialogOpen = true;
		}
	}
</script>

<svelte:window onkeydown={handleWindowKeydown} />
<nav class="mt-8 h-28 w-full">
	<div class="mx-auto w-4/5">
		<div class=" flex justify-between">
			<div class="flex items-center gap-2">
				<img src="/logos/statusy.svg" height="30rem" width="30rem" alt="statusy_logo" />
				<p class="text-2xl font-bold">Statusy</p>
			</div>
			<div class="flex gap-2">
				<Dialog.Root open={isSearchDialogOpen} onOpenChange={() => (isSearchDialogOpen = false)}>
					<form class="">
						<Dialog.Trigger
							type="button"
							class={buttonVariants({
								variant: 'outline',
								class: 'cursor-text text-gray-200',
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
		<Separator class="mt-8" />
	</div>
</nav>

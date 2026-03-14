<script lang="ts">
	import { goto } from '$app/navigation';
	import { StatuspageApi, type Statuspage } from '$lib/api/statuspage/statuspage';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as Kbd from '$lib/components/ui/kbd';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import { SearchIcon } from '@lucide/svelte';

	const statuspageApi = new StatuspageApi();
	const SEARCH_DEBOUNCE_MS = 250;

	let modifierKeyPrefix = $state('');
	let isSearchDialogOpen = $state(false);
	let query = $state('');
	let results = $state<Statuspage[]>([]);
	let isLoading = $state(false);
	let errorMessage = $state<string | null>(null);
	let activeIndex = $state(0);
	let hasSearchResponse = $state(false);

	let searchInput = $state<HTMLInputElement | null>(null);
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;
	let latestRequestId = 0;

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

	function resetSearchState() {
		if (debounceTimer) {
			clearTimeout(debounceTimer);
			debounceTimer = null;
		}

		query = '';
		results = [];
		isLoading = false;
		errorMessage = null;
		activeIndex = 0;
		hasSearchResponse = false;
		latestRequestId += 1;
	}

	function handleDialogOpenChange(open: boolean) {
		isSearchDialogOpen = open;

		if (!open) {
			resetSearchState();
		}
	}

	async function navigateToStatuspage(statuspage: Statuspage) {
		const params = new URLSearchParams({
			page: '1',
			page_size: '10',
			type: 'incidents'
		});

		const href = `/statuspages/${encodeURIComponent(statuspage.slug)}/events?${params.toString()}`;
		isSearchDialogOpen = false;
		resetSearchState();
		await goto(href);
	}

	function handleSearchInputKeydown(event: KeyboardEvent) {
		if (event.key === 'ArrowDown') {
			event.preventDefault();
			if (results?.length > 0) {
				activeIndex = (activeIndex + 1) % results.length;
			}
			return;
		}

		if (event.key === 'ArrowUp') {
			event.preventDefault();
			if (results.length > 0) {
				activeIndex = (activeIndex - 1 + results.length) % results.length;
			}
			return;
		}

		if (event.key === 'Enter' && results.length > 0) {
			event.preventDefault();
			const selected = results[activeIndex];
			if (selected) {
				void navigateToStatuspage(selected);
			}
		}
	}

	$effect(() => {
		if (!isSearchDialogOpen) {
			return;
		}

		const focusTimer = setTimeout(() => {
			searchInput?.focus();
		}, 0);

		return () => {
			clearTimeout(focusTimer);
		};
	});

	$effect(() => {
		if (!isSearchDialogOpen) {
			return;
		}

		const normalizedQuery = query.trim();
		if (debounceTimer) {
			clearTimeout(debounceTimer);
			debounceTimer = null;
		}

		if (!normalizedQuery) {
			results = [];
			isLoading = false;
			errorMessage = null;
			activeIndex = 0;
			hasSearchResponse = false;
			return;
		}

		hasSearchResponse = false;

		debounceTimer = setTimeout(async () => {
			const requestId = ++latestRequestId;
			isLoading = true;
			errorMessage = null;

			try {
				const pages = await statuspageApi.list(normalizedQuery);
				if (requestId !== latestRequestId) {
					return;
				}

				results = pages;
				activeIndex = pages.length === 0 ? 0 : Math.min(activeIndex, pages.length - 1);
				hasSearchResponse = true;
			} catch {
				if (requestId !== latestRequestId) {
					return;
				}

				results = [];
				errorMessage = 'Could not load status pages.';
				hasSearchResponse = false;
			} finally {
				if (requestId === latestRequestId) {
					isLoading = false;
				}
			}
		}, SEARCH_DEBOUNCE_MS);

		return () => {
			if (debounceTimer) {
				clearTimeout(debounceTimer);
				debounceTimer = null;
			}
		};
	});
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
				<Dialog.Root open={isSearchDialogOpen} onOpenChange={handleDialogOpenChange}>
					<div>
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
						<Dialog.Content class="top-[35%] translate-y-[-35%] sm:max-w-xl">
							<div class="grid gap-4">
								<div class="grid gap-3">
									<Label for="search-statuspage" class="text-xl">Search Status Page</Label>
									<Input
										id="search-statuspage"
										name="search-statuspage"
										bind:value={query}
										bind:ref={searchInput}
										onkeydown={handleSearchInputKeydown}
										placeholder=""
										autocomplete="off"
										aria-autocomplete="list"
										aria-controls="statuspage-search-results"
									/>
								</div>

								<div
									class="max-h-72 overflow-auto rounded-md border"
									id="statuspage-search-results"
									role="listbox"
								>
									{#if isLoading}
										<p class="px-3 py-2 text-sm text-muted-foreground">Searching...</p>
									{:else if errorMessage}
										<p class="px-3 py-2 text-sm text-destructive">{errorMessage}</p>
									{:else if query.trim().length === 0}
										<p class="px-3 py-2 text-sm text-muted-foreground">Start typing to search.</p>
									{:else if !hasSearchResponse}
										<p class="px-3 py-2 text-sm text-muted-foreground">Searching...</p>
									{:else if results.length === 0}
										<p class="px-3 py-2 text-sm text-muted-foreground">No status pages found.</p>
									{:else}
										<ul>
											{#each results as result, index (result.id)}
												<li>
													<button
														type="button"
														class={`w-full cursor-pointer px-3 py-2 text-left ${
															index === activeIndex
																? 'bg-accent text-accent-foreground'
																: 'hover:bg-accent/60'
														}`}
														onmouseenter={() => (activeIndex = index)}
														onclick={() => navigateToStatuspage(result)}
														role="option"
														aria-selected={index === activeIndex}
													>
														<p class="text-sm font-medium">{result.name}</p>
													</button>
												</li>
											{/each}
										</ul>
									{/if}
								</div>
							</div>
						</Dialog.Content>
					</div>
				</Dialog.Root>
			</div>
		</div>
		<Separator class="mt-8" />
	</div>
</nav>

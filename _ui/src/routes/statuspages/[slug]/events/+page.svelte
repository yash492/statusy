<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tabs from '$lib/components/ui/tabs';
	import IncidentsTable, { type Incident } from '$lib/derivedcomponents/IncidentsTable.svelte';
	import CheckIcon from '@lucide/svelte/icons/check';
	import ClipboardIcon from '@lucide/svelte/icons/clipboard';
	import RssIcon from '@lucide/svelte/icons/rss';
	import SlackIcon from '@lucide/svelte/icons/slack';
	import type { PaginationState } from '@tanstack/table-core';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
	type TabType = 'incidents' | 'scheduled-maintenance';
	type SubscribeTab = 'rss' | 'slack';
	type CopyKey = 'rss' | 'atom' | 'slack';

	let isSubscribeDialogOpen = $state(false);
	let subscribeTab = $state<SubscribeTab>('rss');
	let copiedKey = $state<CopyKey | null>(null);
	let copyError = $state<string | null>(null);
	let copyResetTimer: ReturnType<typeof setTimeout> | null = null;

	const PAGE_SIZE = 10;
	const initialPagination: PaginationState = {
		pageIndex: Math.max(0, data.page - 1),
		pageSize: data.pageSize
	};

	function toIncidents(raw: typeof data.resp.incidents): Incident[] {
		return raw.map((incident) => ({
			created_at: incident.provider_created_at,
			id: incident.id,
			status: incident.status,
			title: incident.title,
			incident_url: incident.incident_url
		}));
	}

	const incidentData = $derived(toIncidents(data.resp.incidents));
	const feedBaseUrl = $derived(
		`${$page.url.origin}/statuspages/${encodeURIComponent(data.resp.statuspage.slug)}`
	);
	const feedRssPath = $derived(`${feedBaseUrl}/feed.rss`);
	const feedAtomPath = $derived(`${feedBaseUrl}/feed.atom`);
	const slackSnippet = $derived(`/feed subscribe ${feedRssPath}`);

	// If first page is full, we don't know the total — use MAX to keep Next enabled.
	// Once a page returns fewer rows than pageSize, we know the exact total.
	const rowCount = $derived(
		data.resp.incidents.length < PAGE_SIZE
			? (data.page - 1) * data.pageSize + data.resp.incidents.length
			: Number.MAX_SAFE_INTEGER
	);

	async function onPageChange(pagination: PaginationState) {
		const params = new URLSearchParams(window.location.search);
		params.set('page', String(pagination.pageIndex + 1));
		params.set('page_size', String(pagination.pageSize));
		await goto(`?${params.toString()}`, {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	}

	async function onTabChange(type: TabType) {
		const params = new URLSearchParams(window.location.search);
		params.set('type', type);
		params.set('page', '1');
		params.set('page_size', String(data.pageSize));

		await goto(`?${params.toString()}`, {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	}

	function onSubscribeDialogOpenChange(open: boolean) {
		isSubscribeDialogOpen = open;
		if (!open) {
			subscribeTab = 'rss';
			copiedKey = null;
			copyError = null;
			if (copyResetTimer) {
				clearTimeout(copyResetTimer);
				copyResetTimer = null;
			}
		}
	}

	async function copyText(value: string, key: CopyKey) {
		try {
			copyError = null;
			await navigator.clipboard.writeText(value);
			copiedKey = key;

			if (copyResetTimer) {
				clearTimeout(copyResetTimer);
			}

			copyResetTimer = setTimeout(() => {
				copiedKey = null;
				copyResetTimer = null;
			}, 1800);
		} catch {
			copyError = 'Could not copy. Please copy manually.';
		}
	}
</script>

<svelte:head>
	<link rel="alternate" type="application/rss+xml" title="{data.resp.statuspage.name} - RSS Feed" href={feedRssPath} />
	<link rel="alternate" type="application/atom+xml" title="{data.resp.statuspage.name} - Atom Feed" href={feedAtomPath} />
</svelte:head>

<div class="mx-auto w-4/5">
	<div class="w-full">
		<div class="mb-6 flex justify-between md:mb-4">
			<div class="mb-4 flex w-fit items-center gap-2">
				<p class="text-xl font-bold">{data.resp.statuspage.name}</p>
			</div>
			<div>
				<Dialog.Root open={isSubscribeDialogOpen} onOpenChange={onSubscribeDialogOpenChange}>
					<Button class="cursor-pointer" onclick={() => (isSubscribeDialogOpen = true)}>
						Subscribe to Updates
					</Button>
					<Dialog.Content class="sm:max-w-2xl">
						<div class="grid gap-4">
							<Dialog.Header>
								<Dialog.Title>Subscribe to Updates</Dialog.Title>
							</Dialog.Header>

							<Tabs.Root
								value={subscribeTab}
								onValueChange={(value) => (subscribeTab = value as SubscribeTab)}
							>
								<Tabs.List class="grid w-full grid-cols-2 sm:w-[30%]">
									<Tabs.Trigger class="cursor-pointer gap-1.5" value="rss">
										<RssIcon class="size-4" />
										RSS
									</Tabs.Trigger>
									<Tabs.Trigger class="cursor-pointer gap-1.5" value="slack">
										<SlackIcon class="size-4" />
										Slack
									</Tabs.Trigger>
								</Tabs.List>

								<Tabs.Content value="rss" class="grid gap-3 pt-2">
									<div class="rounded-md border p-3">
										<p class="mb-2 text-sm font-medium">RSS Feed</p>
										<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
											<code class="block rounded bg-muted px-2 py-1 text-xs break-all"
												>{feedRssPath}</code
											>
											<Button
												type="button"
												variant="outline"
												class="h-9 w-9 shrink-0 cursor-pointer p-0"
												onclick={() => copyText(feedRssPath, 'rss')}
												title={copiedKey === 'rss' ? 'Copied' : 'Copy RSS URL'}
												aria-label={copiedKey === 'rss' ? 'Copied RSS URL' : 'Copy RSS URL'}
											>
												{#if copiedKey === 'rss'}
													<CheckIcon class="size-4" />
												{:else}
													<ClipboardIcon class="size-4" />
												{/if}
												<span class="sr-only"
													>{copiedKey === 'rss' ? 'Copied RSS URL' : 'Copy RSS URL'}</span
												>
											</Button>
										</div>
									</div>

									<div class="rounded-md border p-3">
										<p class="mb-2 text-sm font-medium">Atom Feed</p>
										<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
											<code class="block rounded bg-muted px-2 py-1 text-xs break-all"
												>{feedAtomPath}</code
											>
											<Button
												type="button"
												variant="outline"
												class="h-9 w-9 shrink-0 cursor-pointer p-0"
												onclick={() => copyText(feedAtomPath, 'atom')}
												title={copiedKey === 'atom' ? 'Copied' : 'Copy Atom URL'}
												aria-label={copiedKey === 'atom' ? 'Copied Atom URL' : 'Copy Atom URL'}
											>
												{#if copiedKey === 'atom'}
													<CheckIcon class="size-4" />
												{:else}
													<ClipboardIcon class="size-4" />
												{/if}
												<span class="sr-only"
													>{copiedKey === 'atom' ? 'Copied Atom URL' : 'Copy Atom URL'}</span
												>
											</Button>
										</div>
									</div>
								</Tabs.Content>

								<Tabs.Content value="slack" class="grid gap-3 pt-2">
									<p class="text-sm text-muted-foreground">
										To receive live status updates in Slack, copy and paste the text below into the
										Slack channel of your choice.
									</p>
									<div></div>
									<div class="flex items-center justify-between rounded-md border p-3">
										<code class="block rounded bg-muted px-2 py-1 text-xs break-all"
											>{slackSnippet}</code
										>
										<div class="flex justify-end">
											<Button
												type="button"
												variant="outline"
												class="h-9 w-9 shrink-0 cursor-pointer p-0"
												onclick={() => copyText(slackSnippet, 'slack')}
												title={copiedKey === 'slack' ? 'Copied' : 'Copy Slack snippet'}
												aria-label={copiedKey === 'slack'
													? 'Copied Slack snippet'
													: 'Copy Slack snippet'}
											>
												{#if copiedKey === 'slack'}
													<CheckIcon class="size-4" />
												{:else}
													<ClipboardIcon class="size-4" />
												{/if}
												<span class="sr-only"
													>{copiedKey === 'slack'
														? 'Copied Slack snippet'
														: 'Copy Slack snippet'}</span
												>
											</Button>
										</div>
									</div>
								</Tabs.Content>
							</Tabs.Root>

							{#if copyError}
								<p class="text-sm text-destructive">{copyError}</p>
							{/if}
						</div>
					</Dialog.Content>
				</Dialog.Root>
			</div>
		</div>

		<div>
			<div>
				<Tabs.Root
					value={data.type}
					onValueChange={(value) => onTabChange(value as TabType)}
					class="w-full"
				>
					<Tabs.List>
						<Tabs.Trigger value="incidents">Incidents</Tabs.Trigger>
						<Tabs.Trigger value="scheduled-maintenances">Scheduled Maintenances</Tabs.Trigger>
					</Tabs.List>
					<Tabs.Content value="incidents">
						<IncidentsTable
							data={incidentData}
							{rowCount}
							paginationState={initialPagination}
							{onPageChange}
						/>
					</Tabs.Content>
					<Tabs.Content value="scheduled-maintenances">
						<!-- <IncidentsTable /> -->
					</Tabs.Content>
				</Tabs.Root>
			</div>
		</div>
	</div>
</div>

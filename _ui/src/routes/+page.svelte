<script lang="ts">
	import { onMount } from 'svelte';

	let modifierKeyPrefix = $state('⌘ + K');

	onMount(() => {
		const isMac =
			navigator.userAgent.toLowerCase().includes('mac') ||
			navigator.userAgent.toLowerCase().includes('iphone');
		modifierKeyPrefix = isMac ? '⌘ + K' : 'Ctrl + K';
	});
</script>

<svelte:head>
	<title>Statusy | Status Page Aggregator</title>
	<meta
		name="description"
		content="Aggregated real-time status pages, incident history, RSS feeds and Slack integration for your system dependencies."
	/>
</svelte:head>

<div class="relative w-full overflow-x-hidden text-white">
	<div class="mx-auto flex w-3/5 flex-col justify-between pt-12 md:pt-20">
		<!-- Hero Section -->
		<header class="relative mb-16 flex flex-col items-center text-center">
			<h1 class="mb-6 max-w-3xl text-4xl font-extrabold tracking-tight text-white sm:text-6xl">
				Keep tabs on all your critical dependencies.
			</h1>

			<p class="mb-8 max-w-2xl text-lg leading-relaxed text-zinc-400">
				Statusy collects, normalizes, and syndicates real-time status page data, incidents, and
				maintenance announcements from top cloud services in one unified dashboard.
			</p>

			<!-- Search Shortcut Text -->
			<p class="mt-8 text-sm text-zinc-500">
				Press
				<kbd
					class="rounded border border-zinc-800 bg-zinc-950 px-1.5 py-0.5 font-mono text-xs text-zinc-400"
				>
					{modifierKeyPrefix}
				</kbd>
				to search anywhere
			</p>
		</header>
	</div>
</div>

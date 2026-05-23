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
	<title>Statusy | The Open Source Status Page Aggregator</title>
	<meta
		name="description"
		content="Aggregated real-time status pages, incident history, RSS feeds and Slack integration for your system dependencies."
	/>
</svelte:head>

<div class="relative min-h-screen w-full overflow-x-hidden bg-black text-white">
	<div class="mx-auto flex min-h-screen w-4/5 flex-col justify-between pt-12 md:pt-20">
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

		<!-- Footer -->
		<footer
			class="flex flex-col items-center justify-between gap-4 border-t border-zinc-800 py-12 text-sm text-zinc-400 md:flex-row"
		>
			<div class="flex items-center gap-2">
				<img src="/logos/statusy.svg" height="24" width="24" alt="logo" />
				<span class="font-bold text-white">Statusy</span>
				<span>© {new Date().getFullYear()}</span>
			</div>
		</footer>
	</div>
</div>

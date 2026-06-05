<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import Navbar from '$lib/components/Navbar.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';
	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<Toaster />

<div class="flex min-h-screen flex-col bg-zinc-950">
	<Navbar />
	<main class="mx-auto mt-4 mb-10 w-full max-w-7xl px-4 sm:px-6 lg:px-24">
		{@render children?.()}
	</main>
</div>

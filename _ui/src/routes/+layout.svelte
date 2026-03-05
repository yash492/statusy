<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	import Separator from '$lib/components/ui/separator/separator.svelte';
	import Navbar from '$lib/derivedcomponents/Navbar.svelte';
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';
	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />

<Navbar />
{@render children()}

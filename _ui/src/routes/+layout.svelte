<script lang="ts">
	import Navbar from '$lib/components/navbar/Navbar.svelte';
	import '../app.css';
	import axios from 'axios';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { browser } from '$app/environment';
	import { Toaster } from 'svelte-french-toast';

	axios.defaults.baseURL = '/api/';
	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser
			}
		}
	});
</script>

<QueryClientProvider client={queryClient}>
	<Toaster />
	<Navbar />
	<main class="mt-20 mx-5">
		<slot />
	</main>
</QueryClientProvider>

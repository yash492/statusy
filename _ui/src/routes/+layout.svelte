<script lang="ts">
	import Navbar from '$lib/components/navbar/Navbar.svelte';
	import '../app.css';
</script>

<Navbar />

<main class="mt-20 mx-5">
	<slot />
</main>

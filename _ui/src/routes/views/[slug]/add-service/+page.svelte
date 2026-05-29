<script lang="ts">
	import { page } from '$app/state';
	import { ServiceForm } from '$lib/derivedcomponents';

	const slug = $derived(page.params.slug ?? '');
</script>

<ServiceForm mode="add" {slug} />

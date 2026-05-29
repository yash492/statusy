<script lang="ts">
	import { page } from '$app/state';
	import { ServiceForm } from '$lib/derivedcomponents';

	const slug = $derived(page.params.slug ?? '');
	const serviceSlug = $derived(page.params.serviceSlug);
</script>

<ServiceForm mode="edit" {slug} {serviceSlug} />

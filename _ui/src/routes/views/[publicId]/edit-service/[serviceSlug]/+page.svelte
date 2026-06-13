<script lang="ts">
	import { page } from '$app/state';
	import { ServiceForm } from '$lib/components/statusy';

	const publicId = $derived(page.params.publicId ?? '');
	const serviceSlug = $derived(page.params.serviceSlug);
</script>

<ServiceForm mode="edit" {publicId} {serviceSlug} />

<script lang="ts">
	import { page } from '$app/state';
	import { ServiceForm } from '$lib/derivedcomponents';

	const publicId = $derived(page.params.publicId ?? '');
</script>

<ServiceForm mode="add" {publicId} />

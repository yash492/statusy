<script lang="ts">
	import Modal from '../modal/Modal.svelte';
	import IntegrationList from './IntegrationList.svelte';
	export let showModal: boolean = false;
	export let name: string;
	export let isIntegrated: boolean;
	export let modalTitle: string;
</script>

<IntegrationList {isIntegrated} {name} bind:showModal />
<Modal bind:showModal {modalTitle}>
	<slot />
</Modal>

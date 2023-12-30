<script lang="ts">
	import type { NavbarItems } from '$lib/components/navbar/navbar';
	export let navbarLinks: NavbarItems[];
	export let isMobileNavbarOpen = false;

	function closeNavDropdown() {
		if (isMobileNavbarOpen) {
			isMobileNavbarOpen = false;
		}
	}
</script>

{#each navbarLinks as navbarLink}
	<a href={navbarLink.link} on:click={closeNavDropdown}>
		<ul
			class="list-none cursor-pointer hover:bg-purple-800 hover:border-purple-800 border-purple-950 py-1 px-2 border-2 rounded"
		>
			{navbarLink.name}
		</ul>
	</a>
{/each}

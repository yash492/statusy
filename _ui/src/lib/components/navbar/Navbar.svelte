<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ArrowPath, Bars3, XMark } from '@steeze-ui/heroicons';
	import { slide } from 'svelte/transition';
	import Navlinks from './Navlinks.svelte';
	import { NavbarLinks } from './NavbarItems';
	let toggleMobileNav = false;
</script>

<nav>
	<div
		class="fixed left-0 top-0 mx-auto flex h-16 w-full items-center justify-between bg-purple-950 px-4 z-50"
	>
		<h1 class="text-xl font-extrabold">Statusy</h1>
		<button class="md:hidden" on:click={() => (toggleMobileNav = !toggleMobileNav)}>
			{#if toggleMobileNav}
				<Icon src={XMark} size={'23px'} class="stroke-white"></Icon>
			{:else}
				<Icon src={Bars3} size={'23px'} class="stroke-white"></Icon>
			{/if}
		</button>

		<div class="left-0 top-0 hidden h-full items-center gap-8 font-bold md:flex">
			<div class="flex items-center bg-neutral-800 border-gray-700 rounded-sm border pr-2">
				<select
					name="refresh"
					id="refresh"
					class="pl-2 py-1 bg-neutral-800 mr-4 hover:cursor-pointer"
				>
					<option value="5">5s</option>
					<option value="10">10s</option>
					<option value="30">30s</option>
					<option value="60">60s</option>
				</select>
				<button class="hover:bg-neutral-700 px-[0.125rem] py-1">
					<Icon src={ArrowPath} size={'20px'} class="stroke-white"></Icon>
				</button>
			</div>
			<Navlinks navbarLinks={NavbarLinks} />
		</div>
		{#if toggleMobileNav}
			<div
				transition:slide={{ delay: 100, duration: 300 }}
				class="absolute left-0 top-16 flex w-full flex-col items-center justify-center gap-4 bg-purple-950 pb-4 transition-all md:hidden"
			>
				<Navlinks navbarLinks={NavbarLinks} bind:isMobileNavbarOpen={toggleMobileNav} />
			</div>
		{/if}
	</div>
</nav>

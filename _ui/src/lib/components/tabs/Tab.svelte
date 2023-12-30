<script lang="ts">
	import { getContext, onDestroy } from 'svelte';
	import type { SelectedTabContext, TabTitlesContext } from './types';
	import { SELECTED_TAB_CONTEXT, TITLE_TABS_CONTEXT } from './types';

	export let title: string;
	export let id: string;

	const tabTitles = getContext<TabTitlesContext>(TITLE_TABS_CONTEXT);
	const selectedTab = getContext<SelectedTabContext>(SELECTED_TAB_CONTEXT);

	tabTitles.registerTitle(id, title);

	onDestroy(() => {
		return tabTitles.unRegisterTitle(id);
	});
</script>

<div>
	{#if $selectedTab === id}
		<div class="px-3 pt-3">
			<slot />
		</div>
	{/if}
</div>

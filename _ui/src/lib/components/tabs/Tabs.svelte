<script lang="ts">
	import { setContext } from 'svelte';
	import {
		SELECTED_TAB_CONTEXT,
		TITLE_TABS_CONTEXT,
		type TabTitlesContext,
		type SelectedTabContext
	} from './types';
	import { writable } from 'svelte/store';

	let titles: { id: string; title: string }[] = [];

	let selectedTab = '1';
	let selectedTabStore = writable(selectedTab);

	function isTabActiveCssClasses(id: string, selectTabId: string) {
		if ($selectedTabStore === id) {
			return 'bg-purple-800';
		}
		return '';
	}

	setContext<TabTitlesContext>(TITLE_TABS_CONTEXT, {
		registerTitle(id: string, title: string) {
			titles.push({ id, title });
			titles = titles;
		},
		unRegisterTitle(id: string) {
			titles = titles.filter((title) => title.id !== id);
		}
	});

	setContext<SelectedTabContext>(SELECTED_TAB_CONTEXT, selectedTabStore);
</script>

<div class="mb-10">
	<div class="border-b border-gray-700 flex w-full text-sm md:text-base">
		{#each titles as title}
			<button
				class={`font-medium rounded-t-md w-fit h-fit px-2 py-1 mr-4 hover:bg-purple-900 ${isTabActiveCssClasses(
					title.id,
					$selectedTabStore
				)}`}
				on:click={() => ($selectedTabStore = title.id)}
			>
				{title.title}
			</button>
		{/each}
	</div>
	<slot />
</div>

import type { Writable } from 'svelte/store';

export type TabTitles = {
	id: string;
	title: string;
};

export type SelectedTabContext = Writable<string>;

export type TabTitlesContext = {
	registerTitle: (id: string, title: string) => void;
	unRegisterTitle: (id: string) => void;
};

export const TITLE_TABS_CONTEXT = 'title-tabs';
export const SELECTED_TAB_CONTEXT = 'selected-tab';

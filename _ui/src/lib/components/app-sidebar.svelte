<script lang="ts">
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { LayoutDashboard } from '@lucide/svelte';

	const defaultView = $derived(page.data.defaultView);
	const viewUrl = $derived(defaultView ? `/views/${defaultView.public_id}` : '/');
	const viewName = $derived(defaultView?.name ?? 'Default View');
	const isWorkspaceActive = $derived(
		(page.url.pathname as string).startsWith('/views/') ||
			(page.url.pathname as string).startsWith('/statuspages/')
	);
</script>

<Sidebar.Root>
	<Sidebar.Header class="border-b border-zinc-900/50 px-5 py-6">
		<a href={viewUrl} class="flex items-center gap-3 select-none">
			<div
				class="flex size-9 items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900 shadow-[0_0_12px_rgba(255,255,255,0.05)]"
			>
				<img src="/logos/statusy.svg" class="size-5" alt="Statusy Logo" />
			</div>
			<span class="text-xl font-bold tracking-tight bg-linear-to-r from-white via-zinc-200 to-zinc-400 bg-clip-text text-transparent">Statusy</span>
		</a>
	</Sidebar.Header>

	<Sidebar.Content class="px-3 py-4 flex flex-col gap-4">
		<Sidebar.Group>
			<Sidebar.GroupLabel class="px-2 text-[10px] font-bold uppercase tracking-wider text-zinc-500 select-none">Views</Sidebar.GroupLabel>
			<Sidebar.GroupContent class="mt-1.5">
				<Sidebar.Menu>
					<Sidebar.MenuItem>
						<!-- Main View button -->
						<Sidebar.MenuButton 
							isActive={isWorkspaceActive}
							class="rounded-lg hover:bg-zinc-900/50 hover:text-white transition-all duration-150 py-2 px-3 text-zinc-400"
						>
							{#snippet child({ props })}
								<a href={viewUrl} class="flex items-center gap-2.5 text-sm font-semibold" {...props}>
									<LayoutDashboard class="size-4" />
									<span>{viewName}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>

	<!-- Footer with Github Link and Sleek Layout -->
	<Sidebar.Footer class="flex items-center justify-between border-t border-zinc-900/50 px-5 py-4">
		<a
			href="https://github.com"
			target="_blank"
			rel="noopener noreferrer"
			class="flex items-center gap-2 text-xs font-semibold text-zinc-500 transition-colors hover:text-zinc-300"
		>
			<svg
				class="size-4 shrink-0"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				><path
					d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"
				/></svg
			>
			<span>GitHub</span>
		</a>
		<span class="text-[10px] font-medium text-zinc-600 select-none">v2.0.0</span>
	</Sidebar.Footer>
</Sidebar.Root>

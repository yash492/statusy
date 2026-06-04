<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { ChevronsUpDown, LayoutDashboard } from '@lucide/svelte';
	import BellRing from '@lucide/svelte/icons/bell-ring';

	// Menu items.
	const items = [
		{
			title: 'Notifications',
			url: '',
			icon: BellRing
		}
	];
</script>

<Sidebar.Root>
	<Sidebar.Content class="ml-5">
		<div class="mt-10 mb-4 flex items-center gap-2">
			<img src="/logos/statusy.svg" height="30rem" width="30rem" alt="statusy_logo" />
			<a href="/" class="text-2xl font-bold">Statusy</a>
		</div>
		<Sidebar.Group>
			<Sidebar.GroupContent class="mb-4">
				<Sidebar.Menu>
					<Sidebar.MenuItem>
						<Sidebar.MenuButton variant="default">
							{#snippet child({ props })}
								<a href={''} {...props}>
									<LayoutDashboard /> <span>{'Deepli View'}</span>
									<ChevronsUpDown />
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				</Sidebar.Menu>
			</Sidebar.GroupContent>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each items as item (item.title)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton variant="outline">
								{#snippet child({ props })}
									<a href={item.url} {...props}>
										<item.icon />
										<span>{item.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Footer class="ml-5">Github</Sidebar.Footer>
</Sidebar.Root>

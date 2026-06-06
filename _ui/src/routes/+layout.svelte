<script lang="ts">
	import { page } from '$app/state';
	import favicon from '$lib/assets/favicon.svg';
	import Navbar from '$lib/components/Navbar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';

	let { children } = $props();

	const breadcrumbs = $derived.by(() => {
		const pathname = page.url.pathname;
		const data = page.data;

		if (pathname === '/') {
			return [];
		}

		const list: Array<{ label: string; href?: string }> = [];

		if (pathname.startsWith('/views/')) {
			list.push({ label: 'Views', href: '/' });

			const viewName = data.view?.name ?? 'View';
			const viewId = page.params.publicId;
			const viewUrl = `/views/${viewId}`;

			if (pathname === viewUrl) {
				list.push({ label: viewName });
			} else {
				list.push({ label: viewName, href: viewUrl });

				if (pathname.endsWith('/add-service')) {
					list.push({ label: 'Add Service' });
				} else if (pathname.includes('/edit-service/')) {
					list.push({ label: 'Edit Service' });
				} else if (pathname.endsWith('/notifications')) {
					list.push({ label: 'Configure Notifications' });
				}
			}
		} else if (pathname.startsWith('/statuspages/')) {
			const defaultView = data.defaultView;
			const viewUrl = defaultView ? `/views/${defaultView.public_id}` : '/';

			list.push({ label: 'Views', href: viewUrl });

			const statuspageName = data.resp?.statuspage?.name ?? 'Status Page';
			list.push({ label: statuspageName });
		}

		return list;
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<Toaster />

<div class="flex min-h-screen flex-col bg-zinc-950 text-white">
	<Navbar />
	<main class="mx-auto mt-4 mb-10 w-full max-w-7xl px-4 sm:px-6 lg:px-48">
		{#if breadcrumbs.length > 0}
			<div class="mb-5">
				<Breadcrumb.Root>
					<Breadcrumb.List>
						{#each breadcrumbs as crumb, i}
							{#if i > 0}
								<Breadcrumb.Separator class="text-zinc-600" />
							{/if}
							<Breadcrumb.Item>
								{#if crumb.href}
									<Breadcrumb.Link href={crumb.href} class="text-zinc-400 hover:text-zinc-100 transition-colors">
										{crumb.label}
									</Breadcrumb.Link>
								{:else}
									<Breadcrumb.Page class="text-zinc-100">
										{crumb.label}
									</Breadcrumb.Page>
								{/if}
							</Breadcrumb.Item>
						{/each}
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>
		{/if}
		{@render children?.()}
	</main>
</div>

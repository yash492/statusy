<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { onMount } from 'svelte';

	// Predefined pool of all possible services that can be added (excluding Stripe)
	const ALL_AVAILABLE_SERVICES = [
		{
			id: 2,
			name: 'OpenAI',
			slug: 'openai',
			provider: 'Incident.io',
			status: 'degraded',
			components: ['API Core', 'API Gateway', 'ChatGPT Web', 'ChatGPT Mobile'],
			componentGroups: [
				{ name: 'API', components: ['API Core', 'API Gateway'] },
				{ name: 'ChatGPT', components: ['ChatGPT Web', 'ChatGPT Mobile'] }
			],
			lastChecked: '5 mins ago',
			lastIncident: 'Degraded performance in ChatGPT response times',
			history: [
				...Array(10).fill('operational'),
				'degraded',
				'operational',
				'operational',
				'degraded',
				'degraded'
			]
		},
		{
			id: 3,
			name: 'Plivo',
			slug: 'plivo',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['SMS API', 'SMS Gateway', 'Voice API', 'Voice Portal'],
			componentGroups: [
				{ name: 'SMS', components: ['SMS API', 'SMS Gateway'] },
				{ name: 'Voice', components: ['Voice API', 'Voice Portal'] }
			],
			lastChecked: '1 min ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 4,
			name: 'Twilio',
			slug: 'twilio',
			provider: 'Atlassian Statuspage',
			status: 'major_outage',
			components: ['WhatsApp Messaging', 'Console UI', 'Console API'],
			componentGroups: [
				{ name: 'Messaging', components: ['WhatsApp Messaging'] },
				{ name: 'Console', components: ['Console UI', 'Console API'] }
			],
			lastChecked: '30s ago',
			lastIncident: 'Major Outage on SMS delivery in EU regions',
			history: [...Array(12).fill('operational'), 'operational', 'operational', 'major_outage']
		},
		{
			id: 5,
			name: 'GitHub',
			slug: 'github',
			provider: 'Atlassian Statuspage',
			status: 'degraded',
			components: ['Git Operations', 'GitHub Actions', 'Codespaces'],
			componentGroups: [
				{ name: 'Core', components: ['Git Operations'] },
				{ name: 'Developer Tools', components: ['GitHub Actions', 'Codespaces'] }
			],
			lastChecked: '1 min ago',
			lastIncident: 'Degraded performance on GitHub Actions runs',
			history: [
				...Array(11).fill('operational'),
				'degraded',
				'operational',
				'operational',
				'degraded'
			]
		},
		{
			id: 6,
			name: 'Cloudflare',
			slug: 'cloudflare',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['DNS Resolvers', 'CDN Edge', 'Workers API', 'Workers KV'],
			componentGroups: [
				{ name: 'Network', components: ['DNS Resolvers', 'CDN Edge'] },
				{ name: 'Compute', components: ['Workers API', 'Workers KV'] }
			],
			lastChecked: '2 mins ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 7,
			name: 'DigitalOcean',
			slug: 'digitalocean',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['Droplets', 'App Platform', 'Volumes'],
			componentGroups: [
				{ name: 'Compute', components: ['Droplets', 'App Platform'] },
				{ name: 'Storage', components: ['Volumes'] }
			],
			lastChecked: '4 mins ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 8,
			name: 'Datadog',
			slug: 'datadog',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['APM Tracing', 'Log Ingestion', 'Synthetics Tests'],
			componentGroups: [
				{ name: 'Observability', components: ['APM Tracing', 'Log Ingestion'] },
				{ name: 'Testing', components: ['Synthetics Tests'] }
			],
			lastChecked: '1 min ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 9,
			name: 'SolarWinds Observability',
			slug: 'solarwindsobservability',
			provider: 'Atlassian Statuspage',
			status: 'major_outage',
			components: ['Alerting', 'Ingest Pipeline'],
			componentGroups: [
				{ name: 'Alerting', components: ['Alerting'] },
				{ name: 'Pipeline', components: ['Ingest Pipeline'] }
			],
			lastChecked: '1 min ago',
			lastIncident: 'Ingest latency spike and processing delays',
			history: [...Array(13).fill('operational'), 'degraded', 'major_outage']
		},
		{
			id: 10,
			name: 'Cursor',
			slug: 'cursor',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['Copilot API', 'Backend Sync'],
			componentGroups: [
				{ name: 'Copilot', components: ['Copilot API'] },
				{ name: 'Sync', components: ['Backend Sync'] }
			],
			lastChecked: '3 mins ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 11,
			name: 'CircleCI',
			slug: 'circleci',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['Build Runners', 'API'],
			componentGroups: [
				{ name: 'Build', components: ['Build Runners'] },
				{ name: 'API', components: ['API'] }
			],
			lastChecked: '5 mins ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 12,
			name: 'Claude',
			slug: 'claude',
			provider: 'Atlassian Statuspage',
			status: 'degraded',
			components: ['Claude API', 'Claude.ai Web App'],
			componentGroups: [
				{ name: 'API', components: ['Claude API'] },
				{ name: 'Chat', components: ['Claude.ai Web App'] }
			],
			lastChecked: '2 mins ago',
			lastIncident: 'Claude API rate limits returning elevated 5xx errors',
			history: [...Array(14).fill('operational'), 'degraded']
		},
		{
			id: 13,
			name: 'New Relic',
			slug: 'newrelic',
			provider: 'Atlassian Statuspage',
			status: 'operational',
			components: ['Metrics Ingestion', 'User Interface'],
			componentGroups: [
				{ name: 'Ingestion', components: ['Metrics Ingestion'] },
				{ name: 'UI', components: ['User Interface'] }
			],
			lastChecked: '4 mins ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		},
		{
			id: 14,
			name: 'HashiCorp',
			slug: 'hashicorp',
			provider: 'Incident.io',
			status: 'operational',
			components: ['HCP Portal', 'Vault Secrets'],
			componentGroups: [
				{ name: 'Platform', components: ['HCP Portal'] },
				{ name: 'Secrets', components: ['Vault Secrets'] }
			],
			lastChecked: '3 mins ago',
			lastIncident: null,
			history: Array(15).fill('operational')
		}
	];

	function getDefaultServicesForSlug(vSlug: string): any[] {
		if (vSlug === 'payment-gateways' || vSlug === 'payments') {
			return ALL_AVAILABLE_SERVICES.filter((s) => [2, 3, 4].includes(s.id));
		} else if (vSlug === 'core-infrastructure' || vSlug === 'infra') {
			return ALL_AVAILABLE_SERVICES.filter((s) => [5, 6, 7, 8, 9].includes(s.id));
		} else if (vSlug === 'developer-tools' || vSlug === 'dev-tools') {
			return ALL_AVAILABLE_SERVICES.filter((s) => [10, 11, 12, 13, 14].includes(s.id));
		}
		return [];
	}

	const slug = $derived(page.params.slug);
	const serviceSlug = $derived(page.params.serviceSlug);

	let localServices = $state<any[]>([]);
	let activeService = $state<any | null>(null);
	let selectedComponents = $state<string[]>([]);
	let componentMode = $state<'all' | 'custom'>('all');

	// Original service definition to extract all possible components
	let originalService = $state<any | null>(null);

	onMount(() => {
		const stored = localStorage.getItem(`statusy_view_${slug}`);
		if (stored) {
			localServices = JSON.parse(stored);
		} else {
			localServices = getDefaultServicesForSlug(slug || '');
		}

		activeService = localServices.find((s) => s.slug === serviceSlug) || null;
		if (activeService) {
			originalService = ALL_AVAILABLE_SERVICES.find((s) => s.slug === serviceSlug) || activeService;
			selectedComponents = [...activeService.components];

			const hasAll = originalService.components.every((c: string) =>
				selectedComponents.includes(c)
			);
			if (hasAll && selectedComponents.length === originalService.components.length) {
				componentMode = 'all';
			} else {
				componentMode = 'custom';
			}
		}
	});

	function toggleComponent(component: string) {
		if (selectedComponents.includes(component)) {
			selectedComponents = selectedComponents.filter((c) => c !== component);
		} else {
			selectedComponents = [...selectedComponents, component];
		}
	}

	function toggleGroup(group: { name: string; components: string[] }) {
		const allChecked = group.components.every((c: string) => selectedComponents.includes(c));
		if (allChecked) {
			selectedComponents = selectedComponents.filter((c: string) => !group.components.includes(c));
		} else {
			const toAdd = group.components.filter((c: string) => !selectedComponents.includes(c));
			selectedComponents = [...selectedComponents, ...toAdd];
		}
	}

	function selectAll(currentService: any) {
		selectedComponents = [...currentService.components];
	}

	function deselectAll() {
		selectedComponents = [];
	}

	function saveService() {
		if (activeService) {
			localServices = localServices.map((s) => {
				if (s.slug === serviceSlug) {
					return {
						...s,
						components: [...selectedComponents]
					};
				}
				return s;
			});
			localStorage.setItem(`statusy_view_${slug}`, JSON.stringify(localServices));
			void goto(`/views/${slug}`);
		}
	}

	function handleCancel() {
		void goto(`/views/${slug}`);
	}
</script>

<svelte:head>
	<title>Edit Service | Statusy</title>
</svelte:head>

<div class="mx-auto w-3/5 pt-4 text-white">
	<!-- Header -->
	<div class="mb-6 max-w-xl">
		<h1 class="text-3xl font-extrabold tracking-tight text-white sm:text-4xl">Edit Service</h1>
		<p class="mt-2 text-sm text-zinc-400">
			Customize the component subscriptions for this service dashboard view.
		</p>
	</div>

	<!-- Form Card (Borderless & Background Merged) -->
	<div class="max-w-xl">
		{#if activeService}
			<div class="grid gap-6">
				<!-- Service Name (Read Only) -->
				<div class="grid gap-2">
					<Label class="text-sm font-semibold text-zinc-300">Service Name</Label>
					<div class="rounded-lg bg-zinc-900/50 px-4 py-3 text-sm font-medium text-zinc-200">
						{activeService.name}
					</div>
				</div>

				<!-- Component Checklist with Hierarchical groups -->
				{#if originalService}
					<div class="grid gap-4 pt-5">
						<Label class="text-sm font-semibold text-zinc-300">Monitored Components</Label>

						<!-- Custom Radio Buttons and Component checklist in border container -->
						<div
							class="flex flex-col gap-4 rounded-lg border border-zinc-800/40 bg-zinc-900/20 p-4"
						>
							<label class="group flex cursor-pointer items-start gap-3">
								<input
									type="radio"
									name="component-mode"
									value="all"
									checked={componentMode === 'all'}
									onchange={() => {
										componentMode = 'all';
										selectAll(originalService);
									}}
									class="sr-only"
								/>
								<div
									class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full border border-zinc-700 transition-all group-hover:border-zinc-500 {componentMode ===
									'all'
										? 'border-emerald-500 bg-emerald-500/10'
										: ''}"
								>
									{#if componentMode === 'all'}
										<div class="size-2 rounded-full bg-emerald-500"></div>
									{/if}
								</div>
								<div class="flex flex-col">
									<span
										class="text-sm font-medium text-zinc-200 transition-colors group-hover:text-white"
										>Monitor all components</span
									>
									<span class="mt-0.5 text-xs text-zinc-500"
										>Automatically monitor all current and future components for this service</span
									>
								</div>
							</label>

							<label class="group flex cursor-pointer items-start gap-3">
								<input
									type="radio"
									name="component-mode"
									value="custom"
									checked={componentMode === 'custom'}
									onchange={() => {
										componentMode = 'custom';
									}}
									class="sr-only"
								/>
								<div
									class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full border border-zinc-700 transition-all group-hover:border-zinc-500 {componentMode ===
									'custom'
										? 'border-emerald-500 bg-emerald-500/10'
										: ''}"
								>
									{#if componentMode === 'custom'}
										<div class="size-2 rounded-full bg-emerald-500"></div>
									{/if}
								</div>
								<div class="flex flex-col">
									<span class="font-medium text-zinc-200 transition-colors group-hover:text-white"
										>Customize component selection</span
									>
									<span class="mt-0.5 text-zinc-500"
										>Select specific component groups or individual components to monitor</span
									>
								</div>
							</label>

							{#if componentMode === 'custom'}
								<div class="ml-7 grid max-h-72 grid-cols-1 gap-3 overflow-y-auto pr-1">
									{#each originalService.componentGroups as group}
										{@const groupChecked = group.components.every((c: string) =>
											selectedComponents.includes(c)
										)}
										{@const groupSomeChecked =
											group.components.some((c: string) => selectedComponents.includes(c)) &&
											!groupChecked}

										<div class="py-2">
											<!-- Group Header -->
											<button
												type="button"
												onclick={() => toggleGroup(group)}
												class="flex w-full cursor-pointer items-center gap-3 text-left transition-all hover:text-white"
											>
												<div
													class="flex size-4 shrink-0 items-center justify-center rounded transition-all {groupChecked
														? 'bg-emerald-500 text-zinc-950'
														: groupSomeChecked
															? 'bg-emerald-500/20 text-emerald-400'
															: 'bg-zinc-900'}"
												>
													{#if groupChecked}
														<svg
															xmlns="http://www.w3.org/2000/svg"
															fill="none"
															viewBox="0 0 24 24"
															stroke-width="3"
															stroke="currentColor"
															class="size-2.5"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																d="m4.5 12.75 6 6 9-13.5"
															/>
														</svg>
													{:else if groupSomeChecked}
														<div class="size-1.5 rounded-sm bg-emerald-400"></div>
													{/if}
												</div>
												<span class="text-sm font-bold text-zinc-200">{group.name} Group</span>
											</button>

											<!-- Child Components -->
											<div class="mt-2 ml-7 grid gap-2 pl-3">
												{#each group.components as component}
													{@const componentChecked = selectedComponents.includes(component)}
													<button
														type="button"
														onclick={() => toggleComponent(component)}
														class="flex cursor-pointer items-center gap-3 text-left transition-all hover:text-zinc-200 {componentChecked
															? 'text-zinc-300'
															: 'text-zinc-500'}"
													>
														<div
															class="flex size-3.5 shrink-0 items-center justify-center rounded transition-all {componentChecked
																? 'bg-emerald-500/80 text-zinc-950'
																: 'bg-zinc-900'}"
														>
															{#if componentChecked}
																<svg
																	xmlns="http://www.w3.org/2000/svg"
																	fill="none"
																	viewBox="0 0 24 24"
																	stroke-width="3.5"
																	stroke="currentColor"
																	class="size-2"
																>
																	<path
																		stroke-linecap="round"
																		stroke-linejoin="round"
																		d="m4.5 12.75 6 6 9-13.5"
																	/>
																</svg>
															{/if}
														</div>
														<span class="text-xs font-medium">{component}</span>
													</button>
												{/each}
											</div>
										</div>
									{/each}
								</div>
							{/if}
						</div>
					</div>
				{/if}
			</div>

			<!-- Footer Buttons -->
			<div class="mt-8 flex items-center justify-end gap-3 pt-5">
				<Button
					variant="ghost"
					class="cursor-pointer px-5 hover:bg-zinc-900 hover:text-white"
					onclick={handleCancel}
				>
					Cancel
				</Button>
				<Button
					class="cursor-pointer px-6"
					disabled={selectedComponents.length === 0}
					onclick={saveService}
				>
					Save Changes
				</Button>
			</div>
		{:else}
			<div class="py-6 text-center text-zinc-500">Loading service details...</div>
		{/if}
	</div>
</div>

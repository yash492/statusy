<script lang="ts">
	import {
		NotificationsApi,
		type ViewNotification,
		type ViewNotificationRequest
	} from '$lib/api/notifications/notifications';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import Bell from '@lucide/svelte/icons/bell';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Search from '@lucide/svelte/icons/search';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const notificationsApi = new NotificationsApi();

	let notificationConfigs = $state<ViewNotification[]>(data.notifications);

	// Search State
	let searchQuery = $state('');

	const filteredConfigs = $derived(
		notificationConfigs.filter(
			(config) =>
				config.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				config.type.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	// Pagination State
	const itemsPerPage = 5;
	let currentPage = $state(1);

	$effect(() => {
		const _ = searchQuery;
		currentPage = 1;
	});

	const paginatedConfigs = $derived(
		filteredConfigs.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
	);
	const totalConfigsCount = $derived(filteredConfigs.length);
	const totalPages = $derived(Math.ceil(totalConfigsCount / itemsPerPage) || 1);

	// Modal and Form State
	let isFormModalOpen = $state(false);
	let editingConfigId = $state<number | null>(null);
	let newConfigType = $state<ViewNotificationRequest['type']>('slack');
	let newConfigName = $state('');
	let newConfigValue = $state('');
	let newConfigHeaders = $state<Array<{ id: string; key: string; value: string }>>([]);
	let isSaving = $state(false);

	function addHeaderField() {
		newConfigHeaders = [
			...newConfigHeaders,
			{ id: `hdr-${Date.now()}-${Math.random()}`, key: '', value: '' }
		];
	}

	function deleteHeaderField(id: string) {
		newConfigHeaders = newConfigHeaders.filter((h) => h.id !== id);
	}

	function startAddNew() {
		editingConfigId = null;
		newConfigType = 'slack';
		newConfigName = '';
		newConfigValue = '';
		newConfigHeaders = [];
		isFormModalOpen = true;
	}

	function startEdit(conf: ViewNotification) {
		editingConfigId = conf.id;
		newConfigType = conf.type as ViewNotificationRequest['type'];
		newConfigName = conf.name;

		const config = conf.config as Record<string, any>;

		if (config.headers) {
			newConfigHeaders = Object.entries(config.headers).map(([key, value]) => ({
				id: `hdr-${Date.now()}-${Math.random()}`,
				key,
				value: String(value)
			}));
		} else {
			newConfigHeaders = [];
		}

		if (conf.type === 'pagerduty') {
			newConfigValue = config.routing_key || '';
		} else if (conf.type === 'webhook') {
			newConfigValue = config.url || '';
		} else {
			newConfigValue = config.webhook_url || '';
		}
		isFormModalOpen = true;
	}

	function buildConfigFields(): Record<string, any> {
		const configFields: Record<string, any> = {};
		if (newConfigType === 'pagerduty') {
			configFields.routing_key = newConfigValue;
		} else if (newConfigType === 'webhook') {
			configFields.url = newConfigValue;

			const headersObj: Record<string, string> = {};
			newConfigHeaders.forEach((h) => {
				if (h.key.trim()) {
					headersObj[h.key.trim()] = h.value;
				}
			});
			configFields.headers = headersObj;
		} else {
			configFields.webhook_url = newConfigValue;
		}
		return configFields;
	}

	async function saveConfig() {
		if (!newConfigName.trim()) {
			toast.error('Display Name is required');
			return;
		}
		if (!newConfigValue.trim()) {
			toast.error('Configuration value is required');
			return;
		}

		isSaving = true;
		const body: ViewNotificationRequest = {
			name: newConfigName.trim(),
			type: newConfigType,
			config: buildConfigFields()
		};

		if (editingConfigId) {
			const [result, err] = await notificationsApi.edit(data.view.public_id, editingConfigId, body);
			if (err) {
				toast.error(err.message || 'Failed to update notification');
				isSaving = false;
				return;
			}
			notificationConfigs = notificationConfigs.map((c) =>
				c.id === editingConfigId ? result! : c
			);
			toast.success('Notification updated');
		} else {
			const [result, err] = await notificationsApi.create(data.view.public_id, body);
			if (err) {
				toast.error(err.message || 'Failed to add notification');
				isSaving = false;
				return;
			}
			notificationConfigs = [...notificationConfigs, result!];
			toast.success('Notification added');
		}

		isSaving = false;
		isFormModalOpen = false;
		editingConfigId = null;
	}

	async function deleteConfig(id: number) {
		const [, err] = await notificationsApi.delete(data.view.public_id, id);
		if (err) {
			toast.error(err.message || 'Failed to delete notification');
			return;
		}
		notificationConfigs = notificationConfigs.filter((c) => c.id !== id);
		toast.success('Notification removed');
		if (editingConfigId === id) {
			isFormModalOpen = false;
			editingConfigId = null;
		}
		// Adjust page number if the current page becomes empty
		if (currentPage > 1 && paginatedConfigs.length === 0) {
			currentPage -= 1;
		}
	}
</script>

<svelte:head>
	<title>Configure Notifications | {data.view.name}</title>
</svelte:head>

<div class="mx-auto max-w-4xl px-4 py-6">
	<!-- Navigation and Header -->
	<div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex flex-col gap-1">
			<h1
				class="flex items-center gap-2.5 text-2xl font-extrabold tracking-tight text-white sm:text-3xl"
			>
				<Bell class="size-7 text-indigo-400" />
				Configure Notifications
			</h1>
			<p class="max-w-xl text-sm text-zinc-400">
				Manage notification settings for <strong>{data.view.name}</strong>. Alerts will only be
				dispatched for the specific services and components monitored by this view.
			</p>
		</div>
	</div>

	<!-- Subscribed Notifications Section -->
	<div class="space-y-4">
		<div class="border-b border-zinc-900 pb-3">
			<h2 class="text-sm font-semibold tracking-wider text-zinc-500 uppercase">
				Subscribed Notifications
			</h2>
		</div>

		<!-- Search & Add Bar -->
		<div class="mb-4 flex items-center justify-between gap-3">
			<div class="relative w-full max-w-sm">
				<input
					type="text"
					placeholder="Search notifications..."
					bind:value={searchQuery}
					class="w-full rounded-lg border border-zinc-800 bg-zinc-900/40 py-2 pr-4 pl-9 text-sm text-white placeholder-zinc-500 outline-none focus:border-zinc-700 focus:ring-1 focus:ring-zinc-700"
				/>
				<div class="absolute inset-y-0 left-0 flex items-center pl-3 text-zinc-500">
					<Search class="size-4" />
				</div>
			</div>

			<Button
				variant="outline"
				size="sm"
				onclick={startAddNew}
				class="h-9 shrink-0 cursor-pointer border-zinc-800 bg-zinc-900/50 px-4 font-semibold text-white hover:bg-zinc-800"
			>
				<Plus class="mr-1.5 size-4" />
				Add Notification
			</Button>
		</div>

		<!-- List of configurations -->
		<div class="grid grid-cols-1 gap-3">
			{#each paginatedConfigs as config (config.id)}
				<div
					class="flex items-center justify-between rounded-xl border border-zinc-800 bg-zinc-900/10 p-4 transition-all hover:bg-zinc-900/30"
				>
					<div class="flex min-w-0 flex-col gap-1 pr-4">
						<div class="flex items-center gap-2.5">
							<span class="truncate text-base font-bold text-white">{config.name}</span>
							<span
								class="rounded-full border border-zinc-700/30 bg-zinc-800/80 px-2 py-0.5 text-[10px] font-bold tracking-wider text-zinc-400 uppercase"
							>
								{config.type}
							</span>
						</div>
					</div>

					<div class="flex shrink-0 items-center gap-2">
						<button
							onclick={() => startEdit(config)}
							class="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-zinc-800 bg-zinc-900/50 text-zinc-400 transition-all hover:border-zinc-700 hover:text-white"
							title="Edit Notification"
						>
							<Pencil class="size-3.5" />
						</button>
						<button
							onclick={() => deleteConfig(config.id)}
							class="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-red-500/10 bg-red-950/10 text-red-400 transition-all hover:bg-red-900/20 hover:text-red-300"
							title="Delete Notification"
						>
							<Trash2 class="size-3.5" />
						</button>
					</div>
				</div>
			{:else}
				<div
					class="flex flex-col items-center justify-center rounded-xl border border-dashed border-zinc-800 py-16 text-center text-zinc-500"
				>
					<Bell class="size-8 text-zinc-600 mb-3" />
					<p class="text-sm font-medium">No notification configurations found.</p>
				</div>
			{/each}
		</div>

		<!-- Pagination Footer -->
		{#if totalConfigsCount > 0}
			<div
				class="mt-6 flex items-center justify-between border-t border-zinc-900 pt-4 text-sm text-zinc-400"
			>
				<div>
					Showing <span class="font-medium text-white">{(currentPage - 1) * itemsPerPage + 1}</span>
					to
					<span class="font-medium text-white"
						>{Math.min(currentPage * itemsPerPage, totalConfigsCount)}</span
					>
					of <span class="font-medium text-white">{totalConfigsCount}</span> notifications
				</div>
				<div class="flex gap-2">
					<Button
						variant="outline"
						size="sm"
						class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
						disabled={currentPage === 1 || totalPages <= 1}
						onclick={() => (currentPage -= 1)}
					>
						Previous
					</Button>
					<Button
						variant="outline"
						size="sm"
						class="cursor-pointer border-zinc-800 hover:bg-zinc-900 hover:text-white"
						disabled={currentPage === totalPages || totalPages <= 1}
						onclick={() => (currentPage += 1)}
					>
						Next
					</Button>
				</div>
			</div>
		{/if}
	</div>
</div>

<!-- Add/Edit Notification Modal -->
<Dialog.Root
	open={isFormModalOpen}
	onOpenChange={(open) => {
		if (!open) {
			isFormModalOpen = false;
			editingConfigId = null;
		}
	}}
>
	<Dialog.Content class="border-zinc-800 bg-zinc-950 text-white shadow-xl sm:max-w-[450px]">
		<Dialog.Header>
			<Dialog.Title class="text-lg font-bold text-white">
				{editingConfigId ? 'Edit Notification' : 'Add Notification'}
			</Dialog.Title>
			<Dialog.Description class="text-zinc-400">
				Configure the settings for your notification destination.
			</Dialog.Description>
		</Dialog.Header>

		<div class="mt-4 space-y-4">
			<div class="flex flex-col gap-1.5">
				<label for="channel-type" class="text-xs font-semibold text-zinc-400"
					>Notification Type</label
				>
				<select
					id="channel-type"
					bind:value={newConfigType}
					class="cursor-pointer rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-sm text-white outline-none focus:border-zinc-700"
				>
					<option value="slack">Slack Webhook</option>
					<option value="discord">Discord Webhook</option>
					<option value="msteams">MS Teams Webhook</option>
					<option value="pagerduty">PagerDuty Event v2</option>
					<option value="solarwinds_incident_response">Solarwinds Incident Response Webhook</option>
					<option value="webhook">Custom Webhook</option>
				</select>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="channel-name" class="text-xs font-semibold text-zinc-400">Display Name</label>
				<input
					id="channel-name"
					type="text"
					placeholder="e.g. Engineering Slack Alerts"
					bind:value={newConfigName}
					class="placeholder:text-zinc-650 rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-sm text-white outline-none focus:border-zinc-700"
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="channel-value" class="text-xs font-semibold text-zinc-400">
					{#if newConfigType === 'pagerduty'}
						Routing/Integration Key
					{:else if newConfigType === 'webhook'}
						Target URL
					{:else}
						Webhook URL
					{/if}
				</label>
				<div class="relative flex items-center">
					{#if newConfigType === 'webhook'}
						<button
							disabled
							type="button"
							class="absolute left-2.5 cursor-not-allowed rounded border border-zinc-700/50 bg-zinc-800/80 px-2 py-0.5 text-[10px] font-bold text-zinc-400 select-none"
						>
							POST
						</button>
					{/if}
					<input
						id="channel-value"
						type="text"
						placeholder={newConfigType === 'pagerduty' ? 'pd-service-key-xxx' : 'https://...'}
						bind:value={newConfigValue}
						class="placeholder:text-zinc-650 w-full rounded-lg border border-zinc-800 bg-zinc-900 py-2 text-sm text-white outline-none focus:border-zinc-700 {newConfigType ===
						'webhook'
							? 'pr-3 pl-16'
							: 'px-3'}"
					/>
				</div>
			</div>

			{#if newConfigType === 'webhook'}
				<!-- Custom Headers -->
				<div class="flex flex-col gap-2">
					<div class="mt-1 flex items-center justify-between border-t border-zinc-900 pt-3">
						<span class="text-xs font-semibold text-zinc-400">Custom HTTP Headers</span>
						<button
							type="button"
							onclick={addHeaderField}
							class="flex cursor-pointer items-center gap-1 text-[11px] font-semibold text-indigo-400 transition-colors hover:text-indigo-300"
						>
							<Plus class="size-3" />
							Add Header
						</button>
					</div>

					{#if newConfigHeaders.length === 0}
						<p class="text-[11px] text-zinc-600 italic">No custom headers added yet.</p>
					{:else}
						<div class="max-h-36 space-y-2 overflow-y-auto pr-1">
							{#each newConfigHeaders as header (header.id)}
								<div class="flex items-center gap-2">
									<input
										type="text"
										placeholder="Header-Name"
										bind:value={header.key}
										class="w-1/2 rounded-lg border border-zinc-800 bg-zinc-900 px-2.5 py-1.5 font-mono text-xs text-white outline-none placeholder:text-zinc-600 focus:border-zinc-700"
									/>
									<input
										type="text"
										placeholder="value"
										bind:value={header.value}
										class="w-1/2 rounded-lg border border-zinc-800 bg-zinc-900 px-2.5 py-1.5 font-mono text-xs text-white outline-none placeholder:text-zinc-600 focus:border-zinc-700"
									/>
									<button
										type="button"
										onclick={() => deleteHeaderField(header.id)}
										class="cursor-pointer p-1 text-zinc-500 transition-colors hover:text-red-400"
										title="Remove Header"
									>
										<Trash2 class="size-3.5" />
									</button>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<div class="mt-6 flex justify-end gap-2.5 border-t border-zinc-900 pt-4">
			<Button
				variant="outline"
				class="cursor-pointer border-zinc-800 text-zinc-300 hover:bg-zinc-900"
				onclick={() => {
					isFormModalOpen = false;
					editingConfigId = null;
				}}
			>
				Cancel
			</Button>
			<Button
				class="cursor-pointer bg-indigo-600 text-white hover:bg-indigo-500"
				onclick={saveConfig}
				disabled={isSaving}
			>
				{isSaving ? 'Saving...' : 'Save Notification'}
			</Button>
		</div>
	</Dialog.Content>
</Dialog.Root>

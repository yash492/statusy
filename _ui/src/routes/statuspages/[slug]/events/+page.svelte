<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import IncidentsTable from '$lib/derivedcomponents/IncidentsTable.svelte';

	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<div class="mx-auto w-4/5">
	<div class="w-full">
		<div class="mb-6 flex justify-between md:mb-4">
			<div class="mb-4 flex w-fit items-center gap-2">
				<div class="rounded-4xl border-2 bg-white">
					<img src="/provider_logo/plivo.png" width="30px" height="30px" alt="plivo-logo" />
				</div>
				<p class="text-xl font-bold">{data.resp.statuspage.name}</p>
			</div>
			<div>
				<Button class="cursor-pointer">Subscribe to Updates</Button>
			</div>
		</div>

		<div>
			<div>
				<Tabs.Root value="incidents" class="w-full">
					<Tabs.List>
						<Tabs.Trigger value="incidents">Incidents</Tabs.Trigger>
						<Tabs.Trigger value="scheduled-maintenance">Scheduled Maintenances</Tabs.Trigger>
					</Tabs.List>
					<Tabs.Content value="incidents">
						<IncidentsTable />
					</Tabs.Content>
					<Tabs.Content value="scheduled-maintenance">
						<IncidentsTable />
					</Tabs.Content>
				</Tabs.Root>
			</div>
		</div>
	</div>
</div>

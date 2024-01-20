<script lang="ts">
	export let normalisedStatus: string;
	export let status: string;

	let bgColor = '';

	$: switch (normalisedStatus) {
		case 'triggered':
			bgColor = `bg-purple-800`;
			break;
		case 'in_progress':
			bgColor = 'bg-yellow-800';
			break;
		case 'resolved':
			bgColor = 'bg-green-800';
			break;
	}
</script>

<div class={`w-fit px-3 text-center font-semibold ${bgColor} rounded-lg`}>
	<p>{status}</p>
</div>

<script lang="ts">
	export let name: string;
	export let link: string;
</script>

<p>
	<a class="hover:text-blue-300" target="_blank" href={link}>{name}</a>
</p>

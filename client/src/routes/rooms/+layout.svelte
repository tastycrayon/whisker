<script lang="ts">
	import { refreshRooms, roomStore } from '$lib/store';
	import { onMount } from 'svelte';

	onMount(() => {
		if ($roomStore.data.length === 0) refreshRooms();
	});
</script>

<slot />

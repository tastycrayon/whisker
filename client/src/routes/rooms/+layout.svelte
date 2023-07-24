<script lang="ts">
	import { goto } from '$app/navigation';
	import { LOGIN_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
	import { refreshRooms, roomStore } from '$lib/store';
	import { onMount } from 'svelte';

	$: if (!$currentUser) goto(LOGIN_PATH);

	onMount(() => {
		if ($roomStore.data.length === 0) refreshRooms();
	});
</script>

<slot />

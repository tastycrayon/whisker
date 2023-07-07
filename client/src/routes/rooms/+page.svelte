<script lang="ts">
	import { goto } from '$app/navigation';
	import Icon from '$components/icon.svelte';
	import Rooms from '$components/room/rooms.svelte';
	import { DEFAULT_ROOM, ROOM_PATH } from '$lib/constant';
	import { currentRoom, refreshRooms, roomStore } from '$lib/store';
	import { modalStore, type ModalComponent, type ModalSettings } from '@skeletonlabs/skeleton';
	import { onDestroy, onMount } from 'svelte';

	const modal: ModalSettings = {
		type: 'confirm',
		title: `Go to ${$currentRoom}`,
		body: 'Are you sure you wish to proceed?',
		modalClasses: 'card max-w-[400px]',
		// TRUE if confirm pressed, FALSE if cancel pressed
		response: (r: boolean) => {
			if (r) goto(ROOM_PATH + '/' + $currentRoom);
		}
	};
	let loaded = false;
	// get current room information
	$: console.log({ ff: $roomStore });
	$: if (!$roomStore.loading && $roomStore.data.length != 0 && !loaded) {
		loaded = true;
		modal.title = $roomStore.data.find((e) => e.slug == $currentRoom)?.name || DEFAULT_ROOM;
		modalStore.trigger(modal);
	}
	onMount(() => {
		if ($roomStore.data.length === 0) refreshRooms();
	});
	onDestroy(() => {
		modalStore.close();
	});
</script>

<div class="card p-4 max-w-[640px] mx-auto space-y-4 my-4">
	<div class="inline-flex items-center">
		<h1 class="h1">Rooms</h1>
		<Icon name="cat" class="text-purple-600" width="48px" height="48px" />
	</div>
	<Rooms />
</div>

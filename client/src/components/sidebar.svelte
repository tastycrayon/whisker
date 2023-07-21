<script lang="ts">
	import { Autocomplete, drawerStore, type AutocompleteOption } from '@skeletonlabs/skeleton';
	import Rooms from './room/rooms.svelte';
	import { refreshRooms, roomStore } from '$lib/store';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { ROOM_PATH } from '$lib/constant';

	let inputDemo = '';
	$: roomOptions = $roomStore.data.map((e) => ({
		label: e.name,
		value: e.slug,
		keywords: [e.description, e.createdBy, e.type].join(', '),
		meta: { slug: e.slug }
	}));

	let isFocused = false;
	function onRoomSelection(event: any): void {
		inputDemo = event.detail.label;
		goto(ROOM_PATH + '/' + event.detail.meta.slug);
		inputDemo = '';
		isFocused = false;
		drawerStore.close();
	}

	const onFocus = () => (isFocused = true);
	onMount(() => {
		if ($roomStore.data.length === 0) refreshRooms();
	});

	const handleDrawyerClose = () => {
		drawerStore.close();
	};
</script>

<!-- Header -->
<header class="border-b border-surface-500/30 p-4 relative">
	<input
		bind:value={inputDemo}
		on:focus={onFocus}
		class="input"
		type="search"
		placeholder="Search..."
	/>
	{#if isFocused}
		<div class="absolute top-100 left-0 w-full p-4 z-[10]">
			<div class="card w-full max-w-sm max-h-48 p-4 overflow-y-auto" tabindex="-1">
				<Autocomplete
					bind:input={inputDemo}
					options={roomOptions}
					on:selection={onRoomSelection}
					on:click={handleDrawyerClose}
					class="text-sm"
				/>
			</div>
		</div>
	{/if}
</header>
<Rooms />

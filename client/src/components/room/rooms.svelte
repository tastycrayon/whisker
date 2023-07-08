<script lang="ts">
	import {
		Avatar,
		modalStore,
		type ModalComponent,
		type ModalSettings
	} from '@skeletonlabs/skeleton';
	import CreateRoom from './create-room.svelte';
	import { onMount } from 'svelte';
	import { refreshRooms, roomStore } from '$lib/store';
	import Icon from '../icon.svelte';
	import { RoomType } from '$lib/types';
	import RoomItem from './roomlist-item.svelte';
	import { PERSONAL_ROOM_PLACEHOLDER_COUNT, PUBLIC_ROOM_PLACEHOLDER_COUNT } from '$lib/constant';

	// modal
	const modalComponent: ModalComponent = {
		// Pass a reference to your custom component
		ref: CreateRoom,
		// Add the component properties as key/value pairs
		props: { background: 'bg-red-500' },
		// Provide a template literal for the default component slot
		slot: '<p>Skeleton</p>'
	};
	const handleModal = () => {
		const modal: ModalSettings = {
			type: 'component',
			// Pass the component directly:
			component: modalComponent
		};
		modalStore.trigger(modal);
	};
	let publicOpen = true;
	let personalOpen = true;
</script>

<!-- List -->
<div class="p-4 space-y-4 overflow-y-auto">
	<nav class="list-nav">
		<div class="w-full inline-flex justify-between items-center pb-2">
			<button class="btn btn-sm variant-ghost-surface" on:click={handleModal}>
				<span><Icon name="chat" width="16px" height="16px" /></span>
				<span>Create Room</span>
			</button>
			<button
				class="btn-icon btn-icon-sm variant-soft"
				disabled={$roomStore.loading}
				on:click={refreshRooms}
			>
				<span><Icon name="refresh" width="16px" height="16px" /></span>
			</button>
		</div>
		<hr />
		<div class="w-full inline-flex justify-between items-center py-1">
			<small class="opacity-50">DEFAULT</small>
			<button
				class="btn-icon btn-icon-sm variant-soft"
				on:click={() => {
					publicOpen = !publicOpen;
				}}
			>
				{#if publicOpen}
					<span><Icon name="minus" width="8px" height="8px" /></span>
				{:else}
					<span><Icon name="plus" width="8px" height="8px" /></span>
				{/if}
			</button>
		</div>
		<hr />
		<ul class={`list ${publicOpen ? '' : 'hidden'}`}>
			{#if $roomStore.loading}
				<div class="space-y-5 py-2">
					{#each new Array(PUBLIC_ROOM_PLACEHOLDER_COUNT).fill(0) as _}
						<div class="animate-pulse placeholder h-10" />
					{/each}
				</div>
			{:else if $roomStore.data}
				{#each $roomStore.data.filter((e) => e.type == RoomType.PublicRoom) as room}
					<RoomItem {room} />
				{/each}
			{/if}
		</ul>
		<hr />
		<div class="w-full inline-flex justify-between items-center py-1">
			<small class="opacity-50">PUBLIC</small>
			<button
				class="btn-icon btn-icon-sm variant-soft"
				on:click={() => {
					personalOpen = !personalOpen;
				}}
			>
				{#if personalOpen}
					<span><Icon name="minus" width="8px" height="8px" /></span>
				{:else}
					<span><Icon name="plus" width="8px" height="8px" /></span>
				{/if}
			</button>
		</div>
		<hr />
		<ul class={`list ${personalOpen ? '' : 'hidden'}`}>
			{#if $roomStore.loading}
				<div class="space-y-5 py-2">
					{#each new Array(PERSONAL_ROOM_PLACEHOLDER_COUNT).fill(0) as _}
						<div class="animate-pulse placeholder h-10" />
					{/each}
				</div>
			{:else if $roomStore.data}
				{#each $roomStore.data.filter((e) => e.type == RoomType.PersonalRoom) as room}
					<RoomItem {room} />
				{/each}
			{/if}
		</ul>
	</nav>
</div>

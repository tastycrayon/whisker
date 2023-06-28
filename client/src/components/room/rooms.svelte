<script lang="ts">
	import {
		Avatar,
		modalStore,
		type ModalComponent,
		type ModalSettings
	} from '@skeletonlabs/skeleton';
	import JoinRoom from './create-room.svelte';
	import { onMount } from 'svelte';
	import { refreshRooms, roomStore } from '$lib/store';
	import Icon from '../icon.svelte';
	import { RoomType } from '$lib/types';
	import RoomItem from './roomlist-item.svelte';

	// modal
	const modalComponent: ModalComponent = {
		// Pass a reference to your custom component
		ref: JoinRoom,
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

	onMount(async () => {
		await refreshRooms();
		console.log({ first: $roomStore });
	});
	$: defaultOpen = true;
	$: publicOpen = true;
</script>

<!-- List -->
<div class="p-4 space-y-4 overflow-y-auto">
	<nav class="list-nav">
		<div class="w-full inline-flex justify-between items-center pb-2">
			<button class="btn btn-sm variant-filled" on:click={handleModal}>
				<span>💀</span>
				<span>Create Room</span>
			</button>
			<button
				class="btn-icon btn-icon-sm variant-soft"
				disabled={$roomStore.loading}
				on:click={refreshRooms}
			>
				<span><Icon name="arrow-path" width="16px" height="16px" /></span>
			</button>
		</div>
		<hr />
		<div class="w-full inline-flex justify-between items-center py-1">
			<small class="opacity-50">DEFAULT</small>
			<button
				class="btn-icon btn-icon-sm variant-soft"
				on:click={() => {
					defaultOpen = !defaultOpen;
				}}
			>
				{#if defaultOpen}
					<span><Icon name="minus" width="8px" height="8px" /></span>
				{:else}
					<span><Icon name="plus" width="8px" height="8px" /></span>
				{/if}
			</button>
		</div>
		<hr />
		<ul class={`list ${defaultOpen ? '' : 'hidden'}`}>
			{#if $roomStore.loading}
				<div class="space-y-5 py-2">
					{#each new Array(7).fill(0) as _}
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
					{#each new Array(3).fill(0) as _}
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

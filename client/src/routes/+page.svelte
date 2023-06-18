<script lang="ts">
	import { COOKIE_OPTIONS } from '$lib/constant';
	import { currentUser, pb } from '$lib/pocketbase';
	import { json } from '@sveltejs/kit';
	import type { SendOptions } from 'pocketbase';
	import type { IRoom } from '$lib/types';
	import { onMount } from 'svelte';
	import { modalStore, type ModalComponent, type ModalSettings } from '@skeletonlabs/skeleton';
	import Login from '$components/login.svelte';

	const modalComponent: ModalComponent = {
		ref: Login,
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
</script>

<!-- <form method="POST" on:submit|preventDefault={JoinRoomFormHandler}>
	<input type="text" name={FormFieldKey.RoomID} id={FormFieldKey.RoomID} />
	<button type="submit">join room</button>
	<br />
</form>

<form method="POST" on:submit|preventDefault={CreateRoomFormHandler}>
	<input type="text" name={FormFieldKey.RoomID} id={FormFieldKey.RoomID} />
	<button type="submit">create room</button>
	<br />
</form>
<br /> -->
<!-- 
{JSON.stringify({ rooms })}
<button on:click={fetchRooms}>Get Rooms</button> -->

<!-- {#if rooms}
	{#each rooms as room}
		<a href={'rooms/' + room.roomId}>{room.roomName}</a> &emsp;
	{/each}
{/if} -->
<h1>test</h1>
{#if !$currentUser}
	<button on:click={handleModal}>login</button>
{/if}

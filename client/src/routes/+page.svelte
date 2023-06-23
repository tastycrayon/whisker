<script lang="ts">
	import { COOKIE_OPTIONS } from '$lib/constant';
	import { currentUser, pb } from '$lib/pocketbase';
	import { json } from '@sveltejs/kit';
	import type { SendOptions } from 'pocketbase';
	import type { IRoom } from '$lib/types';
	import { onMount } from 'svelte';
	import { modalStore, type ModalComponent, type ModalSettings } from '@skeletonlabs/skeleton';
	import Login from '$components/login.svelte';
	import Rooms from '$components/rooms.svelte';
	import AuthProviders from '$components/auth-providers.svelte';
	import Icon from '$components/icon.svelte';

	export let user: any;

	const modalComponent: ModalComponent = {
		// Pass a reference to your custom component
		ref: AuthProviders,
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
<div class="container mx-auto flex justify-center items-center flex-auto h-full">
	<div class="p-4 max-w-[640px] space-y-3">
		<h1 class="h1">
			Whisker<span
				class="bg-gradient-to-br from-pink-500 to-violet-500 bg-clip-text text-transparent box-decoration-clone"
				>Chat</span
			>
		</h1>
		<p class="text-sm">Revitalize your passion. Connect. Heal. Code.</p>
		<br />
		<p>
			Introducing "CodeHeal," a revolutionary chat app designed exclusively for programmers who are
			battling burnout and depression. We understand the unique challenges and pressures that come
			with a career in programming, and we're here to support you on your journey to recovery and
			well-being. CodeHeal provides a safe, judgment-free space where programmers can connect with
			like-minded individuals, share their experiences, and find solace in a supportive community.
		</p>
		{#if !$currentUser}
			<button on:click={handleModal}>login</button>
		{/if}
		<button type="button" class="btn variant-filled m-0" on:click={handleModal}>
			<Icon name="login" fill width="24px" height="24px" />
			<span>SIGN IN</span>
		</button>
	</div>
</div>

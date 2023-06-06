<script lang="ts">
	import { COOKIE_OPTIONS } from '$lib/constant';
	import { pb } from '$lib/pocketbase';
	import { json } from '@sveltejs/kit';
	import type { SendOptions } from 'pocketbase';
	import type { IRoom } from '$lib/types';
	import { onMount } from 'svelte';
	enum FormFieldKey {
		RoomID = 'roomID'
	}
	const JoinRoomFormHandler = async (e: SubmitEvent) => {
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);
		const roomID = formData.get(FormFieldKey.RoomID)?.toString()?.toLowerCase();
		if (!roomID) return; // TODO: validate later
		// try {
		// 	console.log({ token, record });
		// } catch (err) {
		// 	console.log('Error: ', err);
		// }
	};

	const CreateRoomFormHandler = async (e: SubmitEvent) => {
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);
		const roomID = formData.get(FormFieldKey.RoomID)?.toString()?.toLowerCase();
		if (!roomID) return; // TODO: validate later
		try {
			const res = await pb.send('/rooms', {
				method: 'POST',
				body: JSON.stringify({
					roomname: 'room' + roomID,
					roomId: roomID
				})
			});
		} catch (err) {
			console.log('Error: ', err);
		}
	};

	const fetchRooms = async <T>() => {
		type Result = { loading: boolean; error: undefined | Error; data: T | undefined };
		const result: Result = { loading: true, error: undefined, data: undefined };
		try {
			const { data, error } = await pb.send('/rooms', {
				method: 'GET'
			});
			if (error) throw error;
			(result.loading = false), (result.data = data as T);
			return result;
		} catch (err: any) {
			(result.loading = false), (result.error = new Error(err));
			return result;
		}
	};
	let rooms: IRoom[] = [];
	onMount(async () => {
		const { data, error } = await fetchRooms<IRoom[]>();
		if (data) rooms = data;
	});
	$: console.log(rooms);
</script>

<form method="POST" on:submit|preventDefault={JoinRoomFormHandler}>
	<input type="text" name={FormFieldKey.RoomID} id={FormFieldKey.RoomID} />
	<button type="submit">join room</button>
	<br />
</form>

<form method="POST" on:submit|preventDefault={CreateRoomFormHandler}>
	<input type="text" name={FormFieldKey.RoomID} id={FormFieldKey.RoomID} />
	<button type="submit">create room</button>
	<br />
</form>
<br />

{JSON.stringify({ rooms })}
<button on:click={fetchRooms}>Get Rooms</button>

{#if rooms}
	{#each rooms as room}
		<a href={'rooms/' + room.roomId}>{room.roomName}</a> &emsp;
	{/each}
{/if}

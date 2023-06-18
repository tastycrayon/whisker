<script lang="ts">
	import { navigating } from '$app/stores';
	import { pb } from '$lib/pocketbase';
	import { refreshRooms } from '$lib/store';
	import { modalStore } from '@skeletonlabs/skeleton';

	enum FormFieldKey {
		RoomID = 'roomID'
	}

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
			if (res) modalStore.close();
			refreshRooms();
		} catch (err) {
			console.log('Error: ', err);
		}
	};
</script>

<div class="card p-4 max-w-[400px] mx-auto space-y-4 my-4">
	<h1 class="h1">Create Room</h1>
	<form class="space-y-2" method="POST" on:submit|preventDefault={CreateRoomFormHandler}>
		<label class="label">
			<small>Room Name</small>
			<input
				class="input"
				type="text"
				placeholder="Room Name"
				name={FormFieldKey.RoomID}
				id={FormFieldKey.RoomID}
			/>
		</label>
		<button class="btn variant-filled" type="submit">SUBMIT</button>
	</form>
</div>

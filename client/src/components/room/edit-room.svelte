<script lang="ts">
	import { ROOM_PATH } from '$lib/constant';
	import { currentUser, pb } from '$lib/pocketbase';
	import {
		currentRoom,
		deleteRoom,
		refreshParticipants,
		refreshRooms,
		roomStore
	} from '$lib/store';
	import {
		FileDropzone,
		clipboard,
		modalStore,
		type ToastSettings,
		toastStore
	} from '@skeletonlabs/skeleton';
	import Icon from '../icon.svelte';
	import type { ClientResponseError } from 'pocketbase';
	import { RoomType, type FormResError, type IRoom } from '$lib/types';
	import { goto } from '$app/navigation';
	import { firstConfirm } from '$lib/util';
	export let room: IRoom;
	enum FormFieldKey {
		RoomName = 'name',
		RoomDescription = 'description',
		RoomCover = 'cover',
		RoomCreatedBy = 'createdBy',
		Type = 'type'
	}
	let roomName = room.name;
	let roomCover: undefined | HTMLImageElement;
	let loading = false;
	let error: FormResError<FormFieldKey>;

	const makeErrObj = (code: FormFieldKey | 'unknown' = 'unknown', message = '') =>
		(error = {
			code,
			message
		});
	$: formatError = (key: FormFieldKey) => (error?.code && error.code == key ? error.message : '');
	$: setErrorClass = (keys: string[]) =>
		`input  ${error?.code && keys.includes(error.code) ? 'input-error' : ''}`;

	const onFileChange = (event: Event) => {
		const input = event.target as HTMLInputElement;
		const file = input?.files != null && input?.files.length > 0 ? input.files[0] : null;
		if (file) {
			const reader = new FileReader();
			reader.addEventListener('load', () => {
				if (roomCover && reader.result) roomCover.setAttribute('src', reader.result as string);
			});
			reader.readAsDataURL(file);
			return;
		}
	};

	const updateRoomFormHandler = async (e: SubmitEvent) => {
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		if (!$currentUser) return makeErrObj('unknown', 'Authentication failed, Please sign in.');
		formData.append(FormFieldKey.RoomCreatedBy, $currentUser.id);

		const cover = formData.get(FormFieldKey.RoomCover) as File;
		if (cover.size < 1000) formData.delete(FormFieldKey.RoomCover);

		const roomName = formData.get(FormFieldKey.RoomName)?.toString();

		if (!roomName) return makeErrObj(FormFieldKey.RoomName, 'Room name can not be empty.');
		if (roomName.length < 3)
			return makeErrObj(FormFieldKey.RoomName, 'Room name must have at least 3 letter.');

		formData.append(FormFieldKey.Type, RoomType.PersonalRoom);
		try {
			loading = true;
			const res = await pb.collection('rooms').update<IRoom>(room.id, formData);
			const toast: ToastSettings = {
				message: `Updated ${res.name} <a href="${
					ROOM_PATH + '/' + res.slug
				}" class="btn btn-sm variant-filled">View</a>`,
				background: 'variant-filled-success'
			};
			toastStore.trigger(toast);
			await refreshRooms();
		} catch (err: unknown) {
			const e = (err as ClientResponseError).response.data;
			const keys = [
				FormFieldKey.RoomName,
				FormFieldKey.RoomDescription,
				FormFieldKey.RoomCover,
				FormFieldKey.RoomCreatedBy
			];
			if (e !== undefined && keys.includes((Object.keys(e).pop() || '') as any))
				for (const key of keys) {
					if (e[key] !== undefined) {
						if (e[key].code === 'validation_not_unique') e[key].message = 'Room already exists.';
						makeErrObj(key, e[key].message);
					}
				}
			else makeErrObj('unknown', 'Failed to update room.');
		} finally {
			loading = false;
		}
	};
	let { loading: roomLoading, error: roomStoreError, data } = $roomStore;
	const handleDelete = async () => {
		modalStore.clear();
		const confirmDeletion = await firstConfirm();
		if (!confirmDeletion) return;
		deleteRoom(room.id);
		goto(ROOM_PATH);
		modalStore.close();
	};
</script>

<form class="space-y-2" method="POST" on:submit|preventDefault={updateRoomFormHandler}>
	<!-- cover -->
	<div class="inline-flex gap-2">
		<span>
			<figure
				class="avatar flex aspect-square justify-center items-center overflow-hidden bg-surface-400-500-token w-16 h-16 rounded-full"
			>
				<img bind:this={roomCover} class="aspect-square object-cover object-center" alt="" />
			</figure>
		</span>
		<FileDropzone
			accept=".jpg, .jpeg, .png, .webp, .gif"
			on:change={onFileChange}
			name={FormFieldKey.RoomCover}
			class="text-sm"
		/>
	</div>
	<small class="text-error-500">{formatError(FormFieldKey.RoomCover)}</small>
	<!-- room name -->
	<label class="label">
		<small class="inline-flex items-center">
			<span class="whitespace-nowrap">Room Name</span>
		</small>
		<input
			class={setErrorClass([FormFieldKey.RoomName, 'unknown'])}
			type="text"
			placeholder="My awesome room"
			name={FormFieldKey.RoomName}
			id={FormFieldKey.RoomName}
			bind:value={roomName}
		/>
		<small class="text-error-500">{formatError(FormFieldKey.RoomName)}</small>
	</label>

	<!-- room description -->
	<label class="label">
		<small class="inline-flex items-center">
			<span class="whitespace-nowrap">Description</span>
		</small>
		<textarea
			name={FormFieldKey.RoomDescription}
			class={setErrorClass([FormFieldKey.RoomDescription, 'unknown'])}
			rows="4"
			placeholder="Lorem ipsum dolor sit amet consectetur adipisicing elit."
			value={room.description}
		/>
		<small class="text-error-500">{formatError(FormFieldKey.RoomDescription)}</small>
	</label>
	<div class="flex justify-between">
		<button class="btn btn-sm variant-ghost-primary" disabled={loading} type="submit">UPDATE</button
		>
		<button
			class="btn btn-sm variant-ghost-error"
			type="button"
			on:click={handleDelete}
			disabled={roomLoading}>DELETE</button
		>
	</div>
</form>

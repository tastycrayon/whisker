<script lang="ts">
	import { ROOM_PATH } from '$lib/constant';
	import { currentUser, pb } from '$lib/pocketbase';
	import { refreshRooms } from '$lib/store';
	import {
		FileDropzone,
		clipboard,
		modalStore,
		type ToastSettings,
		toastStore
	} from '@skeletonlabs/skeleton';
	import Icon from '../icon.svelte';
	import type { ClientResponseError } from 'pocketbase';
	import type { FormResError, IRoom } from '$lib/types';

	enum FormFieldKey {
		RoomName = 'name',
		RoomSlug = 'slug',
		RoomDescription = 'description',
		RoomCover = 'cover',
		RoomCreatedBy = 'createdBy'
	}
	let roomName = '';
	let isTouched = false;
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
	const handleOnBlur = () => (isTouched = true);
	const convertToSlug = (e: string): string => {
		return e
			.toLowerCase()
			.replace(/ /g, '-')
			.replace(/[^\w-]+/g, '');
	};
	const getFullSlug = (value: string) => ROOM_PATH + '/' + convertToSlug(value);

	const CreateRoomFormHandler = async (e: SubmitEvent) => {
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);
		// const roomSlug = formData.get(FormFieldKey.RoomSlug)?.toString()?.toLowerCase();
		if (!$currentUser) return makeErrObj('unknown', 'Authentication failed, Please sign in.');
		formData.append(FormFieldKey.RoomCreatedBy, $currentUser.id);

		const roomName = formData.get(FormFieldKey.RoomName)?.toString();

		if (!roomName) return makeErrObj(FormFieldKey.RoomName, 'Room name can not be empty.');
		if (roomName.length < 3)
			return makeErrObj(FormFieldKey.RoomName, 'Room name must have at least 3 letter.');

		const roomSlug = convertToSlug(roomName);
		formData.append(FormFieldKey.RoomSlug, roomSlug);

		// const description = formData.get(FormFieldKey.RoomDescription)?.toString();
		// const cover = formData.get(FormFieldKey.RoomCover)?.toString();
		// const user = formData.get(FormFieldKey.RoomCreatedBy)?.toString();
		// // console.log({ roomName, roomSlug, description, cover, user });
		try {
			loading = true;
			const res = await pb.collection('rooms').create<IRoom>(formData);
			const toast: ToastSettings = {
				message: `New room ${res.name} <a href="${getFullSlug(
					res.slug
				)}" class="btn btn-sm variant-filled">View</a>`,
				background: 'variant-filled-success'
			};
			toastStore.trigger(toast);
			if (res) modalStore.close();
			refreshRooms();
		} catch (err: unknown) {
			const e = (err as ClientResponseError).response.data;
			const keys = [
				FormFieldKey.RoomName,
				FormFieldKey.RoomSlug,
				FormFieldKey.RoomDescription,
				FormFieldKey.RoomCover,
				FormFieldKey.RoomCreatedBy
			];
			if (e !== undefined)
				for (const key of keys) {
					if (e[key] !== undefined) {
						if (e[key].code === 'validation_not_unique') e[key].message = 'Room already exists.';
						makeErrObj(key, e[key].message);
					}
				}
		} finally {
			loading = false;
		}
	};
</script>

<div class="card p-4 max-w-[400px] mx-auto space-y-4 my-4">
	<h1 class="h3">Create Room</h1>
	<form class="space-y-2" method="POST" on:submit|preventDefault={CreateRoomFormHandler}>
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
				class={setErrorClass([FormFieldKey.RoomName, FormFieldKey.RoomSlug, 'unknown'])}
				type="text"
				placeholder="My awesome room"
				name={FormFieldKey.RoomName}
				id={FormFieldKey.RoomName}
				on:focus={handleOnBlur}
				bind:value={roomName}
			/>
			<small class="text-error-500">{formatError(FormFieldKey.RoomName)}</small>
		</label>
		<!-- room slug -->
		{#if isTouched}
			<button
				type="button"
				class="btn btn-sm code h-6 inline-flex items-center justify-between w-full text-left"
				use:clipboard={getFullSlug(roomName)}
			>
				<span class="break-words break-all">
					{ROOM_PATH + '/' + convertToSlug(roomName)}
				</span>&nbsp;
				<span>
					<Icon name="clipboard" width="16px" height="16px" />
				</span>
			</button>
			<small class="text-error-500">{formatError(FormFieldKey.RoomSlug)}</small>
		{:else}
			<div class="placeholder h-6" />
		{/if}
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
			/>
			<small class="text-error-500">{formatError(FormFieldKey.RoomDescription)}</small>
		</label>
		<button class="btn btn-sm variant-ghost-primary" disabled={loading} type="submit">SUBMIT</button
		>
	</form>
</div>

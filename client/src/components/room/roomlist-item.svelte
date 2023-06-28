<script lang="ts">
	import { page } from '$app/stores';
	import { ROOM_PATH } from '$lib/constant';
	import type { IRoom } from '$lib/types';
	import {
		Avatar,
		modalStore,
		type ModalSettings,
		type ModalComponent
	} from '@skeletonlabs/skeleton';
	import RoomInfo from './roomlist-info.svelte';
	import { participantStore } from '$lib/store';

	export let room: IRoom;
	$: classesActive = (href: string) =>
		href === $page.params.roomSlug ? '!bg-primary-500 text-white' : '';

	// modal
	const modalComponent: ModalComponent = {
		ref: RoomInfo,
		props: { room },
		slot: '<p>Skeleton</p>'
	};
	const handleModal = () => {
		const modal: ModalSettings = {
			type: 'component',
			component: modalComponent
		};
		modalStore.trigger(modal);
	};
</script>

<li class="my-3">
	<a href={`${ROOM_PATH}/${room.slug}`} class={'w-full ' + classesActive(room.slug)}>
		<span class="w-8">
			<Avatar src="https://source.unsplash.com/1/48x48" width="w-8" />
		</span>
		<span class="flex-auto whitespace-nowrap overflow-hidden text-ellipsis">{room.name}</span>
		<button class="w-1" on:click|preventDefault={handleModal}>⋮</button>
	</a>
</li>

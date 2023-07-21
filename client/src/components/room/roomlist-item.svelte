<script lang="ts">
	import { page } from '$app/stores';
	import { ROOM_PATH } from '$lib/constant';
	import { CollectionName, type IRoom } from '$lib/types';
	import {
		Avatar,
		modalStore,
		type ModalSettings,
		type ModalComponent,
		drawerStore
	} from '@skeletonlabs/skeleton';
	import RoomInfo from './roomlist-tab.svelte';
	import { generateAvatar } from '$lib/util';

	export let room: IRoom;
	$: classesActive = (href: string) =>
		href === $page.params.roomSlug ? '!bg-primary-500 text-white' : '';
	const cover = generateAvatar(CollectionName.Room, room.id, room.cover);
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

	const handleDrawyerClose = () => {
		drawerStore.close();
	};
</script>

<li class="my-3">
	<a
		href={`${ROOM_PATH}/${room.slug}`}
		on:click={handleDrawyerClose}
		class={'w-full ' + classesActive(room.slug)}
	>
		<span class="w-8">
			<Avatar src={cover} width="w-8" />
		</span>
		<span class="flex-auto whitespace-nowrap overflow-hidden text-ellipsis">{room.name}</span>
		<button class="btn-icon btn-icon-sm w-1" on:click|preventDefault={handleModal}>⋮</button>
	</a>
</li>

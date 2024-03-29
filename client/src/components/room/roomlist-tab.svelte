<script lang="ts">
	import { CollectionName, type IRoom } from '$lib/types';
	import { Avatar, Tab, TabGroup, clipboard } from '@skeletonlabs/skeleton';
	import Icon from '../icon.svelte';
	import RoomParticipants from './room-participants.svelte';
	import { participantStore, refreshSingleParticipant } from '$lib/store';
	import { ROOM_PATH } from '$lib/constant';
	import { generateAvatar, timeSince } from '$lib/util';
	import EditRoom from './edit-room.svelte';
	import { currentUser } from '$lib/pocketbase';

	export let room: IRoom;
	let showParticipantCount = false;
	let tabSet: number = 0;
	$: roomOwner = $participantStore.data.find((p) => p.id == room.createdBy);
	$: if (!roomOwner) refreshSingleParticipant(room.createdBy);

	const cover = generateAvatar(CollectionName.Room, room.id, room.cover);

	$: participantCount = $participantStore.data.reduce((acc, e) => {
		return e.roomSlug == room.slug ? acc + 1 : acc;
	}, 0);
</script>

<div class="card bg-initial p-4 w-[400px] mx-auto">
	<div class="flex aspect-square justify-center items-center flex-col">
		<div
			class={`w-full aspect-square bg-cover bg-center rounded`}
			style={`background-image:url('${cover}')`}
		>
			<div class="w-full h-full flex justify-start items-end bg-gradient-to-t from-black">
				<header class="w-full p-4">
					<h3 class="h3 pb-2 text-white" data-toc-ignore="">{room.name}</h3>
					<button
						type="button"
						class="btn btn-sm code h-6 inline-flex items-center justify-between w-full text-left"
						use:clipboard={`${window.location.host}${ROOM_PATH}/${room.slug}`}
					>
						<span class="break-words break-all">
							{ROOM_PATH + '/' + room.slug}
						</span>&nbsp;
						<span>
							<Icon name="clipboard" width="16px" height="16px" />
						</span>
					</button>
				</header>
			</div>
		</div>
	</div>
	<TabGroup>
		<Tab bind:group={tabSet} name="tab1" value={0}>Info</Tab>
		<Tab bind:group={tabSet} name="tab2" value={1}
			>Participants {#if showParticipantCount}{participantCount}{/if}</Tab
		>
		{#if room.createdBy == $currentUser?.id}
			<Tab bind:group={tabSet} name="tab3" value={2}>Edit</Tab>
		{/if}
		<!-- Tab Panels --->
		<svelte:fragment slot="panel">
			{#if tabSet === 0}
				<div class="p-0 space-y-4">
					<!-- <h3 class="h3" data-toc-ignore="">{room.name}</h3> -->

					<article class="whitespace-pre-line text-sm">
						{room.description}
					</article>
				</div>
			{:else if tabSet === 1}
				<div class="overflow-y-auto max-h-[300px]">
					<RoomParticipants {room} bind:showParticipantCount />
				</div>
			{:else if tabSet === 2 && room.createdBy == $currentUser?.id}
				<div class="overflow-y-auto max-h-[300px]">
					<EditRoom {room} />
				</div>
			{/if}
		</svelte:fragment>
	</TabGroup>

	<hr class="opacity-50 my-4" />
	<footer class="flex justify-start items-center space-x-4">
		<Avatar
			width="w-8"
			initials={roomOwner?.username.charAt(0) || ''}
			src={roomOwner ? generateAvatar(CollectionName.User, roomOwner.id, roomOwner.avatar) : 'A'}
			>{roomOwner
				? generateAvatar(CollectionName.User, roomOwner.id, roomOwner.avatar)
				: 'A'}</Avatar
		>
		<div class="flex-auto flex justify-between items-center">
			<h6 class="font-bold">By {roomOwner?.username || 'Admin'}</h6>
			<small>{timeSince(room.created)}</small>
		</div>
	</footer>
</div>

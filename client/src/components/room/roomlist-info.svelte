<script lang="ts">
	import type { IRoom, IUser, ParticipantRoom } from '$lib/types';
	import { Avatar, Tab, TabGroup, clipboard } from '@skeletonlabs/skeleton';
	import Icon from '../icon.svelte';
	import RoomParticipants from './room-participants.svelte';
	import { participantStore } from '$lib/store';
	import { PARTICIPANT_PATH, ROOM_PATH } from '$lib/constant';
	import { pb } from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { generateUserAvatar, timeSince } from '$lib/util';

	export let room: IRoom;
	let tabSet: number = 0;
	let roomOwner: ParticipantRoom | undefined = undefined;
	onMount(async () => {
		try {
			if (room.createdBy == 'Admin') return;
			const res = await pb.send<ParticipantRoom>(PARTICIPANT_PATH + '/' + room.createdBy, {
				method: 'GET'
			});
			roomOwner = res;
		} catch (err) {}
	});

	$: participantCount =
		$participantStore?.data?.reduce((total, item) => {
			return item.roomSlug == room.slug ? total + 1 : total;
		}, 0) || 0;
	const handleDelete = async () => {
		try {
			const result = await pb.collection('rooms').delete(room.id, {});
			console.log({ result });
		} catch (err) {
			console.log({ err });
		}
	};
</script>

<div class="card bg-initial p-4 w-[400px] mx-auto space-y-4 my-4">
	<header>
		<img
			src={generateUserAvatar(room.id, room.cover)}
			class="bg-black/50 w-full aspect-[21/9]"
			alt="Post"
		/>
	</header>
	<button on:click={handleDelete} class="variant-filled">delete</button>
	<TabGroup>
		<Tab bind:group={tabSet} name="tab1" value={0}>Info</Tab>
		<Tab bind:group={tabSet} name="tab2" value={1}>Participants {participantCount}</Tab>
		<!-- Tab Panels --->
		<svelte:fragment slot="panel">
			{#if tabSet === 0}
				<div class="p-0 space-y-4">
					<h3 class="h3" data-toc-ignore="">{room.name}</h3>
					<button
						type="button"
						class="btn btn-sm code h-6 inline-flex items-center justify-between w-full text-left"
						use:clipboard={`${ROOM_PATH}/${room.slug}`}
					>
						<span class="break-words break-all">
							{ROOM_PATH + '/' + room.slug}
						</span>&nbsp;
						<span>
							<Icon name="clipboard" width="16px" height="16px" />
						</span>
					</button>
					<article class="whitespace-pre-line">
						{room.description}
					</article>
				</div>
			{:else if tabSet === 1}
				<div class="overflow-y-auto max-h-[300px]">
					<RoomParticipants />
				</div>
			{/if}
		</svelte:fragment>
	</TabGroup>

	<hr class="opacity-50" />
	<footer class="p-4 flex justify-start items-center space-x-4">
		<Avatar width="w-8" src={roomOwner ? generateUserAvatar(roomOwner.id, roomOwner.avatar) : 'A'}
			>{roomOwner ? generateUserAvatar(roomOwner.id, roomOwner.avatar) : 'A'}</Avatar
		>
		<div class="flex-auto flex justify-between items-center">
			<h6 class="font-bold">By {roomOwner?.name || 'Admin'}</h6>
			<small>{timeSince(room.created)}</small>
		</div>
	</footer>
</div>

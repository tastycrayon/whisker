<script lang="ts">
	import type { IRoom } from '$lib/types';
	import { Tab, TabGroup, clipboard } from '@skeletonlabs/skeleton';
	import Icon from './icon.svelte';
	import RoomParticipants from './room-participants.svelte';
	import { participantStore } from '$lib/store';

	export let room: IRoom;
	let tabSet: number = 0;

	$: participantCount = $participantStore.reduce((total, item) => {
		return item.roomSlug == room.roomSlug ? total + 1 : total;
	}, 0);
</script>

<div class="card bg-initial card-hover p-4 max-w-[400px] mx-auto space-y-4 my-4">
	<header>
		<img
			src="https://source.unsplash.com/vjUokUWbFOs/400x175"
			class="bg-black/50 w-full aspect-[21/9]"
			alt="Post"
		/>
	</header>

	<TabGroup>
		<Tab bind:group={tabSet} name="tab1" value={0}>Info</Tab>
		<Tab bind:group={tabSet} name="tab2" value={1}>Participants {participantCount}</Tab>
		<!-- Tab Panels --->
		<svelte:fragment slot="panel">
			{#if tabSet === 0}
				<div class="p-4 space-y-4">
					<button use:clipboard={room.roomSlug} class="inline-flex items-center">
						<span>/{room.roomSlug}</span>
						<span class="px-3"><Icon name="clipboard" width="16px" height="16px" /></span>
					</button>
					<h3 class="h3" data-toc-ignore="">{room.roomName}</h3>
					<article>
						<p>
							Lorem ipsum dolor sit amet consectetur adipisicing elit. Numquam aspernatur provident
							eveniet eligendi cumque consequatur tempore sint nisi sapiente. Iste beatae laboriosam
							iure molestias cum expedita architecto itaque quae rem.
						</p>
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
		<figure
			class="avatar flex aspect-square text-surface-50 font-semibold justify-center items-center overflow-hidden isolate bg-surface-400-500-token w-8 rounded-full"
			data-testid="avatar"
		>
			<img
				class="avatar-image w-full h-full object-cover"
				style=""
				src="https://source.unsplash.com/YOErFW8AfkI/32x32"
				alt=""
			/>
		</figure>
		<div class="flex-auto flex justify-between items-center">
			<h6 class="font-bold">By {room.createdBy}</h6>
			<small>On 6/11/2023</small>
		</div>
	</footer>
</div>

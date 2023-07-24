<script lang="ts">
	import { Avatar, ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
	import { page } from '$app/stores';
	import { currentUser } from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { refreshParticipants, participantStore } from '$lib/store';
	import type { IRoom } from '$lib/types';

	export let room: IRoom;
	export let showParticipantCount = false;

	onMount(async () => {
		await refreshParticipants(room.slug);
		showParticipantCount = true;
	});
	// $: console.log({ first: $participantStore });
	$: currentParticipant =
		$participantStore.data?.find((e) => e.id == ($currentUser?.id || '')) || '';
	$: currentRoomParticipants = $participantStore.data?.filter((p) => p.roomSlug == room.slug) || [];
</script>

<!-- List -->
<ListBox active="variant-filled-primary">
	{#if $participantStore.data && $currentUser}
		{#if !$participantStore.loading && currentRoomParticipants.length == 0}
			<div class="text-sm text-center">No participants found.</div>
		{/if}
		{#each currentRoomParticipants as p}
			<ListBoxItem bind:group={currentParticipant} name="people" value={p}>
				<svelte:fragment slot="lead">
					<Avatar src={p.avatar} initials={p.username.charAt(0)} width="w-8" />
				</svelte:fragment>
				{p.username}
			</ListBoxItem>
		{/each}
	{:else if $participantStore.loading}
		<div class="space-y-5 py-2">
			{#each new Array(3).fill(0) as _}
				<div class="animate-pulse placeholder h-10" />
			{/each}
		</div>
	{/if}
</ListBox>

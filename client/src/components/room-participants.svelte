<script lang="ts">
	import { PARTICIPANT_PATH } from '$lib/constant';
	import {
		Avatar,
		ListBox,
		ListBoxItem,
		toastStore,
		type ToastSettings
	} from '@skeletonlabs/skeleton';
	import { page } from '$app/stores';
	import { currentUser, pb } from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
	import type { IUser } from '$lib/types';
	// Navigation List
	interface Person {
		avatar: string;
		id: string;
		roomId: string;
		username: string;
	}
	type IParticipantStore = { loading: boolean; error: any | undefined; data: Person[] | undefined };
	let res: IParticipantStore = { loading: true, error: undefined, data: undefined };

	export const getPartiCipants = async () => {
		try {
			const data = await pb.send<Person[] | undefined>(
				`${PARTICIPANT_PATH}/${$page.params.roomSlug}`,
				{
					method: 'GET'
				}
			);
			if (!data || !Array.isArray(data)) throw new Error('Failed loading rooms.');
			res = { loading: false, error: undefined, data: data };
			res.data?.map((d) => {
				d.avatar = `${PUBLIC_POCKETBASE_URL}/api/files/users/${d.id}/${d.avatar}`;
				return d;
			});
		} catch (err: any) {
			const t: ToastSettings = {
				message: err.message,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
			res = { loading: false, error: err, data: undefined };
		}
	};
	onMount(() => {
		getPartiCipants();
	});
	$: currentParticipant = res.data?.find((e) => e.id == ($currentUser?.id || ''));
</script>

<!-- List -->
<ListBox active="variant-filled-primary">
	{#if res.data && $currentUser}
		{#each res.data as p}
			<ListBoxItem bind:group={currentParticipant} name="people" value={p}>
				<svelte:fragment slot="lead">
					<Avatar src={p.avatar} width="w-8" />
				</svelte:fragment>
				{p.username}
			</ListBoxItem>
		{/each}
	{:else if res.loading}
		<div class="space-y-5 py-2">
			{#each new Array(3).fill(0) as _}
				<div class="animate-pulse placeholder h-10" />
			{/each}
		</div>
	{/if}
</ListBox>

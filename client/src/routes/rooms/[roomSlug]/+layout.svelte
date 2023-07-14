<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import DarkMode from '$components/dark-mode.svelte';
	import Icon from '$components/icon.svelte';
	import Rooms from '$components/room/rooms.svelte';
	import Typewriter from '$components/typewriter.svelte';
	import { DEFAULT_IMAGE, LOGIN_PATH, PROFILE_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
	import { WebSocketStore } from '$lib/socketStore';
	import { currentRoom, refreshRooms, roomStore } from '$lib/store';
	import { CollectionName } from '$lib/types';
	import { generateAvatar } from '$lib/util';
	import {
		Avatar,
		LightSwitch,
		ProgressRadial,
		getModeAutoPrefers,
		modeCurrent,
		setModeCurrent,
		setModeUserPrefers
	} from '@skeletonlabs/skeleton';
	import { onDestroy, onMount } from 'svelte';

	$: if (!$currentUser) goto(LOGIN_PATH);
	const { connect, close } = WebSocketStore();
	let loading = true;
	onMount(async () => {
		// if ($roomStore.data.length === 0) refreshRooms();
		await connect();
		loading = false;
		currentRoom.set($page.params.roomSlug);
	});
	onDestroy(() => {
		close();
	});
	$: console.log({ modeCurrent: $modeCurrent, gg: getModeAutoPrefers() });
</script>

<section class="card h-full">
	<div class="chat w-full h-full grid grid-cols-1 lg:grid-cols-[30%_1fr]">
		<!-- Navigation -->
		<div
			class="hidden lg:grid grid-rows-[auto_1fr_auto] max-h-screen border-r border-surface-500/30"
		>
			<!-- Header -->
			<header class="border-b border-surface-500/30 p-4">
				<input class="input" type="search" placeholder="Search..." />
			</header>
			<Rooms />
			<!-- Footer -->
			<footer class="border-t border-surface-500/30 space-y-2">
				{#if loading}
					<div>
						<aside class="w-full inline-flex gap-2 items-center justify-between px-4">
							<small>Connecting <Typewriter slogans={['.', '..', '...']} /></small>
							<p><ProgressRadial width="w-3" /></p>
						</aside>
						<hr />
					</div>
				{/if}

				{#if $currentUser}
					<div class="w-full inline-flex gap-2 items-center px-4 py-2">
						<span class="w-8">
							<a href={PROFILE_PATH}>
								<Avatar
									src={generateAvatar(CollectionName.User, $currentUser.id, $currentUser.avatar) ||
										DEFAULT_IMAGE}
									width="w-8"
								/>
							</a>
						</span>
						<a href={PROFILE_PATH} class="flex-auto">
							<span class="whitespace-nowrap overflow-hidden text-ellipsis"
								>{$currentUser.username}</span
							>
						</a>
						<DarkMode />
					</div>
				{/if}
			</footer>
		</div>
		<slot />
	</div>
</section>

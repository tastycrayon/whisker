<script lang="ts">
	import { goto } from '$app/navigation';
	import { navigating, page } from '$app/stores';
	import Footer from '$components/footer.svelte';
	import Head from '$components/head.svelte';
	import Icon from '$components/icon.svelte';
	import Sidebar from '$components/sidebar.svelte';
	import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import { DEFAULT_ROOM, LOGIN_PATH, ROOM_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
	import { socket } from '$lib/socketStore';
	import { currentRoom, roomStore } from '$lib/store';
	import { type DrawerSettings, drawerStore } from '@skeletonlabs/skeleton';
	import { hasContext, onDestroy, onMount, setContext } from 'svelte';

	$: if (!$currentUser) goto(LOGIN_PATH);
	if (!hasContext('socket')) setContext('socket', socket);

	$: currentRoom.set($page.params.roomSlug || DEFAULT_ROOM);

	const user = $currentUser;

	$: if ($navigating && user) {
		const { from, to } = $navigating;
		const run = async () => {
			if (!from?.params || !to?.params || !from.route || !to.route) return false;
			if (from.params.roomSlug == to.params.roomSlug || !(from.route.id == to.route.id))
				return false;
			// clean message feed
			socket.messageFeed = [];
			console.log('switched');
			await socket.closeConnection();
			const url = `${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${to.params.roomSlug}`;
			await socket.setUrl(url).connect(true);
		};
		run();
	}

	onMount(async () => {
		if (!$currentUser) return;
		try {
			const url = `${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${$page.params.roomSlug}`;
			const res = await socket.setUrl(url).connect();
			// console.log({ res });
		} catch (err) {
			console.log({ err });
		}
	});

	onDestroy(() => {
		socket.closeConnection();
	});

	const handleDrawer = () => {
		const drawerSettings: DrawerSettings = {
			id: 'phone-sidebar',
			width: ' max-w-[640px] '
		};
		drawerStore.open(drawerSettings);
	};

	$: currentRoomTitle = $roomStore.data.find((e) => e.slug == $page.params.roomSlug)?.name;
</script>

<Head title="{currentRoomTitle} - Whisker" />

<section class="h-full">
	<div class="chat w-full h-full grid grid-cols-1 lg:grid-cols-[30%_1fr] relative">
		<!-- drawer icon -->
		<div class="absolute top-0 left-0 p-2 lg:hidden flex items-center gap-4">
			<button class="btn-icon variant-ghost-primary" on:click={handleDrawer}>
				<Icon name="category" width="16px" height="16px" />
			</button>
			{#if currentRoomTitle}
				<small class="whitespace-nowrap overflow-hidden text-ellipsis">{currentRoomTitle}</small>
			{/if}
		</div>
		<!-- drawer -->
		<!-- Navigation -->
		<div
			class="hidden lg:grid grid-rows-[auto_1fr_auto] max-h-screen border-r border-surface-500/30"
		>
			<Sidebar />
			<!-- Footer -->
			<Footer />
		</div>
		<slot />
	</div>
</section>

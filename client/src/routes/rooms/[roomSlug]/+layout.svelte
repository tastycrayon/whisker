<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Footer from '$components/footer.svelte';
	import Icon from '$components/icon.svelte';
	import Sidebar from '$components/sidebar.svelte';
	import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import { DEFAULT_ROOM, LOGIN_PATH, ROOM_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
	import { WS, loading, reader, writer } from '$lib/socketStore';
	import { currentRoom } from '$lib/store';
	import { type DrawerSettings, drawerStore } from '@skeletonlabs/skeleton';
	import { onDestroy, onMount } from 'svelte';

	$: if (!$currentUser) goto(LOGIN_PATH);
	// const { connect, close } = WebSocketStore();
	currentRoom.set($page.params.roomSlug || DEFAULT_ROOM);
	const socket = new WS(
		`${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${$currentRoom}`,
		reader,
		writer,
		loading
	);

	onMount(async () => {
		try {
			const res = await socket.connect();
			// console.log({ res });
		} catch (err) {
			console.log({ err });
		}
	});
	onDestroy(() => {
		socket.closeConnection();
	});

	const drawerSettings: DrawerSettings = {
		id: 'phone-sidebar',
		width: ' max-w-[640px] '
	};
	const handleDrawer = () => {
		drawerStore.open(drawerSettings);
	};
</script>

<section class="h-full">
	<div class="chat w-full h-full grid grid-cols-1 lg:grid-cols-[30%_1fr] relative">
		<!-- drawer icon -->
		<div class="absolute top-0 left-0 p-2 lg:hidden">
			<button class="btn-icon variant-ghost-primary" on:click={handleDrawer}>
				<Icon name="category" width="16px" height="16px" />
			</button>
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

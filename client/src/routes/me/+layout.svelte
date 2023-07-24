<script lang="ts">
	import { goto } from '$app/navigation';
	import DarkMode from '$components/dark-mode.svelte';
	import Icon from '$components/icon.svelte';
	import { LOGIN_PATH, PROFILE_EDIT_PATH, ROOM_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
	import { AppBar, AppShell } from '@skeletonlabs/skeleton';

	$: if (!$currentUser) goto(LOGIN_PATH);
</script>

<AppShell>
	<svelte:fragment slot="header">
		<AppBar>
			<svelte:fragment slot="lead">
				<a href={ROOM_PATH} class="btn-icon btn-icon-sm variant-soft">
					<Icon name="arrow-left" width="18px" height="18px" />
				</a>
			</svelte:fragment>
			<a href={ROOM_PATH}>Back to Rooms</a>
			<svelte:fragment slot="trail">
				<DarkMode />
				<a class="btn-icon btn-icon-sm variant-soft" href={PROFILE_EDIT_PATH} title="Edit Profile">
					<Icon name="edit" width="18px" height="18px" />
				</a>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>
	<!-- (sidebarLeft) -->
	<!-- (sidebarRight) -->
	<!-- (pageHeader) -->
	<!-- Router Slot -->
	<slot />
	<!-- ---- / ---- -->
	<!-- (pageFooter) -->
	<!-- (footer) -->
</AppShell>

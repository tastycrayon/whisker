<script lang="ts">
	import DarkMode from '$components/dark-mode.svelte';
	import Icon from '$components/icon.svelte';
	import { ROOM_PATH, PROFILE_EDIT_PATH, PROFILE_PATH } from '$lib/constant';
	import {
		AppBar,
		AppShell,
		modalStore,
		type ModalComponent,
		type ModalSettings
	} from '@skeletonlabs/skeleton';
</script>

<AppShell>
	<svelte:fragment slot="header">
		<AppBar>
			<svelte:fragment slot="lead">
				<a href={ROOM_PATH} class="btn-icon btn-icon-sm variant-soft">
					<Icon name="arrow-left" width="18px" height="18px" />
				</a>
			</svelte:fragment>
			<a href={ROOM_PATH}>Back to rooms</a>
			<svelte:fragment slot="trail">
				<DarkMode />
				<a class="btn-icon btn-icon-sm variant-soft" href={PROFILE_PATH} title="Profile">
					<Icon name="user" width="18px" height="18px" />
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

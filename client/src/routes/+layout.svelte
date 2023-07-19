<script lang="ts">
	// Your selected Skeleton theme:
	import '../style/theme.postcss';

	// This contains the bulk of Skeletons required styles:
	// NOTE: this will be renamed skeleton.css in the v2.x release.
	import '@skeletonlabs/skeleton/styles/skeleton.css';

	// Finally, your application's global stylesheet (sometimes labeled 'app.css')
	import '../app.postcss';
	// <!-- css setting end -->
	import '../style/style.css';
	import {
		AppShell,
		Drawer,
		Modal,
		Toast,
		drawerStore,
		getModeOsPrefers,
		modeUserPrefers,
		setInitialClassState,
		setModeCurrent
	} from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';
	import Sidebar from '$components/sidebar.svelte';
	import Footer from '$components/footer.svelte';

	// Lifecycle
	onMount(() => {
		// Sync lightswitch with the theme
		if (!('modeUserPrefers' in localStorage)) {
			setModeCurrent(getModeOsPrefers());
		}
	});
	modeUserPrefers.subscribe((value) => {
		setModeCurrent(!!value);
	});
</script>

<svelte:head
	>{@html `<\u{73}cript nonce="%sveltekit.nonce%">(${setInitialClassState.toString()})();</script>`}</svelte:head
>
<Modal zIndex="z-[888]" />
<Toast zIndex="z-[999]" />
<Drawer zIndex="z-[50]">
	{#if $drawerStore.id === 'phone-sidebar'}
		<Footer closeBtn={true} />
		<hr />
		<Sidebar />
	{/if}
</Drawer>
<AppShell>
	<svelte:fragment slot="header">
		<!-- <Navbar {user} /> -->
	</svelte:fragment>
	<!-- (sidebarLeft) -->
	<!-- (sidebarRight) -->
	<!-- (pageHeader) -->
	<!-- Router Slot -->
	<!-- <div class="container px-10 mx-auto"> -->
	<slot />
	<!-- </div> -->
	<!-- ---- / ---- -->
	<!-- (pageFooter) -->
	<!-- <svelte:fragment slot="footer"> -->
	<!-- <Footer /> -->
	<!-- </svelte:fragment> -->
</AppShell>

<script lang="ts">
	import { ABOUT_PATH, ROOM_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
	import { modalStore, type ModalComponent, type ModalSettings } from '@skeletonlabs/skeleton';
	import AuthProviders from '$components/auth-providers.svelte';
	import Icon from '$components/icon.svelte';
	import Typewriter from '$components/typewriter.svelte';
	import Head from '$components/head.svelte';

	const modalComponent: ModalComponent = {
		// Pass a reference to your custom component
		ref: AuthProviders,
		// Add the component properties as key/value pairs
		props: { background: 'bg-red-500' },
		// Provide a template literal for the default component slot
		slot: '<p>Skeleton</p>'
	};
	const handleModal = () => {
		const modal: ModalSettings = {
			type: 'component',
			// Pass the component directly:
			component: modalComponent
		};
		modalStore.trigger(modal);
	};
	// $: if ($currentUser && $currentRoom) goto(ROOM_PATH + '/' + $currentRoom);
</script>

<Head />

<div class="container mx-auto flex justify-center items-center flex-auto h-full">
	<div class="p-4 max-w-[640px] space-y-3">
		<div class="inline-flex items-center">
			<h1 class="h1">
				Whis<span
					class="bg-gradient-to-br from-pink-500 to-violet-500 bg-clip-text text-transparent box-decoration-clone"
					>ker</span
				>
			</h1>
			<Icon name="cat" class="text-purple-600 ml-2" width="72px" height="72px" />
		</div>
		<p class="font-thin">
			Revitalize your passion. <Typewriter slogans={['Connect.', 'Heal.', 'Code.']} />
		</p>

		<br />
		<p>
			Introducing "Whisker," a revolutionary chat app designed exclusively for programmers who are
			battling burnout and depression. We understand the unique challenges and pressures that come
			with a career in programming, and we're here to support you on your journey to recovery and
			well-being. WhiskerChat provides a safe, judgment-free space where programmers can connect
			with like-minded individuals, share their experiences, and find solace in a supportive
			community.
		</p>
		{#if !$currentUser}
			<button type="button" class="btn variant-ghost-primary m-0" on:click={handleModal}>
				<span>Get Started</span>
			</button>
		{:else}
			<a class="btn variant-ghost-primary m-0" href={ROOM_PATH}>
				<span>Get Started</span>
			</a>
		{/if}
		<a class="btn variant-ghost-tertiary m-0" href={ABOUT_PATH}>
			<span>About Us</span>
		</a>
	</div>
</div>

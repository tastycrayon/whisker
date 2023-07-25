<script lang="ts">
	import { goto } from '$app/navigation';
	import Head from '$components/head.svelte';
	import { HOME_PATH, PROFILE_EDIT_PATH } from '$lib/constant';
	import { SignOut, currentUser } from '$lib/pocketbase';
	import { userPreference } from '$lib/store';
	import { CollectionName } from '$lib/types';
	import { generateAvatar, timeSince } from '$lib/util';
	import { SlideToggle } from '@skeletonlabs/skeleton';

	const handleSignOut = () => {
		SignOut(() => {
			goto(HOME_PATH, { invalidateAll: true });
		});
	};

	let welToggle = $userPreference.hideWelcomeMessage;
	const handleToggler = () => {
		userPreference.update((prev) => ({ ...prev, hideWelcomeMessage: welToggle }));
	};
</script>

<Head title="Profile - Whisker" />

<div class="container mx-auto">
	<div class="card p-4 max-w-[400px] mx-auto space-y-4 my-4">
		<div class="flex aspect-square justify-center items-center flex-col">
			<div
				class={`w-full aspect-square bg-cover bg-center rounded`}
				style={`background-image:url('${
					generateAvatar(CollectionName.User, $currentUser?.id || '', $currentUser?.avatar) || ''
				}')`}
			>
				<div class="w-full h-full flex justify-start items-end bg-gradient-to-t from-black">
					<header class="w-full p-4">
						<h3 class="h3 pb-2 text-white" data-toc-ignore="">{$currentUser?.username || ''}</h3>
						<small class="text-white font-bold"
							>Registered: {timeSince($currentUser?.created || new Date().toString())}</small
						>
						<br />
						<br />
						<a class="btn btn-sm variant-filled-tertiary" href={PROFILE_EDIT_PATH}>Edit Profile</a>
						<button type="button" class="btn btn-sm variant-filled-warning" on:click={handleSignOut}
							>Sign Out</button
						>
					</header>
				</div>
			</div>
		</div>

		<div class="inline-flex justify-between">
			<div />
			<SlideToggle
				name="slider-large"
				bind:checked={welToggle}
				on:change={handleToggler}
				active="bg-primary-500"
				size="sm"
			>
				Hide miscellaneous messages
			</SlideToggle>
		</div>
	</div>
</div>

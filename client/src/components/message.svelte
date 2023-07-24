<script lang="ts">
	import { getYoutubeVideoId, isImageUrl, isVideoUrl, isYoutubeUrl, timeSince } from '$lib/util';
	import {
		Avatar,
		modalStore,
		type ModalComponent,
		type ModalSettings
	} from '@skeletonlabs/skeleton';
	import User from './user.svelte';

	export let avatar = '';
	export let username = '';
	export let id = '';
	export let created = '';
	export let message = '';
	export let self = false;

	let parsedMessage = message;
	let renderHTML = false;
	if (isImageUrl(message)) {
		renderHTML = true;
		parsedMessage = `<img class="h-auto max-w-[320px] max-h-[240px] rounded-lg" src=${message} alt="image"/>`;
	}
	if (isVideoUrl(message)) {
		renderHTML = true;
		parsedMessage = `<video width="320" controls><source src=${message} type="video/mp4">Your browser does not support HTML video.</video>`;
	}
	if (isYoutubeUrl(message)) {
		const id = getYoutubeVideoId(message);
		if (id) {
			renderHTML = true;
			parsedMessage = `<iframe width="320" height="240" src="//www.youtube.com/embed/${id}" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>`;
		}
	}

	// modal
	const modalComponent: ModalComponent = {
		// Pass a reference to your custom component
		ref: User,
		// Add the component properties as key/value pairs
		props: { username, avatar },
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
</script>

<div class={`w-full min-w-[350px] inline-flex gap-2 ${self ? 'justify-end' : ''}`}>
	{#if !self}
		<div class="w-9">
			<Avatar
				on:click={handleModal}
				src={avatar}
				title={username}
				initials={username.charAt(0)}
				width="h-9"
			/>
		</div>
	{/if}
	<div
		class={`card p-4 space-y-2 rounded-3xl ${
			self ? 'variant-ghost-primary rounded-tr-none' : 'variant-ghost-tertiary rounded-tl-none'
		}`}
	>
		<header class={`flex ${self ? 'justify-end' : ''} items-center gap-3`}>
			<p class="font-bold">{username}</p>
			<small class="opacity-50">{timeSince(created)}</small>
		</header>
		<p class="break-words whitespace-pre-wrap break-all">
			{#if renderHTML}{@html parsedMessage}{:else}{parsedMessage}{/if}
		</p>
	</div>
	{#if self}
		<div class="w-9">
			<Avatar src={avatar} title={username} initials={username.charAt(0)} width="h-9" />
		</div>
	{/if}
</div>

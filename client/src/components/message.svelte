<script lang="ts">
	import { getYoutubeVideoId, isImageUrl, isVideoUrl, isYoutubeUrl, timeSince } from '$lib/util';
	import { Avatar } from '@skeletonlabs/skeleton';

	export let avatar = '';
	export let name = '';
	export let timestamp = '';
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
</script>

<div class={`w-full min-w-[350px] inline-flex gap-2 ${self ? 'justify-end' : ''}`}>
	{#if !self}
		<div class="w-8">
			<Avatar src={avatar} title={name} width="h-8" />
		</div>
	{/if}
	<div
		class={`card p-4 space-y-2 ${
			self ? 'variant-ghost-primary rounded-tr-none' : 'variant-ghost-tertiary rounded-tl-none'
		}`}
	>
		<header class={`flex ${self ? 'justify-end' : ''} items-center gap-3`}>
			<p class="font-bold">{name}</p>
			<small class="opacity-50">{timeSince(timestamp)}</small>
		</header>
		<p class="break-words whitespace-pre-wrap break-all">
			{#if renderHTML}{@html parsedMessage}{:else}{parsedMessage}{/if}
		</p>
	</div>
	{#if self}
		<div class="w-8">
			<Avatar src={avatar} title={name} width="h-8" />
		</div>
	{/if}
</div>

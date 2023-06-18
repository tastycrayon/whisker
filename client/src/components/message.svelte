<script lang="ts">
	import { timeSince } from '$lib/util';
	import { Avatar } from '@skeletonlabs/skeleton';

	export let avatar = '';
	export let name = '';
	export let timestamp = 0;
	export let message = '';
	export let self = false;
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
		<p class="break-words whitespace-pre-wrap break-all">{message}</p>
	</div>
	{#if self}
		<div class="w-8">
			<Avatar src={avatar} title={name} width="h-8" />
		</div>
	{/if}
</div>

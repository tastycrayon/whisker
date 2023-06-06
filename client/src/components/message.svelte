<script lang="ts">
	import { Avatar } from '@skeletonlabs/skeleton';

	export let avatar = '';
	export let name = '';
	export let timestamp = '';
	export let message = '';
	export let color = '';
	export let self = false;
</script>

<div class={`w-full inline-flex gap-2 ${self ? 'justify-end' : ''}`}>
	{#if !self}
		<div class="w-12">
			<Avatar src="https://i.pravatar.cc/?img={avatar}" width="h-12" />
		</div>
	{/if}
	<div
		class={`card p-4 space-y-2 ${color} ${
			self ? 'variant-soft rounded-tr-none' : 'rounded-tl-none'
		}`}
	>
		<header class="flex justify-between items-center">
			<p class="font-bold">{name}</p>
			<small class="opacity-50">{timestamp}</small>
		</header>
		<p class="break-words whitespace-pre-wrap break-all">{message}</p>
	</div>
	<div class="w-12">
		<Avatar src="https://i.pravatar.cc/?img={avatar}" width="h-12" />
	</div>
</div>

<script lang="ts">
	import { Avatar, ProgressRadial } from '@skeletonlabs/skeleton';
	import Typewriter from './typewriter.svelte';
	import { currentUser } from '$lib/pocketbase';
	import { DEFAULT_IMAGE, PROFILE_PATH } from '$lib/constant';
	import { generateAvatar } from '$lib/util';
	import { CollectionName } from '$lib/types';
	import DarkMode from './dark-mode.svelte';
	import { loading } from '$lib/socketStore';
	import Icon from './icon.svelte';
</script>

<footer class="border-t border-surface-500/30 space-y-2">
	{#if !navigator.onLine}
		<div>
			<aside class="w-full inline-flex gap-2 items-center justify-between px-4">
				<small>Connection not available</small>
				<p><Icon name="cell-signal-slash" width="16px" height="16px" /></p>
			</aside>
			<hr />
		</div>
	{/if}
	{#if $loading}
		<div>
			<aside class="w-full inline-flex gap-2 items-center justify-between px-4">
				<small>Connecting <Typewriter slogans={['.', '..', '...']} /></small>
				<p><ProgressRadial width="w-3" /></p>
			</aside>
			<hr />
		</div>
	{/if}

	{#if $currentUser}
		<div class="w-full inline-flex gap-2 items-center px-4 py-2">
			<span class="w-8">
				<a href={PROFILE_PATH}>
					<Avatar
						src={generateAvatar(CollectionName.User, $currentUser.id, $currentUser.avatar) ||
							DEFAULT_IMAGE}
						width="w-8"
					/>
				</a>
			</span>
			<a href={PROFILE_PATH} class="flex-auto">
				<span class="whitespace-nowrap overflow-hidden text-ellipsis">{$currentUser.username}</span>
			</a>
			<DarkMode />
		</div>
	{/if}
</footer>

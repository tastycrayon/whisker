<script>
	import { currentUser, initGoogleAuth } from '$lib/pocketbase';
	import { modalStore } from '@skeletonlabs/skeleton';
	import Icon from './icon.svelte';
	import { goto } from '$app/navigation';
	import { ROOM_PATH } from '$lib/constant';
</script>

<button
	type="button"
	class="btn variant-ringed w-full justify-between"
	on:click={async () => {
		await initGoogleAuth(modalStore.close);
		if ($currentUser?.id) goto(ROOM_PATH);
	}}
>
	<span><Icon name="google" fill width="24px" height="24px" /></span>
	<small>CONTINUE WITH GOOGLE</small>
	<span class="w-6" />
</button>

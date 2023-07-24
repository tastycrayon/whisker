<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import { toastStore, type ToastSettings } from '@skeletonlabs/skeleton';
	import Icon from './icon.svelte';
	import { ClientResponseError } from 'pocketbase';

	let inputVisible = false;
	let btnText = 'SendConfirmation';

	let email = '';
	let loading = false;
	const handleResend = async () => {
		try {
			loading = true;
			await pb.collection('users').requestVerification(email);
			const t: ToastSettings = {
				message: 'Verification email has been sent',
				background: 'variant-filled-success'
			};
			toastStore.trigger(t);
		} catch (err) {
			const t: ToastSettings = {
				message: (err as ClientResponseError).message,
				background: 'variant-filled-surface'
			};
			toastStore.trigger(t);
		} finally {
			loading = false;
		}
	};
</script>

<aside class="alert variant-outline-warning">
	<!-- Icon -->
	<div>
		<Icon name="warning" width="24px" height="24px" />
	</div>
	<!-- Message -->
	<div class="alert-message">
		<small>Check inbox for confirmation email</small>
		{#if inputVisible}
			<form class="flex gap-2" method="post" on:submit|preventDefault={handleResend}>
				<label class="label">
					<input
						bind:value={email}
						type="email"
						name="email"
						required
						class="input text-sm"
						placeholder="johncena@example.com"
					/>
				</label>
				<button type="submit" disabled={loading} class="btn btn-sm variant-outline-primary"
					>Resend</button
				>
			</form>
		{:else}
			<button class="btn btn-sm variant-outline-primary" on:click={() => (inputVisible = true)}
				>Send another confirmation email</button
			>
		{/if}
	</div>
</aside>

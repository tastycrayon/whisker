<script lang="ts">
	import { COOKIE_OPTIONS, ROOM_PATH } from '$lib/constant';
	import { pb } from '$lib/pocketbase';
	import { REGISTER_PATH } from '$lib/constant';
	import Icon from '$components/icon.svelte';
	import Separator from '$components/separator.svelte';
	import GoogleBtn from '$components/google-btn.svelte';
	import { goto } from '$app/navigation';
	import type { FormResError } from '$lib/types';
	import type { ClientResponseError } from 'pocketbase';
	import { ProgressRadial } from '@skeletonlabs/skeleton';

	enum FormFieldKey {
		EmailOrUsername = 'emailOrUsername',
		Password = 'password'
	}

	let error: FormResError<FormFieldKey>;
	let loading = false;
	const makeErrObj = (
		code: FormFieldKey | 'unknown' = 'unknown',
		message = ''
	): FormResError<FormFieldKey> => {
		error = { code, message };
		return error;
	};
	const formHandler = async (e: SubmitEvent) => {
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);
		const emailOrUsername = formData.get(FormFieldKey.EmailOrUsername)?.toString()?.toLowerCase();
		const password = formData.get(FormFieldKey.Password)?.toString();

		if (!emailOrUsername || emailOrUsername == '')
			return makeErrObj(FormFieldKey.EmailOrUsername, 'Email or username can not be empty.');

		// validate password
		if (!password || password == '')
			return makeErrObj(FormFieldKey.Password, 'Password can not be empty.');
		try {
			loading = true;
			const { token, record } = await pb
				.collection('users')
				.authWithPassword(emailOrUsername, password);
			pb.authStore.exportToCookie(COOKIE_OPTIONS);
			goto(ROOM_PATH);
		} catch (err) {
			const e = (err as ClientResponseError).response.data;
			const keys = [FormFieldKey.EmailOrUsername, FormFieldKey.Password];
			if (e !== undefined)
				for (const key of keys) {
					if (e[key] !== undefined) makeErrObj(key, e[key].message);
				}
		} finally {
			loading = false;
		}
	};

	$: formatError = (key: FormFieldKey) => (error?.code && error.code == key ? error.message : '');
	$: setErrorClass = (keys: string[]) =>
		`input  ${error?.code && keys.includes(error.code) ? 'input-error' : ''}`;

	let isPasswordVisible = false;
</script>

<div class="card p-4 max-w-[400px] mx-auto space-y-4 my-4">
	<h1 class="h3 font-bold">SIGN IN</h1>
	{#if error?.code == 'unknown'}<p class="text-error-500">{error.message}</p>{/if}

	<form method="POST" class="space-y-2" on:submit|preventDefault={formHandler}>
		<label class="label">
			<small>Email address or username</small>
			<input
				class={setErrorClass([FormFieldKey.EmailOrUsername, 'unknown'])}
				type="text"
				name={FormFieldKey.EmailOrUsername}
				placeholder="johncena@example.com"
			/>
			<small class="text-error-500">{formatError(FormFieldKey.EmailOrUsername)}</small>
		</label>
		<!-- password -->
		<label class="label">
			<small>Password</small>
			<div class="input-group input-group-divider grid-cols-[1fr_auto]">
				<input
					type={isPasswordVisible ? 'text' : 'password'}
					name={FormFieldKey.Password}
					placeholder="Password"
				/>
				<button
					class="btn btn-sm variant-filled-surface"
					type="button"
					on:click={() => {
						isPasswordVisible = !isPasswordVisible;
					}}
				>
					<Icon name={isPasswordVisible ? 'eye' : 'eye-hide'} width="24px" height="24px" />
				</button>
			</div>
			<small class="text-error-500">{formatError(FormFieldKey.Password)}</small>
		</label>
		<button type="submit" disabled={loading} class="btn variant-filled">
			<span>SUBMIT</span>
			{#if loading}
				<ProgressRadial width="w-4" stroke={150} track="stroke-secondary-500/30" />
			{/if}
		</button>
	</form>
	<div class="w-full">
		<span>Not registered? <a href={REGISTER_PATH} class="anchor">Create Account</a></span>
	</div>

	<Separator text="OR" />
	<GoogleBtn />
	<!-- <Facebook /> -->
</div>

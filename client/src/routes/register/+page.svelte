<script lang="ts">
	import { goto } from '$app/navigation';
	import GoogleBtn from '$components/google-btn.svelte';
	import Icon from '$components/icon.svelte';
	import Separator from '$components/separator.svelte';
	import { COOKIE_OPTIONS, LOGIN_PATH, ROOM_PATH } from '$lib/constant';
	import { currentUser, pb } from '$lib/pocketbase';
	import type { FormResError } from '$lib/types';
	import {
		FileDropzone,
		ProgressRadial,
		toastStore,
		type ToastSettings
	} from '@skeletonlabs/skeleton';
	import type { ClientResponseError } from 'pocketbase';
	import { onMount } from 'svelte';

	enum FormFieldKey {
		Email = 'email',
		Username = 'username',
		Password = 'password',
		PasswordConfirm = 'passwordConfirm',
		Avatar = 'avatar'
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
		const email = formData.get(FormFieldKey.Email)?.toString().trim().toLowerCase();
		const username = formData.get(FormFieldKey.Username)?.toString().trim().toLowerCase();
		const password = formData.get(FormFieldKey.Password)?.toString().trim();
		const passwordConfirm = formData.get(FormFieldKey.PasswordConfirm)?.toString().trim();

		if (!email || email == '')
			return makeErrObj(FormFieldKey.Email, 'Email address can not be empty.');
		if (!username || username == '')
			return makeErrObj(FormFieldKey.Username, 'Username can not be empty.');

		// validate password
		if (!password || password == '')
			return makeErrObj(FormFieldKey.Password, 'Password can not be empty.');
		if (password !== passwordConfirm)
			return makeErrObj(FormFieldKey.Password, 'Password and confirm password do not match.');

		try {
			loading = true;
			const { token, record } = await pb.collection('users').create({
				email,
				username,
				name,
				password,
				passwordConfirm
			});
			pb.authStore.exportToCookie(COOKIE_OPTIONS);
			const t: ToastSettings = {
				message: 'Check email for confirmation.',
				background: 'variant-filled-surface'
			};
			toastStore.trigger(t);
			setTimeout(() => goto(ROOM_PATH), 1500);
		} catch (err) {
			const e = (err as ClientResponseError).response.data;
			const keys = [
				FormFieldKey.Email,
				FormFieldKey.Username,
				FormFieldKey.Password,
				FormFieldKey.PasswordConfirm
			];
			if (e !== undefined && keys.includes((Object.keys(e).pop() || '') as any))
				for (const key of keys) {
					if (e[key] !== undefined) makeErrObj(key, e[key].message);
				}
			else makeErrObj('unknown', 'Failed to sign up.');
		} finally {
			loading = false;
		}
	};

	$: formatError = (key: FormFieldKey) => (error?.code && error.code == key ? error.message : '');
	$: setErrorClass = (keys: string[]) =>
		`input  ${error?.code && keys.includes(error.code) ? 'input-error' : ''}`;

	let isPasswordVisible = false;

	let avatar: undefined | HTMLImageElement;
	const onFileChange = (event: Event) => {
		const input = event.target as HTMLInputElement;
		const file = input?.files != null && input?.files.length > 0 ? input.files[0] : null;
		if (file) {
			const reader = new FileReader();
			reader.addEventListener('load', () => {
				if (avatar && reader.result) avatar.setAttribute('src', reader.result as string);
			});
			reader.readAsDataURL(file);
			return;
		}
	};

	onMount(() => {
		if ($currentUser?.id) goto(ROOM_PATH);
	});
</script>

<div class="card p-4 max-w-[400px] mx-auto space-y-4 my-4">
	<div class="inline-flex items-center">
		<h1 class="h1">Sign Up</h1>
		<Icon name="cat" class="text-purple-600" width="48px" height="48px" />
	</div>
	{#if error?.code == 'unknown'}<p class="text-error-500">{error.message}</p>{/if}
	<form method="POST" class="space-y-2" on:submit|preventDefault={formHandler}>
		<!-- username -->
		<label class="label">
			<small>Username</small>
			<input
				class={setErrorClass([FormFieldKey.Username, 'unknown'])}
				type="text"
				name={FormFieldKey.Username}
				placeholder="johncena911"
			/>
			<small class="text-error-500">{formatError(FormFieldKey.Username)}</small>
		</label>
		<!-- email address -->
		<label class="label">
			<small>Email Address</small>
			<input
				class={setErrorClass([FormFieldKey.Email, 'unknown'])}
				type="email"
				name={FormFieldKey.Email}
				placeholder="johnchena@chinamail.com"
			/>
			<small class="text-error-500">{formatError(FormFieldKey.Email)}</small>
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
		<!-- password confirm -->
		<label class="label">
			<small>Confirm Password</small>
			<input
				class={setErrorClass([FormFieldKey.PasswordConfirm, 'unknown'])}
				type={isPasswordVisible ? 'text' : 'password'}
				name={FormFieldKey.PasswordConfirm}
				placeholder="Must be the same as password"
			/>
			<small class="text-error-500">{formatError(FormFieldKey.PasswordConfirm)}</small>
		</label>
		<div class="h-1" />
		<!-- avatar -->
		<div class="inline-flex gap-2 w-full">
			<span>
				<figure
					class="avatar flex aspect-square justify-center items-center overflow-hidden bg-surface-400-500-token w-16 h-16 rounded-full"
				>
					<img bind:this={avatar} class="aspect-square object-cover object-center" alt="" />
				</figure>
			</span>
			<FileDropzone
				accept=".jpg, .jpeg, .png, .webp, .gif"
				on:change={onFileChange}
				name={FormFieldKey.Avatar}
				class="text-sm"
			/>
		</div>
		<button type="submit" disabled={loading} class="btn variant-filled">
			<span>SUBMIT</span>
			{#if loading}
				<ProgressRadial width="w-4" stroke={150} track="stroke-secondary-500/30" />
			{/if}
		</button>
		<div class="w-full">
			<span>Already have an account? <a href={LOGIN_PATH} class="anchor">Sign in here</a></span>
		</div>
	</form>
	<Separator text="OR" />
	<div class="space-y-2">
		<GoogleBtn />
	</div>
</div>

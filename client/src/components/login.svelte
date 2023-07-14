<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { pb } from '$lib/pocketbase';
	import { HttpStatusCode } from '$lib/statusCodes';
	import { modalStore } from '@skeletonlabs/skeleton';
	import Icon from './icon.svelte';

	interface FormResError {
		key: FormFieldKey | 'unknown'; // error can be unknown field
		message: string;
	}
	enum FormFieldKey {
		EmailOrUsername = 'emailOrUsername',
		Password = 'password'
	}
	const makeErrObj = (message = '', key: FormFieldKey | 'unknown' = 'unknown'): FormResError => ({
		message,
		key
	});

	let errors: FormResError[] = [];
	const handleSubmit = async (e: SubmitEvent) => {
		const formData = new FormData(e.target as HTMLFormElement);
		const emailOrUsername =
			formData.get(FormFieldKey.EmailOrUsername)?.toString().toLowerCase() || '';
		const password = formData.get(FormFieldKey.Password)?.toString() || '';
		// validate emailOrUsername
		if (!emailOrUsername || emailOrUsername == '')
			return (errors = [
				...errors,
				makeErrObj('Empty email address or username', FormFieldKey.EmailOrUsername)
			]);

		// validate password
		if (!password || password == '')
			return (errors = [...errors, makeErrObj('Empty password', FormFieldKey.Password)]);

		try {
			const { token, record } = await pb
				.collection('users')
				.authWithPassword(emailOrUsername.toLowerCase(), password);

			modalStore.close();
			if (token) invalidateAll();
		} catch (err) {
			// in case of user, avoid showing actual error to prevent fishing of information
			errors = [...errors, makeErrObj('Incorrect credentials!', 'unknown')];
		}
	};
	$: errors = errors;
	$: formatError = (key: string) => {
		const item = errors.find((e) => e.key == key);
		if (item) return item.message;
		return '';
	};
	$: setErrorClass = (keys: string[]) => {
		for (const key of keys) {
			const item = errors.find((e) => e.key == key);
			if (item) return `input input-error`;
		}
		return 'input';
	};
	let isPasswordVisible = false;
</script>

<div class="card p-4 max-w-[400px] mx-auto space-y-4 my-4">
	<h1 class="h1">SIGN IN</h1>

	{#if errors.length && errors[0].key == 'unknown'}<p class="text-error-500">
			{errors[0].message}
		</p>{/if}

	<form class="space-y-2" method="POST" on:submit|preventDefault={handleSubmit}>
		<label class="label">
			<small>Email or Username</small>
			<input
				class={setErrorClass(['email', 'unknown'])}
				type="text"
				name="emailOrUsername"
				placeholder="Email or Username"
			/>
			<small class="text-error-500">{formatError('emailOrUsername')}</small>
		</label>
		<label class="label">
			<small>Password</small>
			<div class="input-group input-group-divider grid-cols-[1fr_auto]">
				<input
					type={isPasswordVisible ? 'text' : 'password'}
					name="password"
					placeholder="Password"
				/>
				<button
					type="button"
					on:click={() => {
						isPasswordVisible = !isPasswordVisible;
					}}
				>
					<Icon name="eye" width="24px" fill height="24px" />
				</button>
			</div>
			<small class="text-error-500">{formatError('password')}</small>
		</label>
		<button type="submit" class="btn variant-filled">
			<small>SUBMIT</small>
		</button>

		<!-- <label class="label">
			<span>Username</span>
			<div class="input-group input-group-divider grid-cols-[1fr_auto]">
				<input type="text" placeholder="Enter Username..." />
				<div>
					<Icon name="eye" width="24px" fill height="24px" />
				</div>
			</div>
		</label> -->
	</form>
	<!-- <div class="w-full">
		<small>Not registered? <a href={REGISTER_PATH} class="anchor">Create Account</a></small>
	</div>

	<Separator text="OR" />
	<Google />
	<Facebook /> -->
</div>

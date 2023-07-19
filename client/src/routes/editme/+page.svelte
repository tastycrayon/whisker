<script lang="ts">
	import { currentUser, pb } from '$lib/pocketbase';
	import {
		FileDropzone,
		modalStore,
		type ToastSettings,
		toastStore,
		TabGroup,
		Tab
	} from '@skeletonlabs/skeleton';
	import type { ClientResponseError } from 'pocketbase';
	import { type IUser, type FormResError, CollectionName } from '$lib/types';
	import { goto } from '$app/navigation';
	import { firstConfirm } from '$lib/util';

	let tabSet: number = 0;

	enum FormFieldKey {
		Email = 'email',
		Username = 'username',
		Password = 'password',
		PasswordConfirm = 'passwordConfirm',
		Avatar = 'avatar'
	}

	let userAvatar: undefined | HTMLImageElement;
	let loading = false;
	let error: FormResError<FormFieldKey>;

	const makeErrObj = (code: FormFieldKey | 'unknown' = 'unknown', message = '') =>
		(error = {
			code,
			message
		});
	$: formatError = (key: FormFieldKey) => (error?.code && error.code == key ? error.message : '');
	$: setErrorClass = (keys: string[]) =>
		`input  ${error?.code && keys.includes(error.code) ? 'input-error' : ''}`;

	const onFileChange = (event: Event) => {
		const input = event.target as HTMLInputElement;
		const file = input?.files != null && input?.files.length > 0 ? input.files[0] : null;
		if (file) {
			const reader = new FileReader();
			reader.addEventListener('load', () => {
				if (userAvatar && reader.result) userAvatar.setAttribute('src', reader.result as string);
			});
			reader.readAsDataURL(file);
			return;
		}
	};

	const updateUserFormHandler = async (e: SubmitEvent) => {
		if (!$currentUser) return makeErrObj('unknown', 'Authentication failed, Please sign in.');
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		const cover = formData.get(FormFieldKey.Avatar) as File;
		if (cover.size < 1000) formData.delete(FormFieldKey.Avatar);

		const userName = formData.get(FormFieldKey.Username)?.toString();
		if (!userName || userName.length < 3)
			return makeErrObj(FormFieldKey.Username, 'Username must have at least 3 letter.');

		try {
			loading = true;
			const res = await pb.collection('users').update<IUser>($currentUser.id, formData);
			const toast: ToastSettings = {
				message: `Successfully updated`,
				background: 'variant-filled-success'
			};
			toastStore.trigger(toast);
			if (res) modalStore.close();
		} catch (err: unknown) {
			const e = (err as ClientResponseError).response.data;
			const keys = [FormFieldKey.Username, FormFieldKey.Avatar];
			if (e !== undefined && keys.includes((Object.keys(e).pop() || '') as any))
				for (const key of keys) {
					if (e[key] !== undefined) {
						if (e[key].code === 'validation_not_unique') e[key].message = 'User already exists.';
						makeErrObj(key, e[key].message);
					}
				}
			else makeErrObj('unknown', 'Failed to update user.');
		} finally {
			loading = false;
		}
	};

	let deleteLoading = false;
	const handleDelete = async () => {
		if (!$currentUser?.id) return;
		modalStore.clear();
		const confirmDeletion = await firstConfirm();
		if (!confirmDeletion) return;
		try {
			deleteLoading = true;
			const result = await pb.collection(CollectionName.User).delete($currentUser.id, {});
			if (!result) throw new Error('Failed deleting user.');
			pb.authStore.clear();
			goto('/');
		} catch (err: any) {
			const t: ToastSettings = { message: err.message, background: 'variant-filled-error' };
			toastStore.trigger(t);
		} finally {
			deleteLoading = false;
		}
	};

	// password
	let passwordLoading = false;
	const handlePasswordResetSubmit = async () => {
		try {
			passwordLoading = true;
			const res = await pb.collection('users').requestPasswordReset($currentUser?.email);
			const toast: ToastSettings = {
				message: 'Password reset email has been sent.',
				background: 'variant-filled-success'
			};
			toastStore.trigger(toast);

			modalStore.close();
		} catch (err) {
			const e = (err as ClientResponseError).response.data;
			const keys = [FormFieldKey.Email, FormFieldKey.Password];
			if (e !== undefined && keys.includes((Object.keys(e).pop() || '') as any))
				for (const key of keys) {
					if (e[key] !== undefined) makeErrObj(key, e[key].message);
				}
			else makeErrObj('unknown', 'Failed to reset password');
		} finally {
			passwordLoading = false;
		}
	};
	// email
	let emailLoading = false;
	const handleEmailResetFormSubmit = async (e: SubmitEvent) => {
		if (!$currentUser) return makeErrObj('unknown', 'Authentication failed, Please sign in.');
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		const email = formData.get(FormFieldKey.Email)?.toString();
		if (!email || email.length < 3)
			return makeErrObj(FormFieldKey.Username, 'Email address must have at least 3 letter.');

		try {
			emailLoading = true;
			const res = await pb.collection('users').requestEmailChange(email);
			const toast: ToastSettings = {
				message: 'Confirmation message has been sent to your new email.',
				background: 'variant-filled-success'
			};
			toastStore.trigger(toast);
			modalStore.close();
		} catch (err) {
			const e = (err as ClientResponseError).response.data;
			const keys = [FormFieldKey.Email, FormFieldKey.Password];
			if (e !== undefined && keys.includes((Object.keys(e).pop() || '') as any))
				for (const key of keys) {
					if (e[key] !== undefined) makeErrObj(key, e[key].message);
				}
			else makeErrObj('unknown', 'Failed to change email');
		} finally {
			emailLoading = false;
		}
	};
</script>

<div class="card p-4 max-w-[400px] mx-auto my-4">
	<h2 class="h4">Edit</h2>
	<br />
	<TabGroup>
		<Tab bind:group={tabSet} name="tab1" value={0}>
			<span>Profile</span>
		</Tab>
		<Tab bind:group={tabSet} name="tab2" value={1}>Email</Tab>
		<Tab bind:group={tabSet} name="tab3" value={2}>Account</Tab>
		<!-- Tab Panels --->
		<svelte:fragment slot="panel">
			{#if tabSet === 0}
				{#if error?.code == 'unknown'}<p class="text-error-500">{error.message}</p>{/if}
				<form class="space-y-2" method="POST" on:submit|preventDefault={updateUserFormHandler}>
					<!-- avatar -->
					<small>Change Cover</small>
					<div class="inline-flex gap-2 w-full">
						<span>
							<figure
								class="avatar flex aspect-square justify-center items-center overflow-hidden bg-surface-400-500-token w-16 h-16 rounded-full"
							>
								<img
									bind:this={userAvatar}
									class="aspect-square object-cover object-center"
									alt=""
								/>
							</figure>
						</span>
						<FileDropzone
							accept=".jpg, .jpeg, .png, .webp, .gif"
							on:change={onFileChange}
							name={FormFieldKey.Avatar}
							class="text-sm"
						/>
					</div>
					<small class="text-error-500">{formatError(FormFieldKey.Avatar)}</small>
					<!-- room name -->
					<label class="label">
						<small class="inline-flex items-center">
							<span class="whitespace-nowrap">Username</span>
						</small>
						<input
							class={setErrorClass([FormFieldKey.Username, 'unknown'])}
							type="text"
							placeholder="My awesome room"
							name={FormFieldKey.Username}
							id={FormFieldKey.Username}
							value={$currentUser?.username || ''}
						/>
						<small class="text-error-500">{formatError(FormFieldKey.Username)}</small>
					</label>

					<button class="btn btn-sm variant-filled" disabled={loading} type="submit">UPDATE</button>
				</form>
			{:else if tabSet === 1}
				<form
					class="space-y-2 my-4"
					method="POST"
					on:submit|preventDefault={handleEmailResetFormSubmit}
				>
					<!-- email address -->
					<label class="label">
						<small>Email Address</small>
						<input
							class={setErrorClass([FormFieldKey.Email])}
							type="email"
							name={FormFieldKey.Email}
							placeholder="johnchena@chinamail.com"
						/>
						<small class="text-error-500">{formatError(FormFieldKey.Email)}</small>
					</label>
					<button
						class="btn btn-sm variant-ghost-secondary"
						type="button"
						disabled={passwordLoading}
						on:click={handlePasswordResetSubmit}>CHANGE MY EMAIL</button
					>
				</form>
			{:else if tabSet === 2}
				<div class="space-y-4">
					<button
						class="btn btn-sm variant-ghost-primary"
						type="button"
						disabled={passwordLoading}
						on:click={handlePasswordResetSubmit}>RESET PASSWORD</button
					>
					<br />
					<small>A confirmation will be sent to your email. </small>
					<hr />
					<button
						class="btn btn-sm variant-ghost-error"
						type="button"
						on:click={handleDelete}
						disabled={deleteLoading}>DELETE MY ACCOUNT</button
					>
					<br />
					<small>Permanently remove your record.</small>
				</div>
			{/if}
		</svelte:fragment>
	</TabGroup>
</div>

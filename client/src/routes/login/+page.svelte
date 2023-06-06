<script lang="ts">
	import { COOKIE_OPTIONS } from '$lib/constant';
	import { pb } from '$lib/pocketbase';
	enum FormFieldKey {
		EmailOrUsername = 'emailOrUsername',
		Password = 'password'
	}
	const formHandler = async (e: SubmitEvent) => {
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);
		const emailOrUsername = formData.get(FormFieldKey.EmailOrUsername)?.toString()?.toLowerCase();
		const password = formData.get(FormFieldKey.Password)?.toString();

		if (!emailOrUsername || !password) return; // TODO: validate later
		try {
			const { token, record } = await pb
				.collection('users')
				.authWithPassword(emailOrUsername, password);
			pb.authStore.exportToCookie(COOKIE_OPTIONS);
			console.log({ token, record });
		} catch (err) {
			console.log('Error: ', err);
		}
	};
</script>

<form method="POST" on:submit|preventDefault={formHandler}>
	<input type="text" name="emailOrUsername" id="emailOrUsername" value="mosi" />
	<br />
	<input type="password" name="password" id="password" value="12345678" />
	<br />
	<button type="submit">submit</button>
</form>

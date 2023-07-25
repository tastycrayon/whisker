<script lang="ts">
	import Head from '$components/head.svelte';
	import { GITHUB_URL, HOME_PATH, LOGIN_PATH, ROOM_PATH } from '$lib/constant';
	import { currentUser } from '$lib/pocketbase';
</script>

<Head title="About - Whisker" />

<div class="container mx-auto p-8 space-y-8 flex flex-col justify-center min-h-screen">
	<h1 class="h1 font-bold text-center">About Us</h1>

	<p>
		Introducing "Whisker," a revolutionary chat app designed exclusively for programmers who are
		battling burnout and depression. We understand the unique challenges and pressures that come
		with a career in programming, and we're here to support you on your journey to recovery and
		well-being. Whisker provides a safe, judgment-free space where programmers can connect with
		like-minded individuals, share their experiences, and find solace in a supportive community.
	</p>
	<p>
		With Whisker, you can engage in meaningful conversations, seek advice from experienced
		professionals, and discover practical strategies to cope with burnout and depression. Our
		platform offers a range of specialized features tailored to the needs of programmers, including
		dedicated channels for technical discussions, career guidance, and mental health resources.
		Whether you're seeking emotional support, inspiration, or simply a place to unwind, Whisker is
		your haven in the digital world.
	</p>
	<section>
		{#if $currentUser}
			<a class="btn variant-filled-primary" href={ROOM_PATH}>GO TO ROOMS</a>
		{:else}
			<a class="btn variant-filled-primary" href={LOGIN_PATH}>Sign In</a>
		{/if}
		<a class="btn variant-filled-secondary" href={HOME_PATH}>HOME</a>
		<a class="btn variant-filled-secondary" href={GITHUB_URL}>Github</a>
	</section>
</div>

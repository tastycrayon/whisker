<script lang="ts">
	import { DEFAULT_IMAGE, DEFAULT_OG_IMAGE } from '$lib/constant';

	export let title = 'Whisker - Connect Heal Code';
	export let description = `Introducing "Whisker", a revolutionary chat app designed exclusively for programmers who are battling burnout and depression. We understand the unique challenges and pressures that come with a career in programming, and we're here to support you on your journey to recovery and well-being. Whisker provides a safe, judgment-free space where programmers can connect with like-minded individuals, share their experiences, and find solace in a supportive community.`;
	export let og_image = DEFAULT_OG_IMAGE || DEFAULT_IMAGE;
</script>

<!-- OPENGRAPH -->
<svelte:head>
	<title>{title}</title>
	<meta property="og:title" content={title || 'Whisker - Connect Heal Code'} />
	<meta property="og:description" content={description.trim()} />
	<meta property="og:image" content={og_image} />
</svelte:head>
<!-- OPENGRAPH -->

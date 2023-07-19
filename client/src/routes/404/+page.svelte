<script lang="ts">
	import { HOME_PATH, ROOM_PATH } from '$lib/constant';
</script>

<div class="container mx-auto p-8 space-y-8 text-center flex flex-col justify-center min-h-screen">
	<h1 class="h1 font-bold">404 NOT FOUND</h1>
	<p>Looks like you've found the doorway to the great nothing</p>
	<section>
		<a class="btn variant-filled-primary" href={ROOM_PATH}>GO TO ROOMS</a>
		<a class="btn variant-filled-secondary" href={HOME_PATH}>HOME</a>
	</section>
</div>

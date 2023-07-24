<script lang="ts">
	import { DEFAULT_IMAGE } from '$lib/constant';

	export let username = '';
	export let avatar = '';
</script>

<div class="card p-4 max-w-[400px] w-80 mx-auto space-y-4 my-4">
	<div class="flex aspect-square justify-center items-center flex-col">
		<div
			class={`w-full aspect-square bg-cover bg-center rounded`}
			style="background-image:url('{avatar || DEFAULT_IMAGE}')"
		>
			<div class="w-full h-full flex justify-start items-end bg-gradient-to-t from-black">
				<header class="w-full p-4">
					<h3 class="h3 pb-2 text-white" data-toc-ignore="">{username || ''}</h3>
					<!-- <small>Joined: {timeSince(created || new Date().toString())}</small>
					<br /> -->
				</header>
			</div>
		</div>
	</div>
</div>

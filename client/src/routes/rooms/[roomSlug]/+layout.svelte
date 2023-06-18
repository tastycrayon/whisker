<script>
	import Rooms from '$components/rooms.svelte';
	import { currentUser } from '$lib/pocketbase';
	import { LightSwitch } from '@skeletonlabs/skeleton';
</script>

<section class="card h-full">
	<div class="chat w-full h-full grid grid-cols-1 lg:grid-cols-[30%_1fr]">
		<!-- Navigation -->
		<div
			class="hidden lg:grid grid-rows-[auto_1fr_auto] max-h-screen border-r border-surface-500/30"
		>
			<!-- Header -->
			<header class="border-b border-surface-500/30 p-4">
				<input class="input" type="search" placeholder="Search..." />
			</header>
			<Rooms />
			<!-- Footer -->
			{#if $currentUser}
				<footer class="border-t border-surface-500/30 p-4">
					<div class="inline-flex">
						{$currentUser.name}
						<LightSwitch bgDark="bg-neutral-900/90" />
					</div>
				</footer>
			{/if}
		</div>
		<slot />
	</div>
</section>

<script lang="ts">
	import { Autocomplete, type AutocompleteOption } from '@skeletonlabs/skeleton';
	import Rooms from './room/rooms.svelte';
	import { refreshRooms, roomStore } from '$lib/store';
	import { onMount } from 'svelte';
	let inputDemo = '';

	let flavorOptions: AutocompleteOption[] = $roomStore.data.map((e) => ({
		label: e.name,
		value: e.name,
		keywords: [e.name, e.slug, e.createdBy, e.type, e.description].join(', '),
		meta: { slug: e.slug }
	}));

	function onFlavorSelection(event: any): void {
		console.log({ event });
		inputDemo = event.detail.label;
	}
	let isFocused = false;
	const onFocus = () => (isFocused = true);
	const onBlur = () => (isFocused = false);
	onMount(() => {
		if ($roomStore.data.length === 0) refreshRooms();
	});

	function onClickHandle(params: any) {
		console.log({ params });
	}
</script>

<!-- Header -->
<header class="border-b border-surface-500/30 p-4 relative">
	<input
		bind:value={inputDemo}
		on:focus={onFocus}
		on:blur={onBlur}
		class="input"
		type="search"
		placeholder="Search..."
	/>
	{#if isFocused}
		<div class="absolute top-100 left-0 w-full p-4">
			<div class="card w-full max-w-sm max-h-48 p-4 overflow-y-auto" tabindex="-1">
				<Autocomplete
					bind:input={inputDemo}
					options={flavorOptions}
					on:selection={onFlavorSelection}
					on:click={onClickHandle}
					class="text-sm"
				/>
			</div>
		</div>
	{/if}
</header>
<Rooms />

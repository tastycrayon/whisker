<script lang="ts">
	import { roomStore } from '$lib/store';
	import { Autocomplete, type AutocompleteOption } from '@skeletonlabs/skeleton';

	export let searchValue = '';

	let flavorOptions: AutocompleteOption[] = $roomStore.data.map((e) => ({
		label: e.name,
		value: e.name,
		keywords: [e.name, e.slug, e.createdBy, e.type, e.description].join(', '),
		meta: { slug: e.slug }
	}));

	function onFlavorSelection(event: any): void {
		console.log({ event });
		searchValue = event.detail.label;
	}

	function onClickHandle(params: any) {
		console.log({ params });
	}
</script>

<div class="absolute top-100 left-0 w-full p-4">
	<div class="card w-full max-w-sm max-h-48 p-4 overflow-y-auto" tabindex="-1">
		<Autocomplete
			bind:input={searchValue}
			options={flavorOptions}
			on:selection={onFlavorSelection}
			on:click={() => console.log('Test')}
			class="text-sm"
		/>
	</div>
</div>

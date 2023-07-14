<script lang="ts">
	import { modeUserPrefers, setModeUserPrefers } from '@skeletonlabs/skeleton';
	import Icon from './icon.svelte';
</script>

{#if $modeUserPrefers}
	<button
		class="btn-icon btn-icon-sm variant-soft rounded-full"
		on:click={() => setModeUserPrefers(false)}
	>
		<Icon name="sun" width="16px" height="16px" />
	</button>
{:else}
	<button
		class="btn-icon btn-icon-sm variant-soft rounded-full"
		on:click={() => setModeUserPrefers(true)}
	>
		<Icon name="moon" width="16px" height="16px" />
	</button>
{/if}

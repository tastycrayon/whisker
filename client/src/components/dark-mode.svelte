<script lang="ts">
	import { modeCurrent, setModeCurrent } from '@skeletonlabs/skeleton';
	import Icon from './icon.svelte';
</script>

<!-- <LightSwitch bgDark="bg-neutral-900/90" /> -->
{#if $modeCurrent}
	<button
		class="btn-icon btn-icon-sm variant-soft rounded-full"
		on:click={() => setModeCurrent(false)}
	>
		<Icon name="sun" width="16px" height="16px" />
	</button>
{:else}
	<button
		class="btn-icon btn-icon-sm variant-soft rounded-full"
		on:click={() => setModeCurrent(true)}
	>
		<Icon name="moon" width="16px" height="16px" />
	</button>
{/if}

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	export let slogans: string[] = [];
	let currentText: string = slogans.length ? slogans[0] : '';

	let i = 0,
		j = 0;
	export let speed = 230;

	let typing = '';

	let timeout: number | undefined;
	function typeWriter() {
		if (i >= currentText.length) {
			i = 0;
			typing = '';
			j++;
			if (j === slogans.length) j = 0;
			currentText = slogans[j];
		}
		if (i < currentText.length) {
			typing += currentText.charAt(i);
			i++;
			timeout = setTimeout(typeWriter, speed);
		}
	}
	onMount(() => {
		typeWriter();
	});
	onDestroy(() => {
		if (timeout) clearTimeout(timeout);
	});
</script>

<span class={`${$$restProps.class || ''} type-line`}>{typing}</span>

<style>
	.type-line {
		border-right: 2px solid rgba(255, 255, 255, 0.75);
		animation: blinkTextCursor 500ms steps(44) infinite normal;
	}
	@keyframes blinkTextCursor {
		from {
			border-right-color: rgba(255, 255, 255, 0.75);
		}
		to {
			border-right-color: transparent;
		}
	}
</style>

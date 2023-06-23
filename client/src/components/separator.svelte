<script lang="ts">
	export let text: string = 'OR';
</script>

<div class={'flex items-center w-full ' + ($$props.class ?? '')}>
	<div class="flex-grow"><hr /></div>
	<small class="m-2">{text}</small>
	<div class="flex-grow"><hr /></div>
</div>

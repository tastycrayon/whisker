<script>
	import Icon from '$components/icon.svelte';
	import Rooms from '$components/rooms.svelte';
</script>

<div class="card p-4 max-w-[640px] mx-auto space-y-4 my-4">
	<div class="inline-flex items-center">
		<h1 class="h1">Rooms</h1>
		<Icon name="cat" class="text-purple-600" width="48px" height="48px" />
	</div>
	<Rooms />
</div>

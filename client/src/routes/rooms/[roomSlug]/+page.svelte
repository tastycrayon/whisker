<script lang="ts">
	import Message from '$components/message.svelte';
	import { onMount } from 'svelte';
	import { navigating } from '$app/stores';
	import { currentUser } from '$lib/pocketbase';
	import { currentRoom, participantStore } from '$lib/store';
	import {
		CollectionName,
		type SendEvent,
		type SwapMessage,
		EventType,
		type TextMessage
	} from '$lib/types';
	import Icon from '$components/icon.svelte';
	import { generateAvatar } from '$lib/util';
	import { reader, writer } from '$lib/socketStore';
	import { modalStore, type ModalSettings } from '@skeletonlabs/skeleton';

	const modal: ModalSettings = {
		type: 'alert',
		title: 'Under Construction',
		body: 'This functionality is still under developement',
		response: (r: boolean) => {}
	};

	let elemChat: HTMLElement;
	let messageFeed: TextMessage[] = [];
	let currentMessage = '';

	const roomSwitcher = (): boolean => {
		if (!$navigating) return false;
		const { from, to } = $navigating;
		if (!from?.params || !to?.params || !from.route || !to.route) return false;
		if (from.params.roomSlug == to.params.roomSlug || !(from.route.id == to.route.id)) return false;
		const m: SwapMessage = {
			from: from.params.roomSlug,
			to: to.params.roomSlug
		};
		writer.set({ payload: JSON.stringify(m), type: EventType.Swap });
		writer.set(null);
		currentRoom.set(to.params.roomSlug);
		return true;
	};

	$: if ($navigating) {
		roomSwitcher();
		setTimeout(() => {
			scrollChatBottom('smooth');
		}, 0);
	}

	const sendMessage = () => {
		if (currentMessage === '') return;
		const m: SendEvent = {
			payload: currentMessage,
			type: EventType.Text
		};
		writer.set(m);
		currentMessage = '';
	};
	reader.subscribe((m) => {
		if (!m) return;
		switch (m.type) {
			case EventType.Text:
				messageFeed = [...messageFeed, m.payload];
				break;
			case EventType.ParticipantHistory:
				if (Array.isArray(m.payload)) participantStore.update((s) => ({ ...s, data: m.payload }));
				break;
			case EventType.MessageHistory:
				messageFeed = [];
				if (Array.isArray(m.payload)) messageFeed = m.payload.map((e) => e.payload);
				break;
			default:
				console.error('unknown message type', { m });
				break;
		}
		reader.set(null);
		setTimeout(() => {
			scrollChatBottom('smooth');
		}, 0);
	});
	// When DOM mounted, scroll to bottom
	onMount(() => {
		scrollChatBottom();
	});
	// end
	function scrollChatBottom(behavior?: ScrollBehavior): void {
		if (elemChat) elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}

	function onPromptKeydown(event: KeyboardEvent): void {
		if (['Enter'].includes(event.code)) {
			event.preventDefault();
			sendMessage();
		}
	}

	$: messageList = messageFeed
		// .filter((e) => e.roomSlug == $page.params.roomSlug)
		.sort((a, b) => (parseInt(a.id) > parseInt(b.id) ? 1 : -1));
</script>

<!-- Chat -->
<div class="flex flex-col flex-nowrap justify-end max-h-screen">
	<!-- <div class="grid grid-row-[1fr_auto]"> -->
	<!-- Conversation -->
	<section bind:this={elemChat} class="p-4 overflow-y-auto space-y-4">
		{#each messageList as m, i}
			<Message
				avatar={generateAvatar(CollectionName.User, m.sender.id, m.sender.avatar)}
				name={m.sender.username}
				timestamp={m.created}
				message={m.content.toString()}
				self={m.sender.username == $currentUser?.username}
			/>
		{/each}
	</section>
	<!-- Prompt -->
	<section class="border-t border-surface-500/30 p-4">
		<div class="input-group input-group-divider grid-cols-[auto_1fr_auto] rounded-container-token">
			<button
				class="input-group-shim"
				on:click={() => {
					modalStore.trigger(modal);
				}}>+</button
			>
			<textarea
				bind:value={currentMessage}
				class="bg-transparent border-0 ring-0"
				name="prompt"
				id="prompt"
				placeholder="Write a message..."
				rows="1"
				on:keydown={onPromptKeydown}
			/>
			<button
				class={currentMessage ? 'variant-filled-primary' : 'input-group-shim'}
				on:click={sendMessage}
			>
				<Icon name="paper-airplane" width="16px" height="16px" />
			</button>
		</div>
	</section>
</div>

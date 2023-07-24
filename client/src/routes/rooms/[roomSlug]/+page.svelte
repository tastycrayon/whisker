<script lang="ts">
	import Message from '$components/message.svelte';
	import { navigating, page } from '$app/stores';
	import { getContext, onDestroy } from 'svelte';
	import { currentUser } from '$lib/pocketbase';
	import { participantStore, userPreference } from '$lib/store';
	import { CollectionName, type SendEvent, EventType, type TextMessage } from '$lib/types';
	import Icon from '$components/icon.svelte';
	import { generateAvatar } from '$lib/util';
	import { WS, reader, writer } from '$lib/socketStore';
	import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import { ROOM_PATH } from '$lib/constant';

	const socket = getContext<WS>('socket');

	let elemChat: HTMLElement;
	let currentMessage = '';
	let hideInfo = $userPreference.hideWelcomeMessage;

	const sendMessage = () => {
		if (!socket.isOpened)
			return socket._notifyOpening(
				'Message delivery unsuccessful. Please for the WebSocket to establish a connection.'
			);
		if (currentMessage === '') return;
		const m: SendEvent = {
			payload: currentMessage,
			type: EventType.Text
		};
		writer.set(m);
		currentMessage = '';
	};
	// if not logged in dont connect
	const user = $currentUser;

	const roomSwitcher = async () => {
		if (!$navigating || !user) return;
		const { from, to } = $navigating;
		if (!from?.params || !to?.params || !from.route || !to.route) return false;
		if (from.params.roomSlug == to.params.roomSlug || !(from.route.id == to.route.id)) return false;

		// clean message feed
		socket.messageFeed = [];
		await socket.closeConnection();
		const url = `${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${to.params.roomSlug}`;
		await socket.setUrl(url).connect(true);
	};
	$: if ($navigating) roomSwitcher();

	const unsubscribeReader = reader.subscribe((m) => {
		if (!m) return;
		switch (m.type) {
			case EventType.Info:
				if (hideInfo) break; // do nothing
				socket.messageFeed = [...socket.messageFeed, m.payload];
				break;
			case EventType.Text:
				socket.messageFeed = [...socket.messageFeed, m.payload];
				break;
			case EventType.ParticipantHistory:
				if (!Array.isArray(m.payload)) break;
				participantStore.update((s) => ({ ...s, data: m.payload }));
				break;
			case EventType.MessageHistory:
				socket.messageFeed = [];
				if (!Array.isArray(m.payload)) break;
				const feed = hideInfo ? m.payload.filter((e) => e.type !== EventType.Info) : m.payload; // remove welcome messages
				socket.messageFeed = feed.map((e) => e.payload);

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

	function scrollChatBottom(behavior?: ScrollBehavior): void {
		if (elemChat) elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}

	function onPromptKeydown(event: KeyboardEvent): void {
		if (['Enter'].includes(event.code)) {
			event.preventDefault();
			sendMessage();
		}
	}
	onDestroy(() => {
		unsubscribeReader();
	});
	$: messageList = socket.messageFeed
		// .filter((e) => e.roomSlug == $page.params.roomSlug)
		.sort((a, b) => (parseInt(a.id) > parseInt(b.id) ? 1 : -1));
	// $: console.log({ messageList });
</script>

<!-- Chat -->
<div class="flex flex-col flex-nowrap justify-end max-h-screen">
	<!-- <div class="grid grid-row-[1fr_auto]"> -->
	<!-- Conversation -->
	<section bind:this={elemChat} class="p-2 overflow-y-auto overflow-x-hidden space-y-3">
		{#each messageList as m, i}
			<Message
				avatar={generateAvatar(CollectionName.User, m.sender.id, m.sender.avatar)}
				username={m.sender.username}
				id={m.id}
				created={m.created}
				message={m.content.toString()}
				self={m.sender.id == $currentUser?.id}
			/>
		{/each}
	</section>
	<section class="flex">
		<textarea
			class="textarea rounded-none m-0 p-4 flex-auto"
			bind:value={currentMessage}
			rows="3"
			name="prompt"
			id="prompt"
			placeholder="Write a message..."
			on:keydown={onPromptKeydown}
		/>
		<button
			class="rounded-none btn variant-filled-primary"
			disabled={currentMessage == ''}
			on:click={sendMessage}
		>
			<Icon name="paper-airplane" width="20px" height="20px" />
		</button>
	</section>
</div>

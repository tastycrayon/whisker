<script lang="ts">
	import Message from '$components/message.svelte';
	import { onMount } from 'svelte';
	import { page, navigating } from '$app/stores';
	import { currentUser, pb } from '$lib/pocketbase';
	import { PUBLIC_POCKETBASE_URL, PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import { ROOM_PATH } from '$lib/constant';
	import { messageStore, participantStore } from '$lib/store';
	import { onDestroy } from 'svelte';
	import { MessageType, type IMessage } from '$lib/types';
	import Icon from '$components/icon.svelte';
	import { generateUserAvatar } from '$lib/util';

	let socket: WebSocket | undefined;
	let elemChat: HTMLElement;
	let messageFeed: IMessage[] = [];
	let currentMessage = '';

	const roomSwitcher = () => {
		if (!$navigating) return;
		const { from, to } = $navigating;
		if (!from?.params || !to?.params || !from.route || !to.route) return;
		if (from.params.roomSlug == to.params.roomSlug || !(from.route.id == to.route.id)) return;
		if (socket) {
			socket.send(
				JSON.stringify({
					content: '',
					from: from.params.roomSlug,
					to: to.params.roomSlug,
					messageType: MessageType.Swap
				})
			);
		}
	};
	$: if ($navigating) {
		roomSwitcher();
		setTimeout(() => {
			scrollChatBottom('smooth');
		}, 0);
	}
	onMount(() => {
		socket = new WebSocket(`${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${$page.params.roomSlug}`);
		// Connection opened
		socket.onopen = (event) => {
			console.log("It's open");
		};
		// Listen for messages
		socket.onmessage = (event) => {
			const m = JSON.parse(event.data) as IMessage;
			console.log({ m });
			switch (m.messageType) {
				case MessageType.Swap:
					break;
				case MessageType.Text:
					messageStore.set(m);
					break;
				case MessageType.Welcome:
					messageStore.set(m);
					break;
				case MessageType.Bailout:
					messageStore.set(m);
					break;
				case MessageType.Ping:
					console.log({ ping: m });
					if (Array.isArray(m.content)) participantStore.set(m.content);
					break;
				case MessageType.History:
					messageFeed = [];
					if (Array.isArray(m.content)) m.content.forEach((e) => messageStore.set(e));
					break;
				default:
					messageStore.set(m);
					break;
			}
		};
		socket.onerror = (event) => {
			console.log('error');
		};
		socket.onclose = (event) => {
			console.log('closed');
		};
	});
	onDestroy(() => {
		if (socket) socket.close();
	});

	const sendMessage = (): void => {
		if (socket && socket.readyState <= 1) {
			socket.send(JSON.stringify({ content: currentMessage, messageType: MessageType.Text }));
		}
		currentMessage = '';
	};

	messageStore.subscribe((m: IMessage | null) => {
		if (!m) return;
		messageFeed = [...messageFeed, m];
		messageStore.set(null);
		setTimeout(() => {
			scrollChatBottom('smooth');
		}, 0);
	});

	function scrollChatBottom(behavior?: ScrollBehavior): void {
		elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}

	function onPromptKeydown(event: KeyboardEvent): void {
		if (['Enter'].includes(event.code)) {
			event.preventDefault();
			sendMessage();
		}
	}
	// When DOM mounted, scroll to bottom
	onMount(() => {
		scrollChatBottom();
	});
	$: messageList = messageFeed.filter((e) => e.roomSlug == $page.params.roomSlug);
</script>

<!-- Chat -->
<div class="flex flex-col flex-nowrap justify-end max-h-screen">
	<!-- <div class="grid grid-row-[1fr_auto]"> -->
	<!-- Conversation -->
	<section bind:this={elemChat} class="p-4 overflow-y-auto space-y-4">
		{#each messageList as bubble, i}
			<Message
				avatar={generateUserAvatar(bubble.sender.id, bubble.sender.avatar)}
				name={bubble.sender.userName}
				timestamp={bubble.created}
				message={bubble.content.toString()}
				self={bubble.sender.userName == $currentUser?.username}
			/>
		{/each}
	</section>
	<!-- Prompt -->
	<section class="border-t border-surface-500/30 p-4">
		<div class="input-group input-group-divider grid-cols-[auto_1fr_auto] rounded-container-token">
			<button class="input-group-shim">+</button>
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

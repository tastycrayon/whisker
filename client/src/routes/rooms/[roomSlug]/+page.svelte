<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { pb } from '$lib/pocketbase';
	import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import { ROOM_PATH } from '$lib/constant';
	import { messageStore } from '$lib/store';

	let messageInput = '';
	let messages: any = [];
	let socket: WebSocket | undefined;
	onMount(() => {
		socket = new WebSocket(`${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${$page.params.roomSlug}`);
		// Connection opened
		socket.onopen = (event) => {
			console.log("It's open");
		};
		// Listen for messages
		socket.onmessage = (event) => {
			messageStore.set(event.data);
		};
		socket.onerror = (event) => {};
		socket.onclose = (event) => {};
	});
	const sendMessage = (message: string) => {
		if (socket && socket.readyState <= 1) {
			socket.send(JSON.stringify({ message }));
		}
	};
	const handleSendMessage = () => {
		sendMessage(messageInput);
		messageInput = '';
	};

	messageStore.subscribe((value: string) => {
		messages = [...messages, value];
		messageStore.set('');
	});
	// new start
	// Components
	import { Avatar, CodeBlock, ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
	import Message from '$components/message.svelte';
	import { bubble } from 'svelte/internal';

	// Docs Shell
	// Types
	interface Person {
		id: number;
		avatar: number;
		name: string;
	}
	interface MessageFeed {
		id: number;
		host: boolean;
		avatar: number;
		name: string;
		timestamp: string;
		message: string;
		color: string;
	}

	let elemChat: HTMLElement;
	const lorem =
		'token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb2xsZWN0aW9uSWQiOiJfcGJfdXNlcnNfYXV0aF8iLCJleHAiOjE2ODY5NjQ0NzksImlkIjoidmNwY213ZmdvdG9zYzNnIiwidHlwZSI6ImF1dGhSZWNvcmQifQ.rCvE_TChEzR1hB3Kdsgss6AhfbfNdkIyT5wp_P4xzuotoken eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb2xsZWN0aW9uSWQiOiJfcGJfdXNlcnNfYXV0aF8iLCJleHAiOjE2ODY5NjQ0NzksImlkIjoidmNwY213ZmdvdG9zYzNnIiwidHlwZSI6ImF1dGhSZWNvcmQifQ.rCvE_TChEzR1hB3Kdsgss6AhfbfNdkIyT5wp_P4xzuo';

	// Navigation List
	const people: Person[] = [
		{ id: 0, avatar: 14, name: 'Michael' },
		{ id: 1, avatar: 40, name: 'Janet' },
		{ id: 2, avatar: 31, name: 'Susan' },
		{ id: 3, avatar: 56, name: 'Joey' },
		{ id: 4, avatar: 24, name: 'Lara' },
		{ id: 5, avatar: 9, name: 'Melissa' }
	];
	let currentPerson: Person = people[0];

	// Messages
	let messageFeed: MessageFeed[] = [
		{
			id: 0,
			host: true,
			avatar: 48,
			name: 'Jane',
			timestamp: 'Yesterday @ 2:30pm',
			message: lorem,
			color: 'variant-soft-primary'
		},
		{
			id: 1,
			host: false,
			avatar: 14,
			name: 'Michael',
			timestamp: 'Yesterday @ 2:45pm',
			message: lorem,
			color: 'variant-soft-primary'
		},
		{
			id: 2,
			host: true,
			avatar: 48,
			name: 'Jane',
			timestamp: 'Yesterday @ 2:50pm',
			message: lorem,
			color: 'variant-soft-primary'
		},
		{
			id: 3,
			host: false,
			avatar: 14,
			name: 'Michael',
			timestamp: 'Yesterday @ 2:52pm',
			message: lorem,
			color: 'variant-soft-primary'
		}
	];
	let currentMessage = '';

	function scrollChatBottom(behavior?: ScrollBehavior): void {
		elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}

	function getCurrentTimestamp(): string {
		return new Date().toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });
	}

	function addMessage(): void {
		const newMessage = {
			id: messageFeed.length,
			host: true,
			avatar: 48,
			name: 'Jane',
			timestamp: `Today @ ${getCurrentTimestamp()}`,
			message: currentMessage,
			color: 'variant-soft-primary'
		};
		// Update the message feed
		messageFeed = [...messageFeed, newMessage];
		// Clear prompt
		currentMessage = '';
		// Smooth scroll to bottom
		// Timeout prevents race condition
		setTimeout(() => {
			scrollChatBottom('smooth');
		}, 0);
	}

	function onPromptKeydown(event: KeyboardEvent): void {
		if (['Enter'].includes(event.code)) {
			event.preventDefault();
			addMessage();
		}
	}

	// When DOM mounted, scroll to bottom
	onMount(() => {
		scrollChatBottom();
	});
</script>

<!-- 
<br />
<input type="text" bind:value={message} />
<button on:click={handleSendMessage}>Send</button>
<div>
	{#each messages as m, i}
		<li><pre>{JSON.stringify(m)}</pre></li>
	{/each}
</div> -->
<section class="card h-full">
	<div class="chat w-full h-full grid grid-cols-1 lg:grid-cols-[30%_1fr]">
		<!-- Navigation -->
		<div class="hidden lg:grid grid-rows-[auto_1fr_auto] border-r border-surface-500/30">
			<!-- Header -->
			<header class="border-b border-surface-500/30 p-4">
				<input class="input" type="search" placeholder="Search..." />
			</header>
			<!-- List -->
			<div class="p-4 space-y-4 overflow-y-auto">
				<small class="opacity-50">Contacts</small>
				<ListBox active="variant-filled-primary">
					{#each people as person, i}
						<ListBoxItem bind:group={currentPerson} name="people" value={person}>
							<svelte:fragment slot="lead">
								<Avatar src="https://i.pravatar.cc/?img={person.avatar}" width="w-8" />
							</svelte:fragment>
							{person.name}
						</ListBoxItem>
					{/each}
				</ListBox>
			</div>
			<!-- Footer -->
			<footer class="border-t border-surface-500/30 p-4">(footer)</footer>
		</div>
		<!-- Chat -->
		<div class="flex flex-col flex-nowrap justify-end max-h-screen">
			<!-- <div class="grid grid-row-[1fr_auto]"> -->
			<!-- Conversation -->
			<section bind:this={elemChat} class="max-w-fit p-4 overflow-y-auto space-y-4">
				{#each messageFeed as bubble, i}
					<Message
						avatar={bubble.avatar.toString()}
						name={bubble.name}
						timestamp={bubble.timestamp}
						message={bubble.message}
						color={bubble.color}
						self={bubble.host == true}
					/>
				{/each}
			</section>
			<!-- Prompt -->
			<section class="border-t border-surface-500/30 p-4">
				<div
					class="input-group input-group-divider grid-cols-[auto_1fr_auto] rounded-container-token"
				>
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
						on:click={addMessage}
					>
						<i class="fa-solid fa-paper-plane" />
					</button>
				</div>
			</section>
		</div>
	</div>
</section>

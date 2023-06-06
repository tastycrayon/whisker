import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
import { writable } from 'svelte/store';
import { ROOM_PATH } from './constant';
import { page } from '$app/stores';

export const messageStore = writable('');

export default {
	subscribe: messageStore.subscribe
};

import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
import { writable } from 'svelte/store';
import { ROOM_PATH } from './constant';
import { page } from '$app/stores';
import type { IMessage, IRoom, ISender } from './types';
import { pb } from './pocketbase';
import { toastStore, type ToastSettings } from '@skeletonlabs/skeleton';

export const messageStore = writable<IMessage | null>(null);

export default {
	subscribe: messageStore.subscribe
};

// rooms
type IRoomStore = { loading: boolean, error: any | undefined, data: IRoom[] | undefined }
export const roomStore = writable<IRoomStore>({ loading: false, error: undefined, data: undefined })
export const refreshRooms = async () => {
	try {
		roomStore.set({ loading: true, error: undefined, data: undefined })
		const data = await pb.send('/rooms', {
			method: 'GET'
		});
		if (!data) throw new Error("Failed loading rooms.")
		roomStore.set({ loading: false, error: undefined, data: data.sort((a: IRoom, b: IRoom) => (a.roomSlug > b.roomSlug)) })
	} catch (err: any) {
		const t: ToastSettings = {
			message: err.message, background: "variant-filled-error"
		};
		toastStore.trigger(t);
		roomStore.set({ loading: false, error: err, data: undefined })
	}
};

// participants
export const participantStore = writable<ISender[]>([]);
import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
import { writable } from 'svelte/store';
import { PARTICIPANT_PATH, ROOM_PATH } from './constant';
import { page } from '$app/stores';
import type { IMessage, IRoom, IParticipant, ResponseData } from './types';
import { pb } from './pocketbase';
import { toastStore, type ToastSettings } from '@skeletonlabs/skeleton';

export const messageStore = writable<IMessage | null>(null);

export default {
	subscribe: messageStore.subscribe
};

// rooms
export const roomStore = writable<ResponseData<IRoom[]>>({ loading: false, error: undefined, data: [] })

export const refreshRooms = async () => {
	try {
		roomStore.set({ loading: true, error: undefined, data: [] })
		const data = await pb.send<IRoom[]>(ROOM_PATH, {
			method: 'GET'
		});
		if (!data) throw new Error("Failed loading rooms.")
		roomStore.set({ loading: false, error: undefined, data: data.sort((a, b) => +(a.slug > b.slug)) })
	} catch (err: any) {
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
		roomStore.set({ loading: false, error: err, data: [] })
	}
};

// participants
export const participantStore = writable<ResponseData<IParticipant[]>>({ loading: false, error: undefined, data: [] });

export const refreshParticipants = async (slug: string) => {
	try {
		participantStore.set({ loading: true, error: undefined, data: [] })
		const data = await pb.send<IParticipant[]>(`${PARTICIPANT_PATH}?roomSlug=${slug}`, { method: 'GET' });
		if (!data) throw new Error("Failed loading participants.")
		participantStore.set({ loading: false, error: undefined, data: data })
	} catch (err: any) {
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
		roomStore.set({ loading: false, error: err, data: [] })
	}
};

export const refreshSingleParticipant = async (userId: string) => {
	try {
		participantStore.update((s) => ({ ...s, loading: true }))
		const res = await pb.send<IParticipant>(`${PARTICIPANT_PATH}/${userId}`, { method: 'GET' });
		if (!res) throw new Error("Failed loading participants.")
		participantStore.update((s) => {
			const idx = s.data.findIndex(e => e.id == userId)
			const participants = s.data
			if (idx == -1) {
				participants.push(res)
			} else {
				participants[idx] = res
			}
			return { error: undefined, data: participants, loading: false }
		})
	} catch (err: any) {
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
		roomStore.set({ loading: false, error: err, data: [] })
	}
};
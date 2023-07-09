import { writable, type Writable } from 'svelte/store';
import { DEFAULT_ROOM, PARTICIPANT_PATH, ROOM_PATH } from './constant';
import { page } from '$app/stores';
import { type IRecieveMessage, type IRoom, type IParticipant, type ResponseData, CollectionName } from './types';
import { pb } from './pocketbase';
import { toastStore, type ToastSettings, localStorageStore } from '@skeletonlabs/skeleton';
import type { ClientResponseError } from 'pocketbase';

export const messageStore = writable<IRecieveMessage | null>(null);

export default {
	subscribe: messageStore.subscribe
};

// rooms
export const roomStore = writable<ResponseData<IRoom[]>>({ loading: false, error: undefined, data: [] })

export const refreshRooms = async () => {
	try {
		roomStore.set({ loading: true, error: undefined, data: [] })
		const data = await pb.collection(CollectionName.Room).getFullList<IRoom>()

		if (!data) throw new Error("Failed loading rooms.")
		roomStore.set({ loading: false, error: undefined, data: data.sort((a, b) => +(a.slug > b.slug)) })
	} catch (err: any) {
		if ((err as ClientResponseError).isAbort) return
		roomStore.update(e => ({ loading: false, error: err, data: e.data }))
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
	}
};

export const deleteRoom = async (id: string) => {
	try {
		const result = await pb.collection('rooms').delete(id, {});
		if (!result) throw new Error("Failed deleting room.")
		roomStore.update((s) => {
			const rooms = [...s.data]
			const idx = rooms.findIndex(e => e.id == id)
			if (idx !== -1) rooms.splice(idx, 1)
			return { error: undefined, data: rooms, loading: false }
		})
	} catch (err: any) {
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
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
		if ((err as ClientResponseError).isAbort) return
		participantStore.update(e => ({ loading: false, error: err, data: e.data }))
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
	}
};

export const refreshSingleParticipant = async (userId: string) => {
	try {
		participantStore.update((s) => ({ ...s, loading: true }))
		const res = await pb.send<IParticipant>(`${PARTICIPANT_PATH}/${userId}`, { method: 'GET' });
		if (!res) throw new Error("Failed loading participants.")
		participantStore.update((s) => {
			const participants = [...s.data]
			const idx = participants.findIndex(e => e.id == userId)
			if (idx === -1) {
				participants.push(res)
			} else {
				participants[idx] = { ...res, roomSlug: participants[idx].roomSlug }
			}
			return { error: undefined, data: participants, loading: false }
		})
	} catch (err: any) {
		if ((err as ClientResponseError).isAbort) return
		participantStore.set({ loading: false, error: err, data: [] })
		const t: ToastSettings = { message: err.message, background: "variant-filled-error" };
		toastStore.trigger(t);
	}
};


// user detail
export const currentRoom: Writable<string> = localStorageStore('currentRoom', DEFAULT_ROOM);
interface UserPreference {
	lightMode: boolean
	hideWelcomeMessage: boolean
}
const defaultPreference: UserPreference = { lightMode: true, hideWelcomeMessage: true }
export const userPreference: Writable<UserPreference> = localStorageStore('userPreference', defaultPreference);
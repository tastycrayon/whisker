import { pb } from '$lib/pocketbase';

export const load = async () => {
	// const { data, error } = await pb.send('/rooms', {
	// 	method: 'GET'
	// });
	return { rooms: [] };
};

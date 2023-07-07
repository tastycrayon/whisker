
// import { goto } from "$app/navigation";
// import { ROOM_PATH } from "$lib/constant";
// import type { HandleClientError, } from "@sveltejs/kit";


// export const handle: HandleClientError = () => {
//     console.log("hooks.client.ts")
//     goto(ROOM_PATH)
// }
import { COOKIE_OPTIONS } from '$lib/constant';
import { currentUser, pb } from '$lib/pocketbase';

try {
    pb.authStore.loadFromCookie(document.cookie); // load it first
    // get an up-to-date auth store state by veryfing and refreshing the loaded auth model (if any)
    pb.authStore.isValid && (await pb.collection('users').authRefresh());
} catch (_) {
    // clear the auth store on failed refresh
    pb.authStore.clear();
}
// pb.authStore.loadFromCookie(document.cookie); // load it first
pb.authStore.onChange((_, model) => {
    if (pb.authStore.isValid) {
        currentUser.set(model);
        document.cookie = pb.authStore.exportToCookie(COOKIE_OPTIONS);
    }
});
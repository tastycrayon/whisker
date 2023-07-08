import { PUBLIC_WEBSOCKET_URL } from "$env/static/public";
import type { Unsubscriber, Writable } from "svelte/store";
import { ROOM_PATH } from "./constant";
import { currentRoom } from "./store";
import type { ISendMessage } from "./types";


export const WebSocketStore = (reader: Writable<string | null>, writer: Writable<ISendMessage | null>) => {
    const reopenTimeouts = [2000, 5000, 10000, 30000, 60000];
    let retryCount = 0, timer: number | undefined;
    function retryTimeout() {
        retryCount++;
        if (retryCount > reopenTimeouts.length - 1) retryCount--;
        return reopenTimeouts[retryCount];
    }

    let socket: WebSocket | undefined
    let unsubscribeWriter: Unsubscriber | undefined;

    function close() {
        if (timer) clearTimeout(timer);
        if (!socket) return
        socket.close();
        socket = undefined;
        if (unsubscribeWriter) unsubscribeWriter()
    }
    function reopen() {
        close();
        timer = setTimeout(() => open(), retryTimeout());
    }
    let cRoom: string = ""
    currentRoom.subscribe(room => cRoom = room)

    let openPromise: Promise<WebSocket | undefined>
    const open = (): Promise<WebSocket | undefined> => {
        return new Promise(function (resolve, reject) {
            if (socket) return resolve(socket);
            socket = new WebSocket(`${PUBLIC_WEBSOCKET_URL}${ROOM_PATH}/${cRoom}`); // open websocket
            // socket open
            // writer writes to websocket (send messages)
            unsubscribeWriter = writer.subscribe((m: ISendMessage | null) => {
                if (!m || !socket || socket.readyState !== WebSocket.OPEN) return
                socket.send(JSON.stringify(m));
                writer.set(null)
            });
            // socket open
            socket.onopen = () => {
                retryCount = 0;
                return resolve(socket);
            };
            socket.onerror = (err) => reject(err);
            socket.onclose = () => reopen()
            socket.onmessage = (event) => reader.set(event.data)
        });
    }

    async function connect() {
        if (timer) clearTimeout(timer)
        timer = undefined;

        if (openPromise) return openPromise
        openPromise = open()
        return openPromise
    }
    return {
        connect, open, reopen, socket
    }
}

// /**
//  * Create a writable store based on a web-socket.
//  * Data is transferred as JSON.
//  * Keeps socket open (reopens if closed) as long as there are subscriptions.
//  * @param {string} url the WebSocket url
//  * @param {any} initialValue store value used before 1st. response from server is present
//  * @param {string[]} socketOptions transparently passed to the WebSocket constructor
//  * @return {Store}
//  */
// export function websocketStore(url: string, initialValue: any, socketOptions: string[]) {
//     let socket: WebSocket | undefined, openPromise: Promise<unknown> | undefined, reopenTimeoutHandler: number | undefined;
//     let reopenCount = 0;
//     const subscriptions = new Set<(x: string) => any>();

//     function reopenTimeout() {
//         const n = reopenCount;
//         reopenCount++;
//         return reopenTimeouts[
//             n >= reopenTimeouts.length - 1 ? reopenTimeouts.length - 1 : n
//         ];
//     }

//     function close() {
//         if (reopenTimeoutHandler) {
//             clearTimeout(reopenTimeoutHandler);
//         }

//         if (socket) {
//             socket.close();
//             socket = undefined;
//         }
//     }

//     function reopen() {
//         close();
//         if (subscriptions.size > 0) {
//             reopenTimeoutHandler = setTimeout(() => open(), reopenTimeout());
//         }
//     }

//     async function open() {
//         if (reopenTimeoutHandler) {
//             clearTimeout(reopenTimeoutHandler);
//             reopenTimeoutHandler = undefined;
//         }

//         // we are still in the opening phase
//         if (openPromise) {
//             return openPromise;
//         }

//         socket = new WebSocket(url, socketOptions);

//         socket.onmessage = event => {
//             initialValue = JSON.parse(event.data);
//             subscriptions.forEach(subscription => subscription(initialValue));
//         };

//         socket.onclose = event => reopen();

//         openPromise = new Promise((resolve, reject) => {
//             if (!socket) return
//             socket.onerror = error => {
//                 reject(error);
//                 openPromise = undefined;
//             };
//             socket.onopen = event => {
//                 reopenCount = 0;
//                 resolve(1);
//                 openPromise = undefined;
//             };
//         });
//         return openPromise;
//     }

//     return {
//         set: (value: string) => {
//             const send = () => socket?.send(JSON.stringify(value));
//             if (socket?.readyState !== WebSocket.OPEN) open().then(send);
//             else send();
//         },
//         subscribe: (subscription: (x: string) => void) => {
//             open();
//             subscription(initialValue);
//             subscriptions.add(subscription);
//             return () => {
//                 subscriptions.delete(subscription);
//                 if (subscriptions.size === 0) {
//                     close();
//                 }
//             };
//         }
//     };
// }

// export default websocketStore;
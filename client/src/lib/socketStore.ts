
import { writable, type Unsubscriber, type Writable, } from "svelte/store";
import { NOT_FOUND_PATH } from "./constant";
import type { IRecieveMessage, SendEvent, TextMessage } from "./types";
import { goto } from "$app/navigation";
import { toastStore, type ToastSettings } from "@skeletonlabs/skeleton";

// writer writes to websocket (send messages)
// reader reads from websocket
export const reader = writable<IRecieveMessage | null>(null);
export const writer = writable<SendEvent | null>(null);
export type loadingType = "switching" | "loading" | null
export const loading = writable<loadingType>(null);

const createWebSocket = (url: string): WebSocket => new WebSocket(url)
export class WS {
    _reader: Writable<IRecieveMessage | null>
    _writer: Writable<SendEvent | null>
    _loading: Writable<"switching" | "loading" | null>
    _unsubscribeWriter: Unsubscriber | undefined
    _url: string = ""
    _socket: WebSocket | undefined
    _promise: Promise<WebSocket> | undefined
    _reopenTimeouts = [2000, 5000, 10000, 30000, 60000];
    _retryCount: number = 0;
    _timer: number = 0;
    messageFeed: TextMessage[] = [];


    retryTimeout(): number {
        this._retryCount++;
        if (this._retryCount >= this._reopenTimeouts.length) this._retryCount--;
        return this._reopenTimeouts[this._retryCount];
    }

    get isOpening() {
        return Boolean(this._socket && this._socket.readyState === WebSocket.CONNECTING);
    }
    get isOpened() {
        return Boolean(this._socket && this._socket.readyState === WebSocket.OPEN);
    }
    get isClosing() {
        return Boolean(this._socket && this._socket.readyState === WebSocket.CLOSING);
    }
    get isClosed() {
        return Boolean(!this._socket || this._socket.readyState === WebSocket.CLOSED);
    }
    _onOpen(event: Event) {
        this.setLoadingFalse()
        this._retryCount = 0
    }
    _onError(event: Event) {
        this.setLoadingFalse()
        if (!navigator.onLine) return
        this._notifyError("Error establishing connection")
    }

    _onClose(event: CloseEvent) {
        this.setLoadingFalse()
        switch (event.code) {
            case 1000: // StatusNormalClosure
                if (event.reason === "GOING_AWAY") return
                if (event.reason === "404_NOT_FOUND") goto(NOT_FOUND_PATH)
                break;
            default:
                const message = event.reason === "" ? "Failed to establish websocket connection" : event.reason
                this._notifyError(message)
                this._reOpen()
                break;
        }
    }
    _onMessage(event: MessageEvent<string>) {
        console.log("read")
        this._reader.set(JSON.parse(event.data) as IRecieveMessage)
    }

    _reOpen() {
        this.enableLoading()
        this.closeConnection()
        this._timer = setTimeout(this.connect.bind(this), this.retryTimeout());
    }
    _notifyError(message: string) {
        const t: ToastSettings = { message, background: "variant-filled-surface" };
        setTimeout(() => toastStore.trigger(t), 100)
    }
    _notifyOpening(message: string) {
        const t: ToastSettings = { message, background: "variant-filled-warning" };
        toastStore.trigger(t)
    }
    enableLoading() {
        this._loading.set("loading")
    }
    enableSwitching() {
        this._loading.set("switching")
    }
    setLoadingFalse() {
        this._loading.set(null)
    }
    resetSubscriber() {
        if (!this._unsubscribeWriter) return;
        this._unsubscribeWriter()
        this._unsubscribeWriter = undefined
    }
    setWriter() {
        this._unsubscribeWriter = writer.subscribe((m: SendEvent | null) => {
            if (!m || !this._socket || !this.isOpened) return
            console.log("write")
            this._socket.send(JSON.stringify(m));
            writer.set(null)
        });
    }
    open(): Promise<WebSocket> {
        return new Promise<WebSocket>((res, rej) => {
            if (this._socket && this.isOpened) return res(this._socket)
            if (this._socket && (this.isClosing || this.isClosed)) return rej(this._socket)
            // if (this.isOpening) return // TODO

            this._socket = createWebSocket(this._url)
            console.log("connected")
            this._socket.onopen = this._onOpen.bind(this)
            this._socket.onclose = this._onClose.bind(this)
            this._socket.onerror = this._onError.bind(this)
            this._socket.onmessage = this._onMessage.bind(this)

            this.setWriter()
            res(this._socket)
        })
    };
    setUrl(url: string) {
        this._url = url;
        return this
    }
    connect(isSwitching = false): Promise<WebSocket | undefined> {
        if (this._promise) return this._promise
        if (this._timer) clearTimeout(this._timer)
        // switch otherwise enable loading indicator
        if (isSwitching) this.enableSwitching()
        else this.enableLoading()

        this._timer = 0;
        this._promise = this.open()
        return this._promise

    }
    closeConnection(): Promise<boolean> {
        return new Promise<boolean>((res, rej) => {
            console.log("about to close")
            if (this._timer) clearTimeout(this._timer);
            this.resetSubscriber()
            this.messageFeed = []
            if (!this._socket) return res(true)
            if (this.isOpened) this._socket.close(1000, "GOING_AWAY")
            this._socket = undefined;
            this._promise = undefined;
            res(true)
        })
    }
    constructor(reader: Writable<IRecieveMessage | null>, writer: Writable<SendEvent | null>, loading: Writable<loadingType>) {
        this._reader = reader
        this._writer = writer
        this._loading = loading
    }
}

export const socket = new WS(reader, writer, loading);
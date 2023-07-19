
import { writable, type Unsubscriber, type Writable, } from "svelte/store";
import { NOT_FOUND_PATH } from "./constant";
import type { IRecieveMessage, SendEvent } from "./types";
import { goto } from "$app/navigation";
import { toastStore, type ToastSettings } from "@skeletonlabs/skeleton";

// writer writes to websocket (send messages)
// reader reads from websocket
export const reader = writable<IRecieveMessage | null>(null);
export const writer = writable<SendEvent | null>(null);
export const loading = writable<boolean>(false);

const createWebSocket = (url: string): WebSocket => new WebSocket(url)
export class WS {
    _reader: Writable<IRecieveMessage | null>
    _writer: Writable<SendEvent | null>
    _loading: Writable<boolean>
    _unsubscribeWriter: Unsubscriber | undefined
    _url: string
    _socket: WebSocket | undefined
    _promise: Promise<WebSocket> | undefined
    _reopenTimeouts = [2000, 5000, 10000, 30000, 60000];
    _retryCount: number = 0;
    _timer: number = 0;

    get retryTimeout(): number {
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
        this._notifyError("Error establishing connection")
    }
    _onSend(event: MessageEvent<string>) {
        console.log("read")
        reader.set(JSON.parse(event.data) as IRecieveMessage)
    }
    _onClose(event: CloseEvent) {
        this.setLoadingFalse()
        switch (event.code) {
            case 1000: // StatusNormalClosure
                if (event.reason === "GOING_AWAY") return
                if (event.reason === "404_NOT_FOUND") goto(NOT_FOUND_PATH)
                break;
            default:
                const message = event.reason === "" ? "Failed to establish connection" : event.reason
                this._notifyError(message)
                this._reOpen()
                break;
        }
    }
    _onMessage(event: MessageEvent<string>) { reader.set(JSON.parse(event.data) as IRecieveMessage) }

    _reOpen() {
        this.setLoadingTrue()
        this.closeConnection()
        this._timer = setTimeout(this.connect.bind(this), this.retryTimeout);
    }
    _notifyError(message: string) {
        const t: ToastSettings = { message, background: "variant-filled-surface" };
        setTimeout(() => toastStore.trigger(t), 100)
    }
    setLoadingTrue() {
        this._loading.set(true)
    }
    setLoadingFalse() {
        this._loading.set(false)
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
            // if (this.isOpening) return // umm.. what should I do here? await this.isOpening ?

            if (this._unsubscribeWriter) this._unsubscribeWriter()
            this._socket = createWebSocket(this._url)
            this._socket.onopen = this._onOpen.bind(this)
            this._socket.onclose = this._onClose.bind(this)
            this._socket.onerror = this._onError.bind(this)
            this._socket.onmessage = this._onMessage.bind(this)

            this.setWriter()
            res(this._socket)
        })
    };
    connect(): Promise<WebSocket | undefined> {
        if (this._promise) return this._promise
        if (this._timer) clearTimeout(this._timer)
        this.setLoadingTrue()
        this._timer = 0;
        this._promise = this.open()
        return this._promise

    }
    closeConnection() {
        console.log("about to close")
        if (this._timer) clearTimeout(this._timer);
        if (this._unsubscribeWriter) this._unsubscribeWriter()
        if (!this._socket) return
        if (this.isOpened) this._socket.close(1000, "GOING_AWAY")
        this._socket = undefined;
        this._promise = undefined;
    }
    constructor(url: string, reader: Writable<IRecieveMessage | null>, writer: Writable<SendEvent | null>, loading: Writable<boolean>) {
        this._url = url;
        this._reader = reader
        this._writer = writer
        this._loading = loading
    }
}
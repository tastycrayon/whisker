import { PUBLIC_POCKETBASE_URL } from "$env/static/public";
import { modalStore, type ModalSettings } from "@skeletonlabs/skeleton";
import type { CollectionName } from "./types";

export function timeSince(time: number | string) {
    switch (typeof time) {
        case 'number':
            break;
        case 'string':
            time = +new Date(time);
            break;
        default:
            time = +new Date();
    }
    const localOptopns: Intl.DateTimeFormatOptions = {
        timeStyle: "short"
    }
    const today = new Date()
    const givenTime = new Date(time)
    if (today.getFullYear() === givenTime.getFullYear() &&
        today.getMonth() === givenTime.getMonth() &&
        today.getDate() === givenTime.getDate())
        return `Today at ${givenTime.toLocaleTimeString([], localOptopns)}`

    const yesterday = new Date()
    yesterday.setDate(yesterday.getDate() - 1)
    if (yesterday.getFullYear() === givenTime.getFullYear() &&
        yesterday.getMonth() === givenTime.getMonth() &&
        yesterday.getDate() === givenTime.getDate())
        return `Yesterday at ${givenTime.toLocaleTimeString([], localOptopns)}`

    return givenTime.toLocaleDateString() + ' ' + givenTime.toLocaleTimeString([], localOptopns);
}


export const generateAvatar = (coll: CollectionName, id: string, avatar: string) => {
    if (id === "" || avatar == "") return "/default.jpg"
    return `${PUBLIC_POCKETBASE_URL}/api/files/${coll}/${id}/${avatar}`
}


export const convertToSlug = (e: string): string => {
    return e
        .toLowerCase()
        .replace(/ /g, '-')
        .replace(/[^\w-]+/g, '');
};


export const firstConfirm = (): Promise<boolean> => {
    return new Promise<boolean>((resolve) => {
        const modal: ModalSettings = {
            type: 'confirm',
            title: 'Please Confirm',
            body: 'Are you sure you wish to proceed?',
            response: (r: boolean) => {
                resolve(r);
            }
        };
        modalStore.trigger(modal);
    })
};

export const isImageUrl = (url: string): boolean => {
    const pattern = /^https?:\/\/.+\.(jpg|jpeg|png|webp|avif|gif)$/;
    return pattern.test(url)
}
export const isVideoUrl = (url: string): boolean => {
    const pattern = /^https?:\/\/.+\.(mp4|ogg|webm)$/;
    return pattern.test(url)
}
export const isYoutubeUrl = (url: string): boolean => {
    const pattern = /^(https?\:\/\/)?(www\.youtube\.com|youtu\.be)\/.+$/;
    return pattern.test(url)
}

export const getYoutubeVideoId = (url: string): string | null => {
    let pattern = /(youtu.*be.*)\/(watch\?v=|embed\/|v|shorts|)(.*?((?=[&#?])|$))/gm;
    const ids = pattern.exec(url);
    if (!ids) return null
    return ids[3] || null
}
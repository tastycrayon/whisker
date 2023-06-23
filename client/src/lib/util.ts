import { PUBLIC_POCKETBASE_URL } from "$env/static/public";

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


export const generateUserAvatar = (id: string, avatar: string) =>
    `${PUBLIC_POCKETBASE_URL}/api/files/users/${id}/${avatar}`

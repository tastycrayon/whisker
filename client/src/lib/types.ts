// collection
export enum CollectionName {
	Room = 'rooms',
	User = 'users'
}

export interface ResponseData<T> {
	loading: boolean;
	error: any | undefined;
	data: T
}

// participant
export interface IUser {
	avatar: string
	collectionId: string
	collectionName: string
	created: string
	email: string
	emailVisibility: boolean
	id: string
	name: string
	room: string
	updated: string
	username: string
	verified: boolean
}

// message
export enum MessageType {
	Welcome = 2,
	Ping = 4, // burst 
	Text = 8,
	Bailout = 16,
	Swap = 32,
	History = 64 // burst
}
interface PingType {
	sid: string;
	messageType: MessageType.Ping;
	content: IParticipant[]
	sender: IParticipant;
	roomSlug: string;
	created: string;
}
interface HistoryType {
	sid: string;
	messageType: MessageType.History;
	content: IMessage[]
	sender: IParticipant;
	roomSlug: string;
	created: string;
}
interface NormalType {
	sid: string;
	messageType: MessageType.Welcome | MessageType.Bailout | MessageType.Text | MessageType.Swap;
	content: string;
	sender: IParticipant;
	roomSlug: string;
	created: string;
}

export type IMessage = PingType | HistoryType | NormalType
export interface IParticipant {
	id: string;
	username: string;
	roomSlug: string;
	avatar: string;
}

// room
export enum RoomType {
	PublicRoom = "public",
	PersonalRoom = "personal"
}

export interface IRoom {
	id: string;
	cover: string;
	slug: string
	name: string
	type: string;
	description: string
	createdBy: string
	created: string
}
export interface ParticipantRoom {
	id: string;
	name: string;
	avatar: string;
}

// error handling
export interface FormResError<T> {
	code: T | 'unknown';
	message: string;
}
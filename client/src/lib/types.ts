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
// event 
export enum EventType {
	Info = 1,
	Text = 2,
	Swap = 4,
	ParticipantHistory = 6, // slice
	MessageHistory = 8, // slice
}
export interface SendEvent {
	payload: string,
	type: EventType
}
// message
export interface SwapMessage {
	from: string;
	to: string
}
export interface TextMessage {
	id: string;
	content: string
	sender: IParticipant
	created: string
}

interface MessageHistoryType {
	payload: TextMessageWrap[],
	slug: string,
	type: EventType.MessageHistory
}
interface ParticipantHistoryType {
	payload: IParticipant[],
	slug: string,
	type: EventType.ParticipantHistory
}
interface TextMessageWrap {
	payload: TextMessage,
	slug: string,
	type: EventType.Text | EventType.Info
}

export type IRecieveMessage = MessageHistoryType | ParticipantHistoryType | TextMessageWrap

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
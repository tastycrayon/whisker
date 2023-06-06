package ws

import "github.com/tastycrayon/go-chat/app/queue"

type Room struct {
	RoomId     string                  `json:"roomId"`
	RoomName   string                  `json:"roomName"`
	Clients    map[string]*Participant `json:"clients"`
	CreatedBy  string                  `json:"createdBy"`
	PostsQueue queue.CQueue[Post]
}

package ws

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

//	type Message struct {
//		Message  string `json:"message"`
//		ClientId string `json:"clientId"`
//		RoomId   string `json:"roomId"`
//		Username string `json:"username"`
//	}
type MessageType string

const (
	Welcome MessageType = "welcome"
	Ping    MessageType = "ping"
	Text    MessageType = "text"
	Bailout MessageType = "bailout"
)

type MessageReq struct {
	Content string      `json:"content"`
	Type    MessageType `json:"messageType"`
}

type Message struct {
	Id      string      `json:"id"`
	Content string      `json:"content"`
	Type    MessageType `json:"messageType"`
}

type MessageSender struct {
	ClientId string `json:"clientId"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
	Avatar   string `json:"avatar"`
}

type Post struct {
	Message Message `json:"message"`
	Sender  MessageSender
	RoomId  string    `json:"roomId"`
	Created time.Time `json:"created"`
}

func GenerateId() string {
	str := uuid.NewString()
	str = strings.ReplaceAll(str, "-", "")

	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
	}
	// The slice should now contain random bytes instead of only zeroes.
	return str[:15]
}

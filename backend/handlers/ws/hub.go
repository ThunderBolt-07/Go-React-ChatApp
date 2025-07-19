package ws

import (
	"chat/handlers/types"
	"encoding/json"
	"log"
)

type Hub struct {
	ClientMap     map[string]*Client
	Broadcast     chan []byte
	Register      chan string
	Unregister    chan string
	BroadcastFile chan *types.UrlFile
}

func NewHub() *Hub {
	return &Hub{
		ClientMap:  make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan string),
		Unregister: make(chan string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case name := <-h.Unregister:
			log.Println("client unregistered")
			_, ok := h.ClientMap[name]
			if ok {
				delete(h.ClientMap, name)
			}
		// case name := <-h.Register:
		// 	log.Println("client registered")
		// 	_, ok := h.ClientMap[c]
		// 	if !ok {
		// 		h.ClientMap[c] = true
		// 	}
		case msg := <-h.Broadcast:
			var um SendMessage
			err := json.Unmarshal(msg, &um)
			if err != nil {
				log.Println("erro while decoding byte to json , stopping hub", err)
				break
			}
			log.Println("received in Broadcast", string(msg), um)
			for name, client := range h.ClientMap {

				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.ClientMap, name)
					log.Println("unable to braodcast too client ", client)

				}

			}
		}
	}
}

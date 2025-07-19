package ws

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Name string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

type SendMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// type Entry struct{
// 	Name string
// 	client *Client
// }

const (
	writeDeadline = time.Second * 10
	readDeadline  = time.Second * 60
)

var ChatClient *Client

func HandleWS(w http.ResponseWriter, r *http.Request, H *Hub) error {
	userName := r.URL.Query().Get("username")

	upgrader.CheckOrigin = func(r *http.Request) bool {

		o := os.Getenv("ALLOWED_FRONTEND_ORIGIN")
		log.Println(r.Header["Origin"][0], o)
		if o == "" {
			log.Fatal("uable to fetch allows frontend origin form env")
		}
		log.Println(r.Header["Origin"][0], o)
		return r.Header["Origin"][0] == o
	}
	// origin := r.Header["Origin"]

	// u, _ := url.Parse(origin[0])

	// log.Println(r.Header["Origin"], u.Host, r.Host, origin)
	_, ok := H.ClientMap[userName]
	if ok {
		log.Println("username alrady exists in map")
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("user with same name already exists"))
		return errors.New("username alrady exists")

	}
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error occured while upgrading Connection", err.Error())
		return err
	}
	send := make(chan []byte)
	client := Client{userName, H, Conn, send}
	ChatClient = &client
	H.ClientMap[userName] = &client

	go client.readPump()
	go client.writePump()
	return nil
}

// read from client and write to Hub
func (c *Client) readPump() {

	defer func() {
		log.Println("client closing from read pump ", c)
		c.Hub.Unregister <- c.Name
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(1024)
	c.Conn.SetReadDeadline(time.Now().Add(readDeadline))
	c.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
	c.Conn.SetPongHandler(func(appData string) error {

		err := c.Conn.SetReadDeadline(time.Now().Add(readDeadline))

		if err != nil {
			log.Println("error in pong handler", err)
			return err
		}
		return nil

	})

	for {

		_, p, err := c.Conn.ReadMessage()
		log.Println("rec client with name ", c.Name, string(p))
		message := SendMessage{c.Name, string(p)}
		if err != nil {
			log.Println("err in receiving message in readPump", err)
			break
		}
		bm, err := json.Marshal(message)
		if err != nil {
			log.Println("error while converting p to json byte", err)
			break
		}
		c.Hub.Broadcast <- bm
		log.Println("broadcasted", string(p), err)
	}
}

// read from Hub and write to client
func (c *Client) writePump() {
	defer func() {
		log.Println("client closing from write pump ", c)
		c.Hub.Unregister <- c.Name
		c.Conn.Close()
	}()

	ticker := time.NewTicker(time.Second * 30)
outer:
	for {
		select {
		case msg := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			// err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			// if err!=nil {
			// 	log.Println("error sending message to clined in writePump",err)
			// 	break outer
			// }
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("error sending message to clined in writePump", err)
				break outer
			}

			br, err := w.Write(msg)
			log.Println("bytes written and error ", br, err)
			n := len(c.Send)
			//log.Println("multiple messages ququed sending them ", n)

			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			//log.Println("ticked at ", tick)
			c.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("error while sending ping, closing connection", c, err)
				return
			}
		}
	}

}

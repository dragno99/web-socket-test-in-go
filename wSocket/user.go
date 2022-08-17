package wSocket

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type User struct {
	// UserId
	UserId string
	// room in which user reside
	Room *Room

	// websocket connection fot the user
	Conn *websocket.Conn

	// send channel for the user to recieve messages from the websocket
	Send chan Message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// readPump : pumps message from the websocket connection to the room
func (u *User) readPump() {

	for {
		var msg Message
		err := u.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			u.Room.unregister <- u
			break
		}
		u.Room.broadcast <- msg
	}
}

// writePump : pumps message from the hub to the websocket connetion
func (u *User) writePump() {

	for {
		msg, ok := <-u.Send
		if !ok {
			// the room closed the channel
			u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			u.Room.unregister <- u
			break
		}
		// write message to websocket connection
		err := u.Conn.WriteJSON(msg)

		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			u.Room.unregister <- u
			break
		}
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

// serveWS handles websocket request from the peer
func ServeWS(room *Room, userId string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	user := &User{UserId: userId, Room: room, Conn: conn, Send: make(chan Message)}

	room.register <- user

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go user.writePump()
	go user.readPump()

}

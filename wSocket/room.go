package wSocket

type Room struct {
	// chatId
	ChatId string
	
	// collection of users(websocket connections) in this room
	Users map[*User]bool

	// broadcast channel to send message to all of the users in this room
	broadcast chan Message

	// channel to add/register a new user
	register chan *User

	// channel to remove/unregister a user
	unregister chan *User

	// reason for defining register and unregister channel to avoid race condition
}

func NewRoom(chatId string) *Room {
	return &Room{
		Users:      make(map[*User]bool),
		broadcast:  make(chan Message),
		register:   make(chan *User),
		unregister: make(chan *User),
		ChatId:     chatId,
	}
}

func (r *Room) Run() {
	for {
		select {

		case user := <-r.register:
			{
				r.Users[user] = true
			}
		case user := <-r.unregister:
			{
				user.Conn.Close()
				if _, ok := r.Users[user]; ok {
					delete(r.Users, user)
					close(user.Send)
				}
			}

		case msg := <-r.broadcast:
			{
				for otherUser := range r.Users {
					// user should not send message to himself
					if otherUser.UserId == msg.SenderId {
						continue
					}
					select {

					case otherUser.Send <- msg:

					default:
						{
							close(otherUser.Send)
							delete(r.Users, otherUser)
						}
					}
				}
			}
		}
	}
}

package wSocket

type Message struct {
	Text      string `json:"text"`
	Image     string `json:"image"`
	SenderId  string `json:"senderId"`
	CreatedAt string `json:"createdAt"`
}

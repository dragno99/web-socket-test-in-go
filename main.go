package main

import (
	"encoding/json"
	"log"
	"myWebSocket/wSocket"
	"net/http"

	"github.com/gorilla/mux"
)

var manyRooms = make(map[string]*wSocket.Room)

func joinRooms(w http.ResponseWriter, r *http.Request) {
	chatId := mux.Vars(r)["roomId"]
	userId := mux.Vars(r)["userId"]

	if chatId == "" || userId == "" {
		json.NewEncoder(w).Encode("please provide chatId/userId")
		return
	}

	_, ok := manyRooms[chatId]

	if !ok {
		manyRooms[chatId] = wSocket.NewRoom(chatId)
		go manyRooms[chatId].Run()
	}

	wSocket.ServeWS(manyRooms[chatId], userId, w, r)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/start/{roomId}/{userId}", http.HandlerFunc(joinRooms))

	log.Fatal(http.ListenAndServe(":8080", router))

}

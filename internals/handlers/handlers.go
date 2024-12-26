package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var viewer = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var WsChan = make(chan WsPayload)

var clients = make(map[websocketConnection]string)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsResponse struct {
	Message     string `json:"message"`
	Action      string `json:"action"`
	MessageType string `json:"message_type"`
}

type websocketConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Message  string              `json:"message"`
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Conn     websocketConnection `json:"-"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	// ...
	err := renderPage(w, "/home/chandra/Peer2PeerChat/cmd/web/html/home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// ...

	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client connected to endpoint")

	var response WsResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := websocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)

}

func ListenForWs(conn *websocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error: ", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
		} else {
			payload.Conn = *conn
			WsChan <- payload
		}
	}
}

func ListenForWsChannel() {
	var response WsResponse
	for {
		e := <-WsChan
		_ = e
		response.Message = "Welcome to the server"
		response.Action = "New Action"
		BroadCastToAll(&response)
	}
}

func BroadCastToAll(response *WsResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println(err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmp string, data jet.VarMap) error {
	view, err := viewer.GetTemplate(tmp)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

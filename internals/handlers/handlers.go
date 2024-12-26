package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var viewer = jet.NewSet(
	jet.NewOSFileSystemLoader("/home/chandra/Peer2PeerChat/cmd/web/html"),
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
	Message        string   `json:"message"`
	Action         string   `json:"action"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
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
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// WsEndpoint is a websocket endpoint
// It upgrades the HTTP server connection to the WebSocket protocol
// It sets up a websocket connection and listens for incoming messages
func WsEndpoint(w http.ResponseWriter, r *http.Request) {

	// Creatng a websocket connection(*websocket.Conn)
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
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
		// response.Message = "Welcome to the server"
		// response.Action = "New Action"
		// BroadCastToAll(&response)

		switch e.Action {
		case "username":
			// get a list of all users and send it back via broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			BroadCastToAll(&response)

		case "left":
			delete(clients, e.Conn)
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			BroadCastToAll(&response)

		case "broadcast":
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			response.Action = "broadcast"
			BroadCastToAll(&response)
		}
	}
}

func getUserList() []string {
	var userList []string
	for client := range clients {
		if clients[client] == "" {
			continue
		}
		userList = append(userList, clients[client])
	}
	sort.Strings(userList)
	return userList
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

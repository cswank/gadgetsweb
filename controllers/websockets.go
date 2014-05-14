package controllers

import (
	"time"
	"fmt"
	"log"
	"net/http"
	"bitbucket.org/cswank/gogadgets/models"
	"github.com/gorilla/websocket"
	"github.com/vaughan0/go-zmq"
	"encoding/json"
)

var (
	pingMsg = [][]byte{
		[]byte("ping"),
		[]byte(""),
	}
)

//InSocket (from the client's point of view) is used to send messages
//from server to client
func HandleInSocket(w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	host := params["host"][0]
	ctx, err := zmq.NewContext()
	defer ctx.Close()
	if err != nil {
		return err
	}
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return err
	}
	quitSub := make(chan bool)
	sockIsDone := make(chan bool)
	go getZMQMessage(conn, ctx, host, quitSub)
	<-sockIsDone
	quitSub <- true
	return nil
}

//InSocket (from the client's point of view) is used to send messages
//from client to server
func HandleOutSocket(w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	host := params["host"][0]
	ctx, err := zmq.NewContext()
	defer ctx.Close()
	if err != nil {
		return err
	}
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return err
	}
	quitSub := make(chan bool)
	sockIsDone := make(chan bool)
	go getSocketMessage(conn, ctx, host, sockIsDone)
	<-sockIsDone
	quitSub <- true
	return nil
}

//When a message is received from the zmq socket, it is passed along to the web socket.
func getZMQMessage(conn *websocket.Conn, ctx *zmq.Context, host string, shouldQuit <-chan bool) error {
	sub, chans, err := getSubChannels(ctx, host)
	if err != nil {
		return err
	}
	defer sub.Close()
	defer chans.Close()
	
	for {
		select {
		case msg := <-chans.In():
			sendSocketMessage(conn, msg)
		case <-shouldQuit:
			return nil
		case <-time.After(15 * time.Second):
			sendSocketMessage(conn, pingMsg)
		case err := <-chans.Errors():
			log.Println("get sub err", err)
			return err
		}
	}
	return nil
}

//When a message is received from the web socket it is passed along to the zmq socket.
func getSocketMessage(conn *websocket.Conn, ctx *zmq.Context, host string, done chan<- bool) error {
	pub, err := ctx.Socket(zmq.Pub)
	if err = pub.Connect(fmt.Sprintf("tcp://%s:6111", host)); err != nil {
		return err
	}
	defer pub.Close()
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	requestStatus(pub)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("get sock err", err)
			done <- true
			return err
		}
		if messageType == websocket.TextMessage {
			sendZMQMessage(p, pub)
		} else if messageType == websocket.CloseMessage || messageType == -1 {
			done <- true
			return nil
		}
	}
	return nil
}

//Send a message via the zmq socket.
func sendZMQMessage(input []byte, pub *zmq.Socket) {
	cmd := &command{}
	err := json.Unmarshal(input, cmd)
	if err != nil {
		log.Println(err)
		return
	}
	b, _ := json.Marshal(cmd.Message)
	msg := [][]byte{
		[]byte(cmd.Event),
		b,
	}
	pub.Send(msg)
}

//Send a message via the web socket.
func sendSocketMessage(conn * websocket.Conn, message [][]byte) {
	payload := []string {
		string(message[0]),
		string(message[1]),
	}
	b, _ := json.Marshal(payload)
	conn.WriteMessage(websocket.TextMessage, b)
}

func requestStatus(pub *zmq.Socket) {
	msg := models.Message{
		Type: models.COMMAND,
		Body: "update",
        }
	b, _ := json.Marshal(&msg)
        pub.Send([][]byte{
		[]byte(msg.Type),
		b,
	})
}

func getSubChannels(ctx *zmq.Context, host string) (sub *zmq.Socket, chans *zmq.Channels, err error) {
	//uid := "gadgets ctrl";
	sub, err = ctx.Socket(zmq.Sub)
	if err != nil {
		return sub, chans, err
	}
	if err = sub.Connect(fmt.Sprintf("tcp://%s:6112", host)); err != nil {
		return sub, chans, err
	}
	sub.Subscribe([]byte("update"))
	sub.Subscribe([]byte("method"))
	sub.Subscribe([]byte("info"))
	chans = sub.Channels()
	return sub, chans, err
}

type command struct {
	Event string
	Message map[string]interface{}
}

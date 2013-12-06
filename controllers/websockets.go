package controllers

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/vaughan0/go-zmq"
	"encoding/json"
)


func HandleSocket(w http.ResponseWriter, r *http.Request) error {
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
	go getSockMessage(conn, ctx, sockIsDone)
	go getSubMessage(conn, ctx, quitSub)
	<-sockIsDone
	quitSub <- true
	fmt.Println("finished")
	return nil
}

func sendMessage(conn * websocket.Conn, message [][]byte) {
	payload := []string {
		string(message[0]),
		string(message[1]),
	}
	b, _ := json.Marshal(payload)
	conn.WriteMessage(websocket.TextMessage, b)
}

func requestStatus(conn * websocket.Conn, ctx *zmq.Context) error {
	req, err := ctx.Socket(zmq.Req)
	defer req.Close()
	if err = req.Connect("tcp://192.168.1.16:6113"); err != nil {
		return err
	}
	msg := [][]byte{[]byte("status"), []byte("{}")}
	req.Send(msg)
	response, err := req.Recv()
	if err != nil {
		return err
	}
	sendMessage(conn, response)
	msg = [][]byte{[]byte("commands"), []byte("{}")}
	req.Send(msg)
	response, err = req.Recv()
	sendMessage(conn, response)
	return nil
}

func getSubMessage(conn * websocket.Conn, ctx *zmq.Context, shouldQuit <-chan bool) error {
	requestStatus(conn, ctx)
	sub, chans, err := getSubChannels(ctx)
	if err != nil {
		return err
	}
	defer sub.Close()
	defer chans.Close()
	for {
		select {
		case msg := <-chans.In():
			sendMessage(conn, msg)
		case <-shouldQuit:
			return nil
		case err := <-chans.Errors():
			return err
		}
	}
	return nil
}

func getSubChannels(ctx *zmq.Context) (sub *zmq.Socket, chans *zmq.Channels, err error) {
	uid := "gadgets ctrl";
	sub, err = ctx.Socket(zmq.Sub)
	if err != nil {
		return sub, chans, err
	}
	if err = sub.Connect("tcp://192.168.1.16:6111"); err != nil {
		return sub, chans, err
	}
	sub.Subscribe([]byte("UPDATE"))
	sub.Subscribe([]byte(uid))
	chans = sub.Channels()
	return sub, chans, err
}

type command struct {
	Event string
	Message map[string]interface{}
}


func getSockMessage(conn *websocket.Conn, ctx *zmq.Context, done chan<- bool) error {
	pub, err := ctx.Socket(zmq.Pub)
	if err = pub.Connect("tcp://192.168.1.16:6112"); err != nil {
		return err
	}
	defer pub.Close()
	if err != nil {
		return err
	}
	
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			done <- true
			return err
		}
		if messageType == websocket.TextMessage {
			cmd := &command{}
			err := json.Unmarshal(p, cmd)
			if err == nil {
				b, _ := json.Marshal(cmd.Message)
				msg := [][]byte{
					[]byte(cmd.Event),
					b,
				}
				pub.Send(msg)
			}
			
		} else if messageType == websocket.CloseMessage || messageType == -1 {
			done <- true
			return nil
		}
	}
	return nil
}





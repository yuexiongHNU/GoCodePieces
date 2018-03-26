/*
	declare the client class/struct
*/
package main

import (
	"time"
	"github.com/gorilla/websocket"
	"log"
	"bytes"
	"net/http"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
	readBuffer = 1024
	writeBuffer = 1024
)

var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: readBuffer,
	WriteBufferSize: writeBuffer,
}

type Client struct {
	hub *Hub
	conn *websocket.Conn
	// used to receive the message
	send chan []byte
}

func (cl *Client) readPump() {
	defer func() {
		cl.hub.unregister <- cl
		cl.conn.Close()
	}()
	cl.conn.SetReadLimit(maxMessageSize)
	cl.conn.SetReadDeadline(time.Now().Add(pongWait))
	cl.conn.SetPongHandler(func(string) error { cl.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := cl.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		cl.hub.broadcast <- message
	}
}

func (cl *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cl.conn.Close()
	}()
	for {
		select {
		case message, ok := <-cl.send:
			cl.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				cl.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// return a writer for the next message to send
			w, err := cl.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(cl.send)
			for i :=0; i < n; i++ {
				w.Write(newline)
				w.Write(<-cl.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <- ticker.C:
			cl.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := cl.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				// log.Fatal("Write PingMessage Error:", err)
				return
			}

		}
	}
}

func serverWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
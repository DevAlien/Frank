package servers

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

type SocketIoServer struct {
	Server *socketio.Server
}

func NewSocketIoServer() SocketIoServer {
	var socketIoServer SocketIoServer

	socketIoServer.Server, _ = socketio.NewServer([]string{"websocket"})

	socketIoServer.startServer()

	return socketIoServer
}

func (sis *SocketIoServer) startServer() {
	sis.Server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("bot")
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	sis.Server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", sis.Server)
	log.Println("Serving at localhost:5000...")
	go func() {
		log.Fatal(http.ListenAndServe(":5000", nil))
	}()

}

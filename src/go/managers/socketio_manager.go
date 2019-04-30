package managers

import (
	"strconv"
	"net/http"

	"frank/src/go/helpers/log"
	"frank/src/go/models"

	"github.com/googollee/go-socket.io"
)

type SocketIoServer struct {
	Server *socketio.Server
	config *models.WebSocket
}


func NewSocketIoServer(config *models.WebSocket) SocketIoServer {
	var socketIoServer SocketIoServer

	socketIoServer.config = config
	socketIoServer.Server, _ = socketio.NewServer([]string{"websocket"})

	socketIoServer.startServer()

	return socketIoServer
}

func (sis *SocketIoServer) startServer() {
	sis.Server.On("connection", func(so socketio.Socket) {
		log.Log.Debug("Socket On Connection")
		so.Join("bot")
		so.On("disconnection", func() {
			log.Log.Debug("Socket On Disconnection")
		})
	})
	sis.Server.On("error", func(so socketio.Socket, err error) {
		log.Log.Debug("Error Socket.io %s", err)
	})

	http.Handle("/socket.io/", sis.Server)
	
	port := ":" + strconv.Itoa(5000)
	if sis.config.Port != 0 {
		port = ":" + strconv.Itoa(sis.config.Port)
	}

	go func() {
		log.Log.Debug("Starting Socket.io server at %s", port)
		http.ListenAndServe(port, nil)
	}()

}

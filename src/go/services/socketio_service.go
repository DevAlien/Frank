package services

import (
	"strconv"
	"net/http"

	"frank/src/go/helpers/log"
	"frank/src/go/helpers"
	"frank/src/go/managers"
	"frank/src/go/models"

	"github.com/googollee/go-socket.io"
)

type SocketIoServer struct {
	Server *socketio.Server
	config *models.WebSocket
}

var SocketIo SocketIoServer

func NewSocketIoServer(config *models.WebSocket) SocketIoServer {
	SocketIo = SocketIoServer{}

	SocketIo.config = config
	SocketIo.Server, _ = socketio.NewServer([]string{"websocket"})

	SocketIo.startServer()

	return SocketIo
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

	//add text endpoint to send text messages via socket.io
	sis.Server.On("text", func(msg string) (bool, string) {
		commands := helpers.CheckCommands(msg)
		go managers.ManageCommands(commands)
		return len(commands) > 0, "asd" //TODO return right stuff
	})

	go func() {
		log.Log.Debug("Starting Socket.io server at %s", port)
		http.ListenAndServe(port, nil)
	}()

}

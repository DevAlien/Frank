package managers

import (
	"net"
	"net/http"

	"frank/src/go/helpers/log"
	"frank/src/go/models"



	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
)

type MqttServer struct {
	Server *gmqtt.Server
	config *models.HTTP
}

var defaultMqttPort = 8080
var Mqtt MqttServer

func NewMqttServer(config *models.HTTP) {
	log.Log.Debug("MQTT start")
	if config.Disabled == true {
		return
	}
	log.Log.Debug("Start after")
	Mqtt = MqttServer{}
	Mqtt.config = config
	
	Mqtt.Server = gmqtt.NewServer()
	ln, err := net.Listen("tcp", ":1883")
	if err != nil {
		log.Log.Error(err.Error())
		return
	}
	Mqtt.Server.AddTCPListenner(ln)
	ws := &gmqtt.WsServer{
		Server: &http.Server{Addr: ":8090"},
	}
	wss := &gmqtt.WsServer{
		Server:   &http.Server{Addr: ":8091"},
		CertFile: "./testcerts/server.crt",
		KeyFile:  "./testcerts/server.key",
	}
	Mqtt.Server.AddWebSocketServer(ws, wss)

	Mqtt.Server.RegisterOnPublish(func(client *gmqtt.Client, publish *packets.Publish) bool {
		log.Log.Debug("New message publushied to ", string(publish.TopicName), " => ", string(publish.String()))
		return true
		// if client.ClientOptions().Username == "subscribeonly" {
		// 	client.Close()
		// 	return false
		// }
		// //Only qos1 & qos0 are acceptable(will be delivered)
		// if publish.Qos == packets.QOS_2 {
		// 	return false
		// }
		// return true
	})
	Mqtt.Server.RegisterOnSubscribe(func(client *gmqtt.Client, topic packets.Topic) uint8 {
		log.Log.Debug("New Subscription to ", topic.Name)
		pub := &packets.Publish{
			Qos:       uint8(1),
			TopicName: []byte("/mremond/test-topic-1"),
			Payload:   []byte("Test maa"),
		}
		Mqtt.Server.Publish(pub)
		return topic.Qos
		// if client.ClientOptions().Username == "root" { //alow root user to subscribe whatever he wants
		// 	return topic.Qos
		// } else {
		// 	if topic.Qos <= packets.QOS_1 {
		// 		return topic.Qos
		// 	}
		// 	return packets.QOS_1   //for other users, the maximum QoS level is QoS1
		// }
	})

	go func() {
		log.Log.Debug("Starting MQTT server at %s", "1883")
		Mqtt.Server.Run()
	}()
}

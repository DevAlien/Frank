package services

import (
	"fmt"
	"math/rand"

	"frank/src/go/helpers/log"
	"frank/src/go/models"

	"gosrc.io/mqtt"
)

//MqttClientService struct definition of the Service.
type MqttClientService struct {
	Client *mqtt.ClientManager
	config *models.HTTP
}

//MqttClient single instance for the MqttClient.
var MqttClient MqttClientService

//NewMqttClient Starts a new MqttClient, loads the listeners and handles the messages
func NewMqttClient(config *models.HTTP) {
	if config.Disabled == true {
		return
	}

	MqttClient = MqttClientService{
		config: config,
	}

	address := "tcp://localhost:1883"
	client := mqtt.NewClient(address)
	client.ClientID = fmt.Sprintf("Frank-MQTT-Client-%d", rand.Int())
	log.Log.Debug("Connecting on: %s\n", client.Address)

	messages := make(chan mqtt.Message)
	client.Messages = messages

	postConnect := func(c *mqtt.Client) {
		log.Log.Info("Client Connected")
		name := "/mremond/test-topic-1"
		topic := mqtt.Topic{Name: name, QOS: 0}
		c.Subscribe(topic)
	}

	cm := mqtt.NewClientManager(client, postConnect)
	MqttClient.Client = cm
	log.Log.Debug("Starting MQTT client at %s", address)
	cm.Start()

	go handleMessages(messages, cm.Client)
}

func handleMessages(messages chan mqtt.Message, client *mqtt.Client) {
	for m := range messages {
		log.Log.Debug("DIO")
		log.Log.Debug("Received message from MQTT server on topic %s: %s\n", m.Topic, string(m.Payload))
		client.Publish("test/topic", []byte(fmt.Sprintf("%s %s", string(m.Payload), "RESPONSE")))
	}
}

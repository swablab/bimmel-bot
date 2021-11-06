package handler

import (
	"log"
	"swablab-bot/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttHandler struct {
	config     config.MqttConfig
	mqttClient mqtt.Client
}

func (handler *mqttHandler) SendMessage(message string) {
	handler.mqttClient.Publish("discord", 2, false, message)
}

func (handler *mqttHandler) Close() {
	handler.mqttClient.Disconnect(2)
}

//NewMqttMessageHandler Factory function for mqtt with authentication
func NewMqttMessageHandler(cfg config.MqttConfig) (*mqttHandler, error) {
	handler := new(mqttHandler)
	handler.config = cfg

	opts := mqtt.NewClientOptions().AddBroker(cfg.Host)

	//mqtt
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	if !cfg.AllowAnonymousAuthentication {
		opts.Username = cfg.Username
		opts.Password = cfg.Password
	}

	handler.mqttClient = mqtt.NewClient(opts)

	if token := handler.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	log.Print("successfully connected to mqtt-broker")
	return handler, nil
}

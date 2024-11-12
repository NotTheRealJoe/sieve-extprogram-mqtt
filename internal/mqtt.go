package internal

import (
	"fmt"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var randStrRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func randString(length uint) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = randStrRunes[rand.Intn(len(randStrRunes))]
	}
	return string(b)
}

func MQTTSetup(config Config) (*mqtt.Client, error) {
	// TODO: Automatically try reconnecting
	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
		panic(err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Broker, config.Port))
	opts.SetClientID("sieve-extprogram-mqtt-" + randString(6))
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &client, nil
}

func MQTTPublish(client mqtt.Client, topic string, message []byte) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
}

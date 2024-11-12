package internal

import (
	"crypto/tls"
	"crypto/x509"
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

func newTLSConfig() (*tls.Config, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs: certPool,
	}, nil
}

func MQTTSetup(config Config) (*mqtt.Client, error) {
	// TODO: Automatically try reconnecting
	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
		panic(err)
	}

	opts := mqtt.NewClientOptions()
	var broker string
	if config.UseSSL {
		broker = fmt.Sprintf("ssl://%s:%d", config.Broker, config.Port)
		tlsConfig, err := newTLSConfig()
		if err != nil {
			panic(err)
		}
		opts.SetTLSConfig(tlsConfig)
	} else {
		broker = fmt.Sprintf("tcp://%s:%d", config.Broker, config.Port)
	}
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.AddBroker(broker)
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

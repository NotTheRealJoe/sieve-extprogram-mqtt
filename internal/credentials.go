package internal

type Config struct {
	Broker         string `json:"broker"`
	Port           uint   `json:"port"`
	ClientIDPrefix string `json:"clientIdPrefix"`
	Topic          string `json:"topic"`
}

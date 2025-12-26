package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PGDB   PGDatabase
	Server Server
	MQTT   MQTTBroker
}

type PGDatabase struct {
	DB_Host string
	DB_Name string
	DB_Port string
	DB_User string
	DB_Pass string
}

type MQTTBroker struct {
	Port     string
	Address  string
	ClientID string
}

type Server struct {
	IP_Addr string
}

func NewEnvironment() *Env {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	return &Env{
		PGDB: PGDatabase{
			DB_Host: os.Getenv("DB_HOST"),
			DB_Name: os.Getenv("DB_NAME"),
			DB_Port: os.Getenv("DB_PORT"),
			DB_User: os.Getenv("DB_USER"),
			DB_Pass: os.Getenv("DB_PASS"),
		},
		Server: Server{
			IP_Addr: os.Getenv("IP_ADDR"),
		},
		MQTT: MQTTBroker{
			Port:    os.Getenv("MQTT_PORT"),
			Address: os.Getenv("MQTT_ADDR"),
			ClientID: os.Getenv("MQTT_CLIENT_ID"),
		},
	}
}

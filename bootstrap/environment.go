package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PGDB    PGDatabase
	Server  Server
	Kafka   KafkaBroker
}

type PGDatabase struct {
	DB_Host string
	DB_Name string
	DB_Port string
	DB_User string
	DB_Pass string
}

type KafkaBroker struct {
	Port    string
	Address string
}

type ElasticSearch struct {
	Port    string
	Address string
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
		Kafka: KafkaBroker{
			Port:    os.Getenv("KAFKA_PORT"),
			Address: os.Getenv("KAFKA_ADDR"),
		},
	}
}

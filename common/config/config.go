package config

import "time"

const (
	GatewayPort = "localhost:8081"
	ServerPort  = "localhost:9000"

	SecretKey     = "secret!!hehehe"
	TokenDuration = 15 * time.Minute
)

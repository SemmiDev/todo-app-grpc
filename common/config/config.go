package config

import "time"

const (
	RestServerPort = "localhost:8081"
	GRPCServerPort = "localhost:9000"

	SecretKey     = "secret!!hehehe"
	TokenDuration = 15 * time.Minute
)

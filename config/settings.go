package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

var AMQPConnectionURL = string(os.Getenv("AMQPConnectionURL"))

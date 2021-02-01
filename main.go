package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/natfaulk/distributed_computing_server/internal/filehandling"
	"github.com/natfaulk/distributed_computing_server/internal/nflogger"
	"github.com/natfaulk/distributed_computing_server/internal/server"
)

var logger *log.Logger = nflogger.Make("Main")

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Print("Error loading .env file")
	}

	filehandling.MakeDirectories()
	server.BeginServer()
}

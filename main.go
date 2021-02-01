package main

import (
	"log"

	"github.com/natfaulk/distributed_computing_server/internal/filehandling"
	"github.com/natfaulk/distributed_computing_server/internal/nflogger"
	"github.com/natfaulk/distributed_computing_server/internal/server"
)

var logger *log.Logger = nflogger.Make("Main")

func main() {
	filehandling.MakeDirectories()

	// logger.Println(uuid.NewString())

	// tp := tasks.Taskpool{}
	// tp.SaveToFile()

	// cl := clients.ClientList{}

	// id := uuid.NewString()
	// logger.Println(cl.Clients)
	// cl.AddClient(id)
	// logger.Println(cl.Clients)
	// time.Sleep(5 * time.Second)
	// cl.Ping(id)
	// logger.Println(cl.Clients)

	server.BeginServer()
}

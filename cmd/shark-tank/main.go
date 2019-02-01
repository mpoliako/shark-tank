package main

import (
	"log"

	"gitlab.corp.cloudsimple.com/cloudsimple/csos/incubator/shark-tank/server"
	"gitlab.corp.cloudsimple.com/cloudsimple/csos/incubator/shark-tank/storage"
)

func main() {
	repo, err := storage.NewCassandraRepository("127.0.0.1", 9042)
	if err != nil {
		log.Fatalf("failed to init cassandra repo err: %v", err)
	}
	srv := server.NewServer(repo)
	log.Fatal(srv.Run(":8080"))
}

package api

import (
	"blockchain-transactions/internal/db"
	"github.com/google/uuid"
)

func Start(port int) {
	dbx := db.GetConnection()

	server := newServer(port, dbx, uuid.New().String())
	server.Start()
}

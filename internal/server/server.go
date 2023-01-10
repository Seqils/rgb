package server

import (
	"rgb/internal/database"
	"rgb/internal/store"
	"rgb/internal/conf"
)

func Start(cfg conf.Config) {
	store.SetDBConnection(database.NewDBOptions(cfg))
	
	router := setRouter()

	router.Run(":8080")
}
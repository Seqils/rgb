package server

import (
	"rgb/internal/database"
	"rgb/internal/store"
	"rgb/internal/conf"
)

func Start(cfg conf.Config) {
	jwtSetup(cfg)
	
	store.SetDBConnection(database.NewDBOptions(cfg))
	
	router := setRouter()

	router.Run(":8080")
}
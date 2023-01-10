package store

import(
	"log"

	"github.com/go-pg/pg/v10"
)

var db *pg.DB

func SetDBConnection(dbOpts *pg.Options) {
	if dbOpts == nil {
		log.Panicln("DB options can't be nil")
	} 
	db = pg.Connect(dbOpts)
	log.Printf("%s", db)
} 

func GetDBConnection() *pg.DB { return db }
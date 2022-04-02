package db

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cfg interface {
	Path() string
}

type database struct {
	log zerolog.Logger
	db  *gorm.DB
}

func DB(l zerolog.Logger, c Cfg) (*database, error) {
	db, err := gorm.Open(sqlite.Open(c.Path()), &gorm.Config{})
	if err != nil {
		l.Err(err).Msg("Failed database open")
		return nil, err
	}

	return &database{
		log: l,
		db:  db,
	}, nil
}

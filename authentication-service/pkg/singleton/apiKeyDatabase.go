package singleton

import (
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v3"
)

var (
	db   *DataBaseSingleton
	once sync.Once
)

type ApiKey struct {
	Name       string `yaml:"name"`       // E-Mail address of the user
	Key        string `yaml:"key"`        // API key in no particular format, should only contain base64 characters and have a max length of 100
	ValidUntil string `yaml:"validUntil"` // RFC3339 formatted date
}

type DataBaseSingleton struct {
	ApiKeys []ApiKey `yaml:"apiKeys"`
}

func GetDatabaseInstance(dbPath string) *DataBaseSingleton {
	once.Do(func() {
		db_file, err := os.ReadFile(dbPath)
		if err != nil {
			log.Err(err).Msg("Error reading db.yaml")
			panic(err)
		}
		log.Debug().Msg("db.yaml read")

		db = &DataBaseSingleton{}
		err = yaml.Unmarshal(db_file, db)

		if err != nil {
			log.Err(err).Msg("Error unmarshalling db.yaml")
			panic(err)
		}
		log.Debug().Msg("db.yaml unmarshalled")
	})
	return db
}

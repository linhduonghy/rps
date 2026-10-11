package main

import (
	"log"
	"os"
	"rps/config"
	"rps/migration"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: migrate [up|down|version]")
	}

	cmd := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migration.New(cfg.Database.URL())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := migrator.Run("close"); err != nil {
			log.Println(err)
		}
	}()

	if err := migrator.Run(cmd); err != nil {
		log.Fatal(err)
	}
}

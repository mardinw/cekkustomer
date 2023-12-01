package servers

import (
	"cekkustomer/configs"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sethvargo/go-envconfig"
)

func Migrate(db *sql.DB) error {
	var config configs.AppConfiguration

	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal(err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println(err.Error)
		return err
	}

	startTime := time.Now()

	if config.AppEnv == "production" {
		migrator, err := migrate.NewWithDatabaseInstance(
			"s3://db-migrate-src",
			"postgres", driver,
		)
		if err != nil {
			log.Println(err.Error())
		}

		err = migrator.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Println(err.Error())
			return err
		} else {
			log.Println("database migration applied")
		}
	} else {
		migrator, err := migrate.NewWithDatabaseInstance(
			"file:///home/petr0max/Public/go/src/github.com/cekkustomer/db/migrations",
			"postgres", driver,
		)
		if err != nil {
			log.Println(err.Error())
		}

		err = migrator.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Println(err.Error())
			return err
		} else {
			log.Println("database migrations applied")
		}
	}

	duration := time.Since(startTime)
	log.Printf("migration completed in %v", duration)

	return nil
}

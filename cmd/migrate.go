package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // needed driver, only here
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/conf"
	"github.com/caquillo07/rotom-bot/repository"
)

type migrateLogger zap.Logger

func (l *migrateLogger) Verbose() bool {
	return true
}
func (l *migrateLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func init() {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run the database migrations",
		Run:   runMigrateCommand,
	}
	cmd.Flags().Bool("ignore-no-change", false, "Allow migrations with no change")
	rootCmd.AddCommand(cmd)
}

func runMigrateCommand(cmd *cobra.Command, args []string) {
	ignoreNoChange, err := cmd.Flags().GetBool("ignore-no-change")
	if err != nil {
		log.Fatalln(err)
	}

	config, err := conf.LoadConfig(viper.GetViper())
	if err != nil {
		log.Fatalln(err)
	}

	if err := migrate(config.Database); err != nil {
		// Ignore no change error if ignore-no-change flag is set.
		if err == gomigrate.ErrNoChange && ignoreNoChange == true {
			log.Println(err)
		} else {
			log.Fatalln(err)
		}
	}

	log.Println("Migrations run")
}

// migrate performs a migration with the given configuration
func migrate(config repository.Config) error {
	// Look for migrations folder in working path, else at the same level as a .git folder in a parent
	var migrationFolderName string
	if config.MigrationFolder != "" {
		migrationFolderName = config.MigrationFolder
	} else {
		migrationFolderName = "migrations"
	}
	path := migrationFolderName
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if path, err = filepath.Abs(path); err != nil {
			return err
		}
		path = filepath.Dir(path)
		for {
			if _, err = os.Stat(filepath.Join(path, migrationFolderName)); err == nil {
				break
			}
			newPath := filepath.Dir(path)
			if newPath == path {
				break
			}
			path = newPath
		}

		path = filepath.Join(path, migrationFolderName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatalln("Unable to locate migrations folder")
		}
	}

	db, err := repository.Open(config)
	if err != nil {
		log.Fatalln(err)
	}
	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := gomigrate.NewWithDatabaseInstance(
		"file://"+path,
		config.Driver,
		driver,
	)
	if err != nil {
		return err
	}

	m.Log = &migrateLogger{}

	return m.Up()
}

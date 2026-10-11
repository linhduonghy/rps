package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Command string

const (
	Up      Command = "up"
	Down    Command = "down"
	Version Command = "version"
	Close   Command = "close"
)

func parseCommand(cmd string) (Command, error) {
	switch Command(cmd) {
	case Up, Down, Version, Close:
		return Command(cmd), nil
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

const migrationSource = "file://sql/migrations"

type Migrator struct {
	m *migrate.Migrate
}

func New(databaseURL string) (*Migrator, error) {
	m, err := migrate.New(
		migrationSource,
		databaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return &Migrator{
		m: m,
	}, nil
}

func (m *Migrator) Run(cmd string) error {
	command, _ := parseCommand(cmd)
	switch command {
	case Up:
		return m.up()
	case Down:
		return m.down()
	case Version:
		return m.version()
	case Close:
		return m.close()
	default:
		return fmt.Errorf("unsupported command %q", cmd)
	}
}

func (m *Migrator) up() error {
	if err := m.m.Up(); err != nil &&
		!errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (m *Migrator) down() error {
	if err := m.m.Down(); err != nil &&
		!errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (m *Migrator) version() error {
	version, dirty, err := m.m.Version()
	if err != nil {
		return err
	}

	fmt.Printf(
		"Version: %d (dirty=%v)\n",
		version,
		dirty,
	)

	return nil
}

func (m *Migrator) close() error {
	srcErr, dbErr := m.m.Close()

	if srcErr != nil {
		return srcErr
	}

	return dbErr
}

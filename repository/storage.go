package repository

import "errors"

var (
	// ErrRecordNotFound is a generic error returned when a given record is not
	// found
	ErrRecordNotFound = errors.New("record not found")
)

// Storage will define the required methods required by the storage manager
// inside this bot. This is an interface to allow for future implementations
// of non-sql based storages, since someone self hosting may not have PostgreSQL
// running.
type Storage interface {

	// CreateGuildConfig creates a new config for a guild, returns error on
	// failure
	CreateGuildSettings(config *GuildSettings) error
}

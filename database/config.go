package database

// Config provides database configuration
type Config struct {
	// Database driver
	Driver string

	// Database connection string
	URL string

	// Log will enable or disable query logging
	Log bool

	// Check if there is a custom migrations folder
	MigrationFolder string
}

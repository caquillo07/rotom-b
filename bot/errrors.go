package bot

// botError is a public error. This error should be returned from command
// handlers when we want the user to see it, otherwise a return error should
// be returned
type botError struct {
	title   string
	details string
}

func (e botError) Error() string {
	return e.title + " " + e.details
}

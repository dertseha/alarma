package config

// TimeSpan describes one action.
type TimeSpan struct {
	// ID uniquely identifies the time span.
	ID string `json:"id"`
	// Enabled is set to true to indicate this time span should be active.
	Enabled bool `json:"enabled"`

	// From describes the starting time of the time span. Format: "HH:MM".
	From string `json:"from"`
	// To describes the stopping time of the time span. Format: "HH:MM".
	To string `json:"to"`

	// Path points to the base directory within which audio files are to be searched.
	Path string `json:"path"`
}

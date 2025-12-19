package state

import "time"

// ProfileState represents persisted state of a LinkedIn profile
type ProfileState struct {
	ProfileURL     string
	Visited        bool
	ConnectionSent bool
	Status         string
	CreatedAt      time.Time
	LastActionAt   *time.Time
}

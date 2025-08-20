package entities

import "time"

type SaveData struct {
	RemainingChips int       `json:"remainingChips"`
	LastResetAt    time.Time `json:"lastResetAt"`
}

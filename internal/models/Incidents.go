package models

import "time"

type Incidents struct {
	ID uint `gorm:"primaryKey"`

	ServiceName  string
	Level        string
	Message      string
	Severity     string
	RootCause    string
	Impact       string
	SuggestedFix string
	Confidence   int

	CreatedAt time.Time
}
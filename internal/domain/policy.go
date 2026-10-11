package domain

import "time"

type Policy struct {
	ID           int64
	ResourceType string
	ResourceID   string
	Version      string

	Statements []Statement

	CreatedAt time.Time
	UpdatedAt time.Time
}

package domain

import "time"

type Statement struct {
	ID     int64
	SID    string
	Effect Effect

	Actions    []string
	Resources  []string
	Principals []Principal

	CreatedAt time.Time
}

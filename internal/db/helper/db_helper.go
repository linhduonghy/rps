package helper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToTime(ts pgtype.Timestamp) time.Time {
	if ts.Valid {
		return ts.Time
	}

	return time.Time{}
}

func ToPgText(v string) pgtype.Text {
	return pgtype.Text{
		String: v,
		Valid:  true,
	}
}

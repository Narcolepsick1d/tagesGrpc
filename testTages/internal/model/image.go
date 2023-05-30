package model

import (
	"time"
)

type Image struct {
	Name        string
	CreatedTime time.Time
	UpdatedTime time.Time
}

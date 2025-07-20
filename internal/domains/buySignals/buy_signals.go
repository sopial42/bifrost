package buysignals

import (
	"fmt"
	"time"

	"github.com/bifrost/internal/domains/common"
	"github.com/google/uuid"
)

type Details struct {
	Name       Name
	ID         ID
	BusinessID BusinessID
	Fullname   Fullname
	Pair       common.Pair
	Interval   common.Interval
	Date       time.Time
	Price      float64
	Metadata   map[string]any // TODO use generics over MSI
}

type ID uuid.UUID

func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, uuid.UUID(i).String())), nil
}

// Name is the buySignal name
type Name string

// Fullname is the buySignal name with params
type Fullname string

// BusinessID is used to ensure buysignal uniqueness
type BusinessID string

package buysignals

import (
	"fmt"
	"time"

	"github.com/bifrost/pkg/domains/common"
	"github.com/google/uuid"
)

type Details struct {
	Name       Name           `json:"name"`
	ID         ID             `json:"id"`
	BusinessID BusinessID     `json:"business_id"`
	Fullname   Fullname       `json:"fullname"`
	Pair       common.Pair    `json:"pair"`
	Date       Date           `json:"date"`
	Price      float64        `json:"price"`
	Metadata   map[string]any `json:"metadata"`
}

type ID uuid.UUID

type Date time.Time

func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, uuid.UUID(i).String())), nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format(time.RFC3339))), nil
}

// Name is the buySignal name
type Name string

// Fullname is the buySignal name with params
type Fullname string

// BusinessID is used to ensure buysignal uniqueness
type BusinessID string

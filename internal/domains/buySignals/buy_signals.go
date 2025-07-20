package buysignals

import (
	"fmt"
	"time"

	"github.com/bifrost/internal/domains/common"
	"github.com/google/uuid"
)

type Details struct {
	Name       Name           `json:"name"`
	ID         ID             `json:"id"`
	BusinessID BusinessID     `json:"business_id"`
	Fullname   Fullname       `json:"fullname"`
	Pair       common.Pair    `json:"pair"`
	Date       time.Time      `json:"date"`
	Price      float64        `json:"price"`
	Metadata   map[string]any `json:"metadata"`
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

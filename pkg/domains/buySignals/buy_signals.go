package buysignals

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

type Details struct {
	Name       Name           `json:"name"`
	ID         *ID            `json:"id,omitempty"`
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

func (i *ID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*i = ID(parsed)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format(time.RFC3339))), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return fmt.Errorf("failed to unmarshal date: %w", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}

	*d = Date(parsedTime)
	return nil
}

// Name is the buySignal name
type Name string

// Fullname is the buySignal name with params
type Fullname string

// BusinessID is used to ensure buysignal uniqueness
type BusinessID string

func ParseSignalStrategies(argsSignalStrategies []string) ([]Name, error) {
	signalsStrat := make([]Name, len(argsSignalStrategies))
	errors := []string{}
	for i, ss := range argsSignalStrategies {
		if !AllAvailableSignalStrategies[Name(ss)] {
			errors = append(errors, ss)
		} else {
			signalsStrat[i] = Name(ss)
		}
	}

	if len(errors) > 0 {
		return []Name{}, fmt.Errorf("signalStrategy args not allowed: %s", errors)
	}
	return signalsStrat, nil
}

var MorningStarName Name = "morningStar"
var RsiDivergenceName Name = "rsiDivergence"

var AllAvailableSignalStrategies = map[Name]bool{
	MorningStarName:   true,
	RsiDivergenceName: true,
}

package buysignals

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

const LoggerKeyName = "buy_signals"
const LoggerKeyFullname = "buy_signals_fullname"

type Details struct {
	Name       Name            `json:"name"`
	ID         *ID             `json:"id,omitempty"`
	BusinessID BusinessID      `json:"business_id"`
	Fullname   Fullname        `json:"fullname"`
	Pair       common.Pair     `json:"pair"`
	Interval   common.Interval `json:"interval"`
	Date       Date            `json:"date"`
	Price      float64         `json:"price"`
	Metadata   Metadata        `json:"metadata,omitempty"`
}

var MorningStarName Name = "morningStar"
var RSIDivergenceName Name = "rsiDivergence"

var AllAvailableSignalStrategies = map[Name]bool{
	MorningStarName:   true,
	RSIDivergenceName: true,
}

type Metadata map[string]any

func (m *Metadata) UnmarshalJSON(data []byte) error {
	var metadataMap map[string]any
	if err := json.Unmarshal(data, &metadataMap); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	*m = Metadata(metadataMap)
	return nil
}

func (m *Metadata) SetMetadata(metadata any) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("unable to marshal buySignalMetadata: %w", err)
	}

	// metadata should be a map[string]any
	var metadataMap map[string]any
	err = json.Unmarshal(metadataBytes, &metadataMap)
	if err != nil {
		return fmt.Errorf("unable to unmarshal buySignalMetadata: %w", err)
	}

	*m = metadataMap
	return nil
}

func MetadataToStruct[T any](metadata Metadata) (T, error) {
	var result T

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return result, fmt.Errorf("unable to marshal metadata: %w", err)
	}

	if err = json.Unmarshal(jsonData, &result); err != nil {
		return result, fmt.Errorf("unable to unmarshal to struct: %w", err)
	}

	return result, nil
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
	signalsStrat := make([]Name, 0)
	errors := []string{}
	for _, ss := range argsSignalStrategies {
		if !AllAvailableSignalStrategies[Name(ss)] {
			errors = append(errors, ss)
		} else {
			signalsStrat = append(signalsStrat, Name(ss))
		}
	}

	if len(errors) > 0 {
		return []Name{}, fmt.Errorf("signalStrategy args not allowed: %s", errors)
	}
	return signalsStrat, nil
}

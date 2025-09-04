package candles

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

type Date time.Time

func (d Date) Before(other Date) bool {
	return time.Time(d).Before(time.Time(other))
}

func (d Date) After(other Date) bool {
	return time.Time(d).After(time.Time(other))
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

// String method returns a human-readable date
func (c Date) String() string {
	return GetDateStrFromUnixTimeMilli(time.Time(c))
}

func GetDateStrFromUnixTimeMilli(date time.Time) string {
	return date.In(time.Now().Location()).Format(time.RFC3339)
}

type ID uuid.UUID

type Candle struct {
	ID       *ID             `json:"id,omitempty"`
	Date     Date            `json:"date"`
	Pair     common.Pair     `json:"pair"`
	Interval common.Interval `json:"interval"`
	Open     float64         `json:"open"`
	Close    float64         `json:"close"`
	High     float64         `json:"high"`
	Low      float64         `json:"low"`
	RSI      *RSI            `json:"rsi,omitempty"`
}

type RSI map[RSIPeriod]RSIValue
type RSIPeriod int64
type RSIValue float64

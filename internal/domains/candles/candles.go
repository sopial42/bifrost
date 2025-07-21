package candles

import (
	"fmt"
	"time"

	"github.com/bifrost/internal/domains/common"
	"github.com/google/uuid"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format(time.RFC3339))), nil
}

type ID uuid.UUID

func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, uuid.UUID(i).String())), nil
}

// String method returns a human-readable date
func (c Date) String() string {
	return GetDateStrFromUnixTimeMilli(time.Time(c))
}

func GetDateStrFromUnixTimeMilli(date time.Time) string {
	return date.In(time.Now().Location()).Format(time.RFC3339)
}

type Candle struct {
	ID       ID              `json:"id"`
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

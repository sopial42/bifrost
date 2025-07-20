package positions

import (
	"fmt"

	buySignals "github.com/bifrost/internal/domains/buySignals"
	"github.com/google/uuid"
)

type Details struct {
	ID          ID                  `json:"id"`
	SerialID    int64               `json:"serial_id"`
	Name        Name                `json:"name"`
	Fullname    Fullname            `json:"fullname"`
	BuySignalID buySignals.ID       `json:"buy_signal_id"`
	BuySignal   *buySignals.Details `json:"buy_signal,omitempty"`
	TP          float64             `json:"tp"`
	SL          float64             `json:"sl"`
	Metadata    map[string]any      `json:"metadata"`
}

type ID uuid.UUID

func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, uuid.UUID(i).String())), nil
}

type Name string
type Fullname string
type Metadata map[string]interface{}

package positions

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	buySignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

const LoggerKeyName = "positions"
const LoggerKeyFullname = "positions_fullname"

type Details struct {
	ID          *ID                  `json:"id"`
	SerialID    SerialID            `json:"serial_id"`
	Name        Name                `json:"name"`
	Fullname    Fullname            `json:"fullname"`
	BuySignalID buySignals.ID       `json:"buy_signal_id"`
	BuySignal   *buySignals.Details `json:"buy_signal,omitempty"`
	TP          float64             `json:"tp"`
	SL          float64             `json:"sl"`
	Metadata    map[string]any      `json:"metadata"`
	Ratio       *Ratio              `json:"ratio,omitempty"`
}

type Ratio struct {
	Value float64      `json:"value"`
	Date  candles.Date `json:"date"`
}

type SerialID int64

const (
	FibonacciName  Name = "fibonacci"
	PivotPointName Name = "pivotPoint"
)

var AllAvailablePositionStategies = map[Name]bool{
	FibonacciName:  true,
	PivotPointName: true,
}

type ID uuid.UUID

func (i ID) String() string {
	return uuid.UUID(i).String()
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

type Name string
type Fullname string
type Metadata map[string]interface{}

type Scenario struct {
	Pair          common.Pair
	Interval      common.Interval
	BuySignalName buySignals.Name
	PositionNames []Name
}

var (
	ArgsDefaultLimitStrategy []Name
)

func ParseSignalStrategies(argsPositionName []string) ([]Name, error) {
	names := make([]Name, 0)
	errors := []string{}
	for _, arg := range argsPositionName {
		if !AllAvailablePositionStategies[Name(arg)] {
			errors = append(errors, arg)
		} else {
			names = append(names, Name(arg))
		}
	}

	if len(errors) > 0 {
		return []Name{}, fmt.Errorf("positionNames args not allowed: %s", errors)
	}

	return names, nil
}

func GetScenarios(pairs []common.Pair, intervals []common.Interval, signals []buySignals.Name, positions []Name) *[]Scenario {
	res := make([]Scenario, 0)
	for _, p := range pairs {
		for _, i := range intervals {
			for _, s := range signals {
				res = append(res, Scenario{
					Pair:          p,
					Interval:      i,
					BuySignalName: s,
					PositionNames: positions,
				})
			}
		}
	}

	return &res
}

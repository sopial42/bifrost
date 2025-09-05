package common

import (
	"fmt"
	"time"
)

type Interval string

const IntervalLoggerKey = "interval"

const (
	// S1  Interval = "1s"
	M1  Interval = "1m"
	M3  Interval = "3m"
	M5  Interval = "5m"
	M15 Interval = "15m"
	M30 Interval = "30m"
	H1  Interval = "1h"
	H2  Interval = "2h"
	H4  Interval = "4h"
	H6  Interval = "6h"
	H8  Interval = "8h"
	H12 Interval = "12h"
	D1  Interval = "1d"
	W1  Interval = "1w"
)

func (i Interval) RoundDateToBeginingOfInterval(currentTime time.Time) *time.Time {
	var newTime time.Time

	switch i {
	case W1:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(7 * 24 * time.Hour)
	case D1:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(24 * time.Hour)
	case H12:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(12 * time.Hour)
	case H8:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(8 * time.Hour)
	case H6:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(6 * time.Hour)
	case H4:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(4 * time.Hour)
	case H2:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(2 * time.Hour)
	case H1:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(1 * time.Hour)
	case M1:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(1 * time.Minute)
	case M3:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(3 * time.Minute)
	case M5:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(5 * time.Minute)
	case M15:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(15 * time.Minute)
	case M30:
		newTime = time.Time(currentTime).In(time.UTC).Truncate(30 * time.Minute)
	}

	if newTime.IsZero() {
		return nil
	}

	return &newTime
}

func AddOneInterval(currentTime time.Time, i Interval) *time.Time {
	var newTime time.Time

	switch i {
	case W1:
		newTime = time.Time(currentTime).In(time.UTC).Add(7 * 24 * time.Hour)
	case D1:
		newTime = time.Time(currentTime).In(time.UTC).Add(24 * time.Hour)
	case H12:
		newTime = time.Time(currentTime).In(time.UTC).Add(12 * time.Hour)
	case H8:
		newTime = time.Time(currentTime).In(time.UTC).Add(8 * time.Hour)
	case H6:
		newTime = time.Time(currentTime).In(time.UTC).Add(6 * time.Hour)
	case H4:
		newTime = time.Time(currentTime).In(time.UTC).Add(4 * time.Hour)
	case H2:
		newTime = time.Time(currentTime).In(time.UTC).Add(2 * time.Hour)
	case H1:
		newTime = time.Time(currentTime).In(time.UTC).Add(1 * time.Hour)
	case M1:
		newTime = time.Time(currentTime).In(time.UTC).Add(1 * time.Minute)
	case M3:
		newTime = time.Time(currentTime).In(time.UTC).Add(3 * time.Minute)
	case M5:
		newTime = time.Time(currentTime).In(time.UTC).Add(5 * time.Minute)
	case M15:
		newTime = time.Time(currentTime).In(time.UTC).Add(15 * time.Minute)
	case M30:
		newTime = time.Time(currentTime).In(time.UTC).Add(30 * time.Minute)
	}

	if newTime.IsZero() {
		return nil
	}

	return &newTime
}

var Intervals = []Interval{
	// S1,
	M1,
	M3,
	M5,
	M15,
	M30,
	H1,
	H2,
	H4,
	H6,
	H8,
	H12,
	D1,
	W1,
}

var ArgsDefaultIntervals []string
var AllAvailableInterval map[Interval]bool

func init() {
	AllAvailableInterval = make(map[Interval]bool, len(Intervals))
	ArgsDefaultIntervals = make([]string, len(Intervals))
	for i, itv := range Intervals {
		AllAvailableInterval[itv] = true
		ArgsDefaultIntervals[i] = string(itv)
	}
}

func (i Interval) IsValid() bool {
	_, ok := AllAvailableInterval[i]
	return ok
}

func ParseInterval(arg string) (Interval, error) {
	_, ok := AllAvailableInterval[Interval(arg)]

	if !ok {
		return "", fmt.Errorf("wrong interval: %q", arg)
	} else {
		return Interval(arg), nil
	}
}

func ParseIntervals(argsInterval []string) ([]Interval, error) {
	intervals := make([]Interval, len(argsInterval))
	errors := []string{}
	for i, itv := range argsInterval {
		if !AllAvailableInterval[Interval(itv)] {
			errors = append(errors, itv)
		} else {
			intervals[i] = Interval(itv)
		}
	}

	if len(errors) > 0 {
		return []Interval{}, fmt.Errorf("interval args not allowed: %s", errors)
	}
	return intervals, nil
}

func (i Interval) String() string {
	return string(i)
}

func (i Interval) GetBelowInterval() Interval {
	var belowInterval Interval
	for _, testedInterval := range Intervals {
		if belowInterval == "" {
			belowInterval = testedInterval
		}

		if testedInterval == i {
			return belowInterval
		}

		belowInterval = testedInterval
	}

	return belowInterval
}

package sdk

import (
	"reflect"
	"testing"
	"time"

	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

func Test_createCandlesChunk(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		candles    *[]candles.Candle
		chunkSize  int
		wantChunks *[][]candles.Candle
	}{
		{
			name: "chunk size 1",
			candles: &[]candles.Candle{
				{
					Date:     candles.Date(now),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1000.0,
					Close:    1100.0,
					High:     1200.0,
					Low:      900.0,
				},
				{
					Date:     candles.Date(now.Add(time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1100.0,
					Close:    1200.0,
					High:     1300.0,
					Low:      1000.0,
				},
			},
			chunkSize: 1,
			wantChunks: &[][]candles.Candle{
				{
					{
						Date:     candles.Date(now),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1000.0,
						Close:    1100.0,
						High:     1200.0,
						Low:      900.0,
					},
				},
				{
					{
						Date:     candles.Date(now.Add(time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1100.0,
						Close:    1200.0,
						High:     1300.0,
						Low:      1000.0,
					},
				},
			},
		},
		{
			name: "chunk size 3 with 7 candles",
			candles: &[]candles.Candle{
				{
					Date:     candles.Date(now),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1000.0,
					Close:    1100.0,
					High:     1200.0,
					Low:      900.0,
				},
				{
					Date:     candles.Date(now.Add(time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1100.0,
					Close:    1200.0,
					High:     1300.0,
					Low:      1000.0,
				},
				{
					Date:     candles.Date(now.Add(2 * time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1200.0,
					Close:    1300.0,
					High:     1400.0,
					Low:      1100.0,
				},
				{
					Date:     candles.Date(now.Add(3 * time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1300.0,
					Close:    1400.0,
					High:     1500.0,
					Low:      1200.0,
				},
				{
					Date:     candles.Date(now.Add(4 * time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1400.0,
					Close:    1500.0,
					High:     1600.0,
					Low:      1300.0,
				},
				{
					Date:     candles.Date(now.Add(5 * time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1500.0,
					Close:    1600.0,
					High:     1700.0,
					Low:      1400.0,
				},
				{
					Date:     candles.Date(now.Add(6 * time.Hour)),
					Pair:     common.Pair("BTC/USD"),
					Interval: common.Interval("1h"),
					Open:     1600.0,
					Close:    1700.0,
					High:     1800.0,
					Low:      1500.0,
				},
			},
			chunkSize: 3,
			wantChunks: &[][]candles.Candle{
				{
					{
						Date:     candles.Date(now),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1000.0,
						Close:    1100.0,
						High:     1200.0,
						Low:      900.0,
					},
					{
						Date:     candles.Date(now.Add(time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1100.0,
						Close:    1200.0,
						High:     1300.0,
						Low:      1000.0,
					},
					{
						Date:     candles.Date(now.Add(2 * time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1200.0,
						Close:    1300.0,
						High:     1400.0,
						Low:      1100.0,
					},
				},
				{
					{
						Date:     candles.Date(now.Add(3 * time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1300.0,
						Close:    1400.0,
						High:     1500.0,
						Low:      1200.0,
					},
					{
						Date:     candles.Date(now.Add(4 * time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1400.0,
						Close:    1500.0,
						High:     1600.0,
						Low:      1300.0,
					},
					{
						Date:     candles.Date(now.Add(5 * time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1500.0,
						Close:    1600.0,
						High:     1700.0,
						Low:      1400.0,
					},
				},
				{
					{
						Date:     candles.Date(now.Add(6 * time.Hour)),
						Pair:     common.Pair("BTC/USD"),
						Interval: common.Interval("1h"),
						Open:     1600.0,
						Close:    1700.0,
						High:     1800.0,
						Low:      1500.0,
					},
				},
			},
		},
		{
			name:       "nil input",
			candles:    nil,
			chunkSize:  1,
			wantChunks: nil,
		},
		{
			name:       "empty input",
			candles:    &[]candles.Candle{},
			chunkSize:  1,
			wantChunks: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createChunk(tt.candles, tt.chunkSize)
			if !reflect.DeepEqual(got, tt.wantChunks) {
				t.Errorf("createChunk() = %v, want %v", got, tt.wantChunks)
			}
		})
	}
}

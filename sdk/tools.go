package sdk

import "github.com/sopial42/bifrost/pkg/domains/candles"

func createCandlesChunk(newCandles *[]candles.Candle, chunckSize int) *[][]candles.Candle {
	var chuncks [][]candles.Candle

	if newCandles == nil {
		return nil
	}

	for i := 0; i < len(*newCandles); i += chunckSize {
		chunk := (*newCandles)[i:min(i+chunckSize, len(*newCandles))]
		chuncks = append(chuncks, chunk)
	}

	return &chuncks
}

package common

import (
	"fmt"
)

type Pair string

// ~= Top 50 trading pair binance
// https://coinmarketcap.com/fr/exchanges/binance/
// https://coinmarketcap.com/api/documentation/v1/#operation/getV1ExchangeMarketpairsLatest

const PairLoggerKey = "pair"
const (
	BTCUSDC   Pair = "BTCUSDC"
	BTCUSDT   Pair = "BTCUSDT"
	SOLUSDC   Pair = "SOLUSDC"
	ETHUSDC   Pair = "ETHUSDC"
	ETHUSDT   Pair = "ETHUSDT"
	BNBUSDC   Pair = "BNBUSDC"
	XRPUSDC   Pair = "XRPUSDC"
	PEPEUSDC  Pair = "PEPEUSDC"
	DOGEUSDC  Pair = "DOGEUSDC"
	SOLBTC    Pair = "SOLBTC"
	SUIUSDC   Pair = "SUIUSDC"
	ATOMUSDC  Pair = "ATOMUSDC"
	MATICUSDC Pair = "MATICUSDC"
	NOTUSDC   Pair = "NOTUSDC"
	TRXUSDC   Pair = "TRXUSDC"
	ADAUSDC   Pair = "ADAUSDC"
	HYPEUSDC  Pair = "HYPEUSDC"
	LINKUSDC  Pair = "LINKUSDC"
	XLMUSDC   Pair = "XLMUSDC"
	AVAXUSDC  Pair = "AVAXUSDC"
	BCHUSDC   Pair = "BCHUSDC"
	NEARUSDC  Pair = "NEARUSDC"
)

var Pairs = []Pair{
	BTCUSDC,
	SOLUSDC,
	ETHUSDC,
	BNBUSDC,
	XRPUSDC,
	PEPEUSDC,
	DOGEUSDC,
	SOLUSDC,
	SUIUSDC,
	ATOMUSDC,
	MATICUSDC,
	NOTUSDC,
	TRXUSDC,
	ADAUSDC,
	HYPEUSDC,
	LINKUSDC,
	XLMUSDC,
	AVAXUSDC,
	BCHUSDC,
	NEARUSDC,
}

var AllAvailablePair map[Pair]bool
var ArgsDefaultPairs []string

func init() {
	AllAvailablePair = make(map[Pair]bool, len(Pairs))
	ArgsDefaultPairs = make([]string, len(Pairs))
	for i, pair := range Pairs {
		AllAvailablePair[pair] = true
		ArgsDefaultPairs[i] = string(pair)
	}
}

func ParsePair(arg string) (Pair, error) {
	_, ok := AllAvailablePair[Pair(arg)]

	if !ok {
		return "", fmt.Errorf("wrong pair: %q", arg)
	} else {
		return Pair(arg), nil
	}
}

func ParsePairs(argsPair []string) ([]Pair, error) {
	pairs := make([]Pair, len(argsPair))
	errors := []string{}
	for i, p := range argsPair {
		if !AllAvailablePair[Pair(p)] {
			errors = append(errors, p)
		} else {
			pairs[i] = Pair(p)
		}
	}

	if len(errors) > 0 {
		return []Pair{}, fmt.Errorf("pair args not allowed: %s", errors)
	}
	return pairs, nil
}

func (p Pair) String() string {
	return string(p)
}

func (p Pair) IsValid() bool {
	_, ok := AllAvailablePair[p]
	return ok
}

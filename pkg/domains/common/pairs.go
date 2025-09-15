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
	ADAUSDC   Pair = "ADAUSDC"
	ATOMUSDC  Pair = "ATOMUSDC"
	AVAXUSDC  Pair = "AVAXUSDC"
	BCHUSDC   Pair = "BCHUSDC"
	BNBUSDC   Pair = "BNBUSDC"
	BTCUSDC   Pair = "BTCUSDC"
	DOGEUSDC  Pair = "DOGEUSDC"
	ETHUSDC   Pair = "ETHUSDC"
	HYPEUSDC  Pair = "HYPEUSDC"
	LINKUSDC  Pair = "LINKUSDC"
	MATICUSDC Pair = "MATICUSDC"
	NEARUSDC  Pair = "NEARUSDC"
	NOTUSDC   Pair = "NOTUSDC"
	PEPEUSDC  Pair = "PEPEUSDC"
	SOLBTC    Pair = "SOLBTC"
	SOLUSDC   Pair = "SOLUSDC"
	SUIUSDC   Pair = "SUIUSDC"
	TRXUSDC   Pair = "TRXUSDC"
	XRPUSDC   Pair = "XRPUSDC"
	XLMUSDC   Pair = "XLMUSDC"
)

var Pairs = []Pair{
	ADAUSDC,
	ATOMUSDC,
	AVAXUSDC,
	BCHUSDC,
	BNBUSDC,
	BTCUSDC,
	DOGEUSDC,
	ETHUSDC,
	HYPEUSDC,
	LINKUSDC,
	MATICUSDC,
	NEARUSDC,
	NOTUSDC,
	PEPEUSDC,
	SOLBTC,
	SOLUSDC,
	SUIUSDC,
	TRXUSDC,
	XRPUSDC,
	XLMUSDC,
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

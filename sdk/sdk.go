package sdk

import "github.com/sopial42/bifrost/pkg/sdk"

type Bifrost interface {
	Candles
	BuySignals
}

type client struct {
	*sdk.Client
}

func NewBifrostClient(baseURL string) Bifrost {
	return &client{
		sdk.NewSDKClient(baseURL),
	}
}

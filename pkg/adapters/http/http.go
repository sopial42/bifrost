package http

import (
	"github.com/sopial42/bifrost/pkg/common/sdk"
	"github.com/sopial42/bifrost/pkg/ports"
)

type client struct {
	*sdk.Client
}

func NewBifrostHTTPClient(baseURL string) ports.Client {
	return &client{
		sdk.NewSDKClient(baseURL),
	}
}

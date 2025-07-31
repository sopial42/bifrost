package pinger

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	appErrors "github.com/sopial42/bifrost/pkg/errors"
)

type Pinger interface {
	Ping(c context.Context) error
}

type Pingers struct {
	Pingers []Pinger
}

func newPingers(pinger ...Pinger) *Pingers {
	return &Pingers{
		Pingers: pinger,
	}
}

// Ping checks the health of all registered pingers and returns a 200 OK response if all are healthy
// route should be formatted like "/ping"
func SetNewPingers(echo *echo.Echo, route string, pingers ...Pinger) {
	p := newPingers(pingers...)
	echo.GET(route, p.Ping)
}

func (p *Pingers) Ping(c echo.Context) error {
	ctx := c.Request().Context()
	for _, pinger := range p.Pingers {
		if err := pinger.Ping(ctx); err != nil {
			log.Printf("ping failed: %v", err)
			return appErrors.NewUnexpected("failed to ping", err)
		}
	}

	return c.NoContent(http.StatusOK)
}

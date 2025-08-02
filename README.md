# Bifrost

Bifrost brings your trading data to life.

It’s an API for creating and updating market data, buy signals, and positions.

A candle represents market data (currently OHLC and RSI).
A buy signal defines the date and price at which a buy order was placed.
A position sets the take-profit (TP) and stop-loss (SL) levels for a buy signal.

To support trading strategy analytics, buy signals and positions are always defined with a `name`, `fullname`, and `metadata`.

## Use it 


- Start dependencies using docker
```bash
$ make dependencies
```

Then use it in dev mode or run tests:
- Start server in dev mode
```
$ make run
```

- Run integration tests
```bash
$ make integration
```

## Technical Considerations

- Bifrost follows clean architecture principles
- An SDK is being developed in the `sdk/` directory
- The `pkg/` directory serves as a shared toolkit across other projects

### Architecture tradeoff
A major tradeoff has been made: since everything currently runs in a single process using a single database, all services can directly access other services' tables in the persistence layer. If a service needs to be decoupled (e.g., for scaling), these dependencies will have to be reworked—typically by introducing service-layer clients.

## Upcoming Features

- Compute results for each position
- Handle a position strategy that include multiple positions (eg. TP1, TP2...)

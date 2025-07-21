
-- +migrate Up

CREATE TABLE candles(
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  serial_id       BIGSERIAL,
  date            TIMESTAMP NOT NULL,
  pair            VARCHAR(255) NOT NULL,
  interval        VARCHAR(255) NOT NULL,
  open            DOUBLE PRECISION NOT NULL,
  close           DOUBLE PRECISION NOT NULL,
  high            DOUBLE PRECISION NOT NULL,
  low             DOUBLE PRECISION NOT NULL,
  rsi             JSONB
);


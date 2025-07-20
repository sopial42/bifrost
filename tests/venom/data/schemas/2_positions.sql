
-- +migrate Up

CREATE TABLE positions(
  id              UUID DEFAULT uuid_generate_v4(),
  serial_id       BIGSERIAL PRIMARY KEY,
  buy_signal_id   UUID NOT NULL,
  name            TEXT NOT NULL,
  fullname        TEXT NOT NULL,
  tp              DOUBLE PRECISION,
  sl              DOUBLE PRECISION,
  metadata        JSONB,
  CONSTRAINT FK_buy_signal_id FOREIGN KEY(buy_signal_id) REFERENCES buy_signals(id),
  UNIQUE (buy_signal_id, fullname)
);


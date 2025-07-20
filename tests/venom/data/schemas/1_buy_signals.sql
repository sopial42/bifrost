
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE buy_signals(
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  business_id     TEXT NOT NULL,
  pair            TEXT NOT NULL,
  interval        TEXT NOT NULL,
  name            TEXT NOT NULL,
  fullname        TEXT NOT NULL,
  "date"          TIMESTAMPTZ NOT NULL,
  price           DOUBLE PRECISION,
  metadata JSONB,
  UNIQUE  (business_id, pair, interval, fullname)
);

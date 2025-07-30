
-- +migrate Up

CREATE TABLE ratios(
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  serial_id       BIGSERIAL,
  position_id     UUID NOT NULL,
  ratio           DOUBLE PRECISION,
  date            TIMESTAMPTZ NOT NULL,
  CONSTRAINT FK_position_id FOREIGN KEY(position_id) REFERENCES positions(id),
  UNIQUE (position_id, date)
);


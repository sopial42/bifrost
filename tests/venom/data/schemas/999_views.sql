-- +migrate Up


CREATE OR REPLACE VIEW v_buy_signals_positions AS
SELECT
  bs.pair                          AS pair,
  bs.interval                      AS "buy_interval",
  bs.fullname                      AS buy_fullname,
  bs."date"                        AS buy_date,
  bs.price                         AS buy_price,
  p.fullname                       AS position_fullname,
  p.tp,
  p.sl,
  p.ratio_value,
  p.ratio_date,
  bs.metadata                      AS buy_metadata,
  p.metadata                       AS position_metadata,
  bs.id                            AS buy_id,
  p.id                             AS position_id
FROM buy_signals bs
LEFT JOIN positions p ON p.buy_signal_id = bs.id;

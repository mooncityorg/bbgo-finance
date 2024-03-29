package sqlite3

import (
	"context"

	"github.com/c9s/rockhopper"
)

func init() {
	AddMigration(upAddFuturesKlines, downAddFuturesKlines)

}

func upAddFuturesKlines(ctx context.Context, tx rockhopper.SQLExecutor) (err error) {
	// This code is executed when the migration is applied.

	_, err = tx.ExecContext(ctx, "CREATE TABLE `bybit_futures_klines`\n(\n    `gid`                    INTEGER PRIMARY KEY AUTOINCREMENT,\n    `exchange`               VARCHAR(10)    NOT NULL,\n    `start_time`             DATETIME(3)    NOT NULL,\n    `end_time`               DATETIME(3)    NOT NULL,\n    `interval`               VARCHAR(3)     NOT NULL,\n    `symbol`                 VARCHAR(7)     NOT NULL,\n    `open`                   DECIMAL(16, 8) NOT NULL,\n    `high`                   DECIMAL(16, 8) NOT NULL,\n    `low`                    DECIMAL(16, 8) NOT NULL,\n    `close`                  DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `volume`                 DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `closed`                 BOOLEAN        NOT NULL DEFAULT TRUE,\n    `last_trade_id`          INT            NOT NULL DEFAULT 0,\n    `num_trades`             INT            NOT NULL DEFAULT 0,\n    `quote_volume`           DECIMAL        NOT NULL DEFAULT 0.0,\n    `taker_buy_base_volume`  DECIMAL        NOT NULL DEFAULT 0.0,\n    `taker_buy_quote_volume` DECIMAL        NOT NULL DEFAULT 0.0\n);")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "CREATE TABLE `okex_futures_klines`\n(\n    `gid`           INTEGER PRIMARY KEY AUTOINCREMENT,\n    `exchange`      VARCHAR(10)    NOT NULL,\n    `start_time`    DATETIME(3)    NOT NULL,\n    `end_time`      DATETIME(3)    NOT NULL,\n    `interval`      VARCHAR(3)     NOT NULL,\n    `symbol`        VARCHAR(7)     NOT NULL,\n    `open`          DECIMAL(16, 8) NOT NULL,\n    `high`          DECIMAL(16, 8) NOT NULL,\n    `low`           DECIMAL(16, 8) NOT NULL,\n    `close`         DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `volume`        DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `closed`        BOOLEAN        NOT NULL DEFAULT TRUE,\n    `last_trade_id` INT            NOT NULL DEFAULT 0,\n    `num_trades`    INT            NOT NULL DEFAULT 0\n);")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "CREATE TABLE `binance_futures_klines`\n(\n    `gid`           INTEGER PRIMARY KEY AUTOINCREMENT,\n    `exchange`      VARCHAR(10)    NOT NULL,\n    `start_time`    DATETIME(3)    NOT NULL,\n    `end_time`      DATETIME(3)    NOT NULL,\n    `interval`      VARCHAR(3)     NOT NULL,\n    `symbol`        VARCHAR(7)     NOT NULL,\n    `open`          DECIMAL(16, 8) NOT NULL,\n    `high`          DECIMAL(16, 8) NOT NULL,\n    `low`           DECIMAL(16, 8) NOT NULL,\n    `close`         DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `volume`        DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `closed`        BOOLEAN        NOT NULL DEFAULT TRUE,\n    `last_trade_id` INT            NOT NULL DEFAULT 0,\n    `num_trades`    INT            NOT NULL DEFAULT 0\n);")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "CREATE TABLE `max_futures_klines`\n(\n    `gid`           INTEGER PRIMARY KEY AUTOINCREMENT,\n    `exchange`      VARCHAR(10)    NOT NULL,\n    `start_time`    DATETIME(3)    NOT NULL,\n    `end_time`      DATETIME(3)    NOT NULL,\n    `interval`      VARCHAR(3)     NOT NULL,\n    `symbol`        VARCHAR(7)     NOT NULL,\n    `open`          DECIMAL(16, 8) NOT NULL,\n    `high`          DECIMAL(16, 8) NOT NULL,\n    `low`           DECIMAL(16, 8) NOT NULL,\n    `close`         DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `volume`        DECIMAL(16, 8) NOT NULL DEFAULT 0.0,\n    `closed`        BOOLEAN        NOT NULL DEFAULT TRUE,\n    `last_trade_id` INT            NOT NULL DEFAULT 0,\n    `num_trades`    INT            NOT NULL DEFAULT 0\n);")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "CREATE INDEX `bybit_futures_klines_end_time_symbol_interval` ON `bybit_futures_klines` (`end_time`, `symbol`, `interval`);\nCREATE INDEX `okex_futures_klines_end_time_symbol_interval` ON `okex_futures_klines` (`end_time`, `symbol`, `interval`);\nCREATE INDEX `binance_futures_klines_end_time_symbol_interval` ON `binance_futures_klines` (`end_time`, `symbol`, `interval`);\nCREATE INDEX `max_futures_klines_end_time_symbol_interval` ON `max_futures_klines` (`end_time`, `symbol`, `interval`);")
	if err != nil {
		return err
	}

	return err
}

func downAddFuturesKlines(ctx context.Context, tx rockhopper.SQLExecutor) (err error) {
	// This code is executed when the migration is rolled back.

	_, err = tx.ExecContext(ctx, "DROP INDEX IF EXISTS `bybit_futures_klines_end_time_symbol_interval`;\nDROP INDEX IF EXISTS `okex_futures_klines_end_time_symbol_interval`;\nDROP INDEX IF EXISTS `binance_futures_klines_end_time_symbol_interval`;\nDROP INDEX IF EXISTS `max_futures_klines_end_time_symbol_interval`;\n")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DROP TABLE IF EXISTS `bybit_futures_klines`;")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DROP TABLE IF EXISTS `okex_futures_klines`;")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DROP TABLE IF EXISTS `binance_futures_klines`;")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DROP TABLE IF EXISTS `max_futures_klines`;")
	if err != nil {
		return err
	}

	return err
}

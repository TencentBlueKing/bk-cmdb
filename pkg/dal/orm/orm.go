/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package orm ...
package orm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// DefaultSlowSQLThreshold default slow sql threshold
const DefaultSlowSQLThreshold = 200 * time.Millisecond

// Interface defines all the orm related operations.
type Interface interface {
	WithTx(tx *gorm.DB) Interface
	AutoTxn(ctx context.Context, run TxnFunc) (any, error)
	DB(opts ...func(*gorm.DB) *gorm.DB) *gorm.DB
}

// WithContext set context for db
func WithContext(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx)
	}
}

// UseWriteDB use write db
func UseWriteDB() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(dbresolver.Write)
	}
}

// PreferReadDB prefer read db
func PreferReadDB() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(dbresolver.Read)
	}
}

// WithSessionConfig set session for db
func WithSessionConfig(config *gorm.Session) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Session(config)
	}
}

// New return orm operations.
func New(ctx context.Context, db *gorm.DB, opts ...Option) (Interface, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	ormOpts := new(options)
	for _, opt := range opts {
		opt(ormOpts)
	}

	ro := &runtimeOrm{db: db}

	// initialize metrics with default registerer if not set
	if ormOpts.mc == nil {
		ormOpts.mc = initMetric(prometheus.DefaultRegisterer)
	}

	if ormOpts.slowRequestTime == 0 {
		ormOpts.slowRequestTime = DefaultSlowSQLThreshold
	}

	loggerLevel := logger.Error
	if ormOpts.debug {
		loggerLevel = logger.Info
	}

	loggerConfig := logger.Config{
		SlowThreshold:             ormOpts.slowRequestTime,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
		LogLevel:                  loggerLevel,
	}
	slogger := slog.New(log.Depth(2).Handler())
	// reset logger to use slog logger
	db.Logger = logger.NewSlogLogger(slogger, loggerConfig)

	ro.mc = ormOpts.mc
	ro.slowRequestTime = ormOpts.slowRequestTime

	// register metrics callbacks to record the cost time to exec an orm command.
	if err := ro.registerMetricsCallbacks(); err != nil {
		log.Error(ctx, "register metrics callbacks error", log.E(err))
		return nil, err
	}
	return ro, nil
}

func (o *runtimeOrm) registerMetricsCallbacks() error {
	// register for create command
	err := o.db.Callback().Create().Before("*").Register("metric:create:before", o.beforeCommandCallback("create"))
	if err != nil {
		return fmt.Errorf("register create_metric error: %w", err)
	}
	err = o.db.Callback().Create().After("*").Register("metric:create:after", o.afterCmdCallback("create"))
	if err != nil {
		return fmt.Errorf("register metric:create:after error: %w", err)
	}

	// register for update command
	err = o.db.Callback().Update().Before("*").Register("metric:update:before", o.beforeCommandCallback("update"))
	if err != nil {
		return fmt.Errorf("register metric:update:before error: %w", err)
	}
	err = o.db.Callback().Update().After("*").Register("metric:update:after", o.afterCmdCallback("update"))
	if err != nil {
		return fmt.Errorf("register metric:update:after error: %w", err)
	}

	// register for query command
	err = o.db.Callback().Query().Before("*").Register("metric:query:before", o.beforeCommandCallback("query"))
	if err != nil {
		return fmt.Errorf("register metric:query:before error: %w", err)
	}
	err = o.db.Callback().Query().After("*").Register("metric:query:after", o.afterCmdCallback("query"))
	if err != nil {
		return fmt.Errorf("register metric:query:after error: %w", err)
	}

	// register for delete command
	err = o.db.Callback().Delete().Before("*").Register("metric:delete:before", o.beforeCommandCallback("delete"))
	if err != nil {
		return fmt.Errorf("register metric:delete:before error: %w", err)
	}
	err = o.db.Callback().Delete().After("*").Register("metric:delete:after", o.afterCmdCallback("delete"))
	if err != nil {
		return fmt.Errorf("register metric:delete:after error: %w", err)
	}

	// register for row command
	err = o.db.Callback().Row().Before("*").Register("metric:row:before", o.beforeCommandCallback("row"))
	if err != nil {
		return fmt.Errorf("register metric:row:before error: %w", err)
	}
	err = o.db.Callback().Row().After("*").Register("metric:row:after", o.afterCmdCallback("row"))
	if err != nil {
		return fmt.Errorf("register metric:row:after error: %w", err)
	}

	// register for raw command
	err = o.db.Callback().Raw().Before("*").Register("metric:raw:before", o.beforeCommandCallback("raw"))
	if err != nil {
		return fmt.Errorf("register metric:raw:before error: %w", err)
	}
	err = o.db.Callback().Raw().After("*").Register("metric:raw:after", o.afterCmdCallback("raw"))
	if err != nil {
		return fmt.Errorf("register metric:raw:after error: %w", err)
	}

	return nil
}

type runtimeOrm struct {
	db              *gorm.DB
	mc              *metric
	slowRequestTime time.Duration
}

// WithTx returns a transactional database.
func (o *runtimeOrm) WithTx(tx *gorm.DB) Interface {
	return &runtimeOrm{
		db:              tx,
		mc:              o.mc,
		slowRequestTime: o.slowRequestTime,
	}
}

// DB returns a database with options.
func (o *runtimeOrm) DB(opts ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	db := o.db
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

const metricStartTimeKey = "cmd_start_time"

// beforeCommandCallback record the start time of an orm command.
func (o *runtimeOrm) beforeCommandCallback(cmd string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		db.Set(metricStartTimeKey, time.Now())
	}
}

// afterCmdCallback record the cost time to exec an orm command.
func (o *runtimeOrm) afterCmdCallback(cmd string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Error != nil {
			o.mc.errCounter.With(prometheus.Labels{"cmd": cmd}).Inc()
			return
		}
		startTime := db.Statement.Context.Value(metricStartTimeKey)
		if startTime == nil {
			return
		}
		if _, ok := startTime.(time.Time); ok {
			latency := time.Since(startTime.(time.Time))
			o.mc.cmdLagMS.With(prometheus.Labels{"cmd": cmd}).Observe(float64(latency.Milliseconds()))
		}
	}
}

// TxnFunc is a callback function to process logic tasks between a transaction.
type TxnFunc func(txn *gorm.DB) (any, error)

// ErrRetryTransaction defines errors that need to retry transaction, like deadlock error in upsert scenario,
// could be defined by user.
var ErrRetryTransaction = errors.New("RETRY TRANSACTION ERROR")

// AutoTxn is a wrapper to do all the transaction operations as follows:
// 1. auto launch the transaction
// 2. process the logics, which is a callback run function
// 3. rollback the transaction if 'run' hit an error automatically.
// 4. commit the transaction if no error happens.
func (o *runtimeOrm) AutoTxn(ctx context.Context, run TxnFunc) (any, error) {
	if run == nil {
		return nil, errors.New("transaction function is nil")
	}

	retry, result, err := o.autoTxn(ctx, run)
	if err == nil {
		return result, nil
	}

	if !retry {
		return nil, err
	}

	// if the operation need to retry, retry for at most 3 times, each wait for 50~500ms
	for retryCount := 1; retryCount <= 3; retryCount++ {
		log.Warn(ctx, "retry transaction", "retry_count", retryCount)
		time.Sleep(time.Millisecond * time.Duration(rand.IntN(450)+50))

		retry, result, err = o.autoTxn(ctx, run)
		if err == nil {
			return result, nil
		}

		if !retry {
			return nil, err
		}

		// do next retry
	}

	log.Warn(ctx, "retry transaction exceeds maximum count, **skip**")
	return nil, err
}

func (o *runtimeOrm) autoTxn(ctx context.Context, action TxnFunc) (retry bool, result any, err error) {
	if action == nil {
		return false, nil, errors.New("transaction function is nil")
	}
	err = o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result, err = action(tx)
		return err
	})
	if errors.Is(err, ErrRetryTransaction) {
		return true, nil, err
	}
	return false, result, nil
}

// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/session"
	"github.com/mongodb/mongo-go-driver/x/network/connection"
)

var defaultRegistry = bson.NewRegistryBuilder().Build()

type serverConfig struct {
	clock             *session.ClusterClock
	compressionOpts   []string
	connectionOpts    []connection.Option
	appname           string
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration
	maxConns          uint16
	maxIdleConns      uint16
	registry          *bsoncodec.Registry
}

func newServerConfig(opts ...ServerOption) (*serverConfig, error) {
	cfg := &serverConfig{
		heartbeatInterval: 10 * time.Second,
		heartbeatTimeout:  10 * time.Second,
		maxConns:          100,
		maxIdleConns:      100,
		registry:          defaultRegistry,
	}

	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// ServerOption configures a server.
type ServerOption func(*serverConfig) error

// WithConnectionOptions configures the server's connections.
func WithConnectionOptions(fn func(...connection.Option) []connection.Option) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.connectionOpts = fn(cfg.connectionOpts...)
		return nil
	}
}

// WithCompressionOptions configures the server's compressors.
func WithCompressionOptions(fn func(...string) []string) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.compressionOpts = fn(cfg.compressionOpts...)
		return nil
	}
}

// WithHeartbeatInterval configures a server's heartbeat interval.
func WithHeartbeatInterval(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.heartbeatInterval = fn(cfg.heartbeatInterval)
		return nil
	}
}

// WithHeartbeatTimeout configures how long to wait for a heartbeat socket to
// connection.
func WithHeartbeatTimeout(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.heartbeatTimeout = fn(cfg.heartbeatTimeout)
		return nil
	}
}

// WithMaxConnections configures the maximum number of connections to allow for
// a given server. If max is 0, then there is no upper limit to the number of
// connections.
func WithMaxConnections(fn func(uint16) uint16) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.maxConns = fn(cfg.maxConns)
		return nil
	}
}

// WithMaxIdleConnections configures the maximum number of idle connections
// allowed for the server.
func WithMaxIdleConnections(fn func(uint16) uint16) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.maxIdleConns = fn(cfg.maxIdleConns)
		return nil
	}
}

// WithClock configures the ClusterClock for the server to use.
func WithClock(fn func(clock *session.ClusterClock) *session.ClusterClock) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.clock = fn(cfg.clock)
		return nil
	}
}

// WithRegistry configures the registry for the server to use when creating
// cursors.
func WithRegistry(fn func(*bsoncodec.Registry) *bsoncodec.Registry) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.registry = fn(cfg.registry)
		return nil
	}
}

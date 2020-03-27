// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/rentiansheng/bk_bsonx/mongo/driver/auth"
	"github.com/rentiansheng/bk_bsonx/network/address"
	"github.com/rentiansheng/bk_bsonx/network/connection"
	"github.com/rentiansheng/bk_bsonx/network/description"
	"github.com/stretchr/testify/require"
)

type pool struct {
	connectionError bool
	drainCalled     atomic.Value
	networkError    bool
	desc            *description.Server
}

func (p *pool) Get(ctx context.Context) (connection.Connection, *description.Server, error) {
	if p.connectionError {
		return nil, p.desc, &auth.Error{}
	}
	if p.networkError {
		return nil, p.desc, &connection.NetworkError{}
	}
	return nil, p.desc, nil
}

func (p *pool) Connect(ctx context.Context) error {
	return nil
}

func (p *pool) Disconnect(ctx context.Context) error {
	return nil
}

func (p *pool) Drain() error {
	p.drainCalled.Store(true)
	return nil
}

func NewPool(connectionError bool, networkError bool, desc *description.Server) (connection.Pool, error) {
	p := &pool{
		connectionError: connectionError,
		networkError:    networkError,
		desc:            desc,
	}
	p.drainCalled.Store(false)
	return p, nil
}

func TestSever(t *testing.T) {
	var serverTestTable = []struct {
		name            string
		connectionError bool
		networkError    bool
		hasDesc         bool
	}{
		{"auth_error", true, false, false},
		{"no_error", false, false, false},
		{"network_error_no_desc", false, true, false},
		{"network_error_desc", false, true, true},
	}

	for _, tt := range serverTestTable {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewServer(address.Address("localhost"))
			require.NoError(t, err)

			var desc *description.Server
			descript := s.Description()
			if tt.hasDesc {
				desc = &descript
				require.Nil(t, desc.LastError)
			}
			s.pool, err = NewPool(tt.connectionError, tt.networkError, desc)
			s.connectionstate = connected

			_, err = s.Connection(context.Background())

			if tt.connectionError || tt.networkError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.hasDesc {
				require.Equal(t, desc.Kind, (description.ServerKind)(description.Unknown))
				require.NotNil(t, desc.LastError)
			}
			drained := s.pool.(*pool).drainCalled.Load().(bool)
			require.Equal(t, drained, tt.connectionError || tt.networkError)
		})
	}
}

// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"github.com/mongodb/mongo-go-driver/mongo/readconcern"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/mongodb/mongo-go-driver/mongo/writeconcern"
)

// DefaultCausalConsistency is the default value for the CausalConsistency option.
var DefaultCausalConsistency = true

// SessionOptions represents all possible options for creating a new session.
type SessionOptions struct {
	CausalConsistency     *bool                      // Specifies if reads should be causally consistent. Defaults to true.
	DefaultReadConcern    *readconcern.ReadConcern   // The default read concern for transactions started in the session.
	DefaultReadPreference *readpref.ReadPref         // The default read preference for transactions started in the session.
	DefaultWriteConcern   *writeconcern.WriteConcern // The default write concern for transactions started in the session.
}

// Session creates a new *SessionOptions
func Session() *SessionOptions {
	return &SessionOptions{
		CausalConsistency: &DefaultCausalConsistency,
	}
}

// SetCausalConsistency specifies if a session should be causally consistent. Defaults to true.
func (s *SessionOptions) SetCausalConsistency(b bool) *SessionOptions {
	s.CausalConsistency = &b
	return s
}

// SetDefaultReadConcern sets the default read concern for transactions started in a session.
func (s *SessionOptions) SetDefaultReadConcern(rc *readconcern.ReadConcern) *SessionOptions {
	s.DefaultReadConcern = rc
	return s
}

// SetDefaultReadPreference sets the default read preference for transactions started in a session.
func (s *SessionOptions) SetDefaultReadPreference(rp *readpref.ReadPref) *SessionOptions {
	s.DefaultReadPreference = rp
	return s
}

// SetDefaultWriteConcern sets the default write concern for transactions started in a session.
func (s *SessionOptions) SetDefaultWriteConcern(wc *writeconcern.WriteConcern) *SessionOptions {
	s.DefaultWriteConcern = wc
	return s
}

// MergeSessionOptions combines the given *SessionOptions into a single *SessionOptions in a last one wins fashion.
func MergeSessionOptions(opts ...*SessionOptions) *SessionOptions {
	s := Session()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.CausalConsistency != nil {
			s.CausalConsistency = opt.CausalConsistency
		}
		if opt.DefaultReadConcern != nil {
			s.DefaultReadConcern = opt.DefaultReadConcern
		}
		if opt.DefaultReadPreference != nil {
			s.DefaultReadPreference = opt.DefaultReadPreference
		}
		if opt.DefaultWriteConcern != nil {
			s.DefaultWriteConcern = opt.DefaultWriteConcern
		}
	}

	return s
}

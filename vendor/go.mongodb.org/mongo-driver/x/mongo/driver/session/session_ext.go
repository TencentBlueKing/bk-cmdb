// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session // import "go.mongodb.org/mongo-driver/x/mongo/driver/session"

// GetState get the state of the client session
func (c *Client) GetState() uint8 {
	return uint8(c.state)
}

// SetState set the state of the client session
func (c *Client) SetState(s uint8) {
	c.state = state(s)
}

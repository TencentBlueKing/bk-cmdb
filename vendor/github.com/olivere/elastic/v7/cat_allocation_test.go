// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatAllocation(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatAllocation().Columns("*").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}
	if have := res[0].IP; have == "" {
		t.Fatalf("IP[0]: want != %q, have %q", "", have)
	}
}

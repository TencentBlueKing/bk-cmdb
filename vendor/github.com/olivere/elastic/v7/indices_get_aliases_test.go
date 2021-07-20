// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestAliasesBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Indices  []string
		Aliases  []string
		Expected string
	}{
		{
			[]string{},
			[]string{},
			"/_alias",
		},
		{
			[]string{"index1"},
			[]string{},
			"/index1/_alias",
		},
		{
			[]string{"index1", "index2"},
			[]string{},
			"/index1%2Cindex2/_alias",
		},
		{
			[]string{"index1", "index2"},
			[]string{"alias1"},
			"/index1%2Cindex2/_alias/alias1",
		},
		{
			[]string{},
			[]string{"alias1"},
			"/_alias/alias1",
		},
		{
			[]string{"index1"},
			[]string{"alias1"},
			"/index1/_alias/alias1",
		},
		{
			[]string{"index1"},
			[]string{"alias1", "alias2"},
			"/index1/_alias/alias1%2Calias2",
		},
		{
			[]string{},
			[]string{"alias1", "alias2"},
			"/_alias/alias1%2Calias2",
		},
		{
			[]string{"index1", "index2"},
			[]string{"alias1", "alias2"},
			"/index1%2Cindex2/_alias/alias1%2Calias2",
		},
	}

	for i, test := range tests {
		path, _, err := client.Aliases().Index(test.Indices...).Alias(test.Aliases...).buildURL()
		if err != nil {
			t.Errorf("case #%d: %v", i+1, err)
			continue
		}
		if path != test.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
		}
	}
}

func TestAliases(t *testing.T) {
	var err error

	// client := setupTestClientAndCreateIndex(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	// Some tweets
	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "sandrae", Message: "Cycling is fun."}
	tweet3 := tweet{User: "olivere", Message: "Another unrelated topic."}

	// Add tweets to first index
	_, err = client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	// Add tweets to second index
	_, err = client.Index().Index(testIndexName2).Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Refresh
	_, err = client.Refresh().Index(testIndexName, testIndexName2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Alias should not yet exist
	aliasesResult1, err := client.Aliases().
		Index(testIndexName, testIndexName2).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if len(aliasesResult1.Indices) != 2 {
		t.Errorf("expected len(AliasesResult.Indices) = %d; got %d", 2, len(aliasesResult1.Indices))
	}
	for indexName, indexDetails := range aliasesResult1.Indices {
		if len(indexDetails.Aliases) != 0 {
			t.Errorf("expected len(AliasesResult.Indices[%s].Aliases) = %d; got %d", indexName, 0, len(indexDetails.Aliases))
		}
	}

	// Add both indices to a new alias
	aliasCreate, err := client.Alias().
		Add(testIndexName, testAliasName).
		Action(
			NewAliasAddAction(testAliasName).
				Index(testIndexName2).
				IsWriteIndex(true),
		).
		//Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !aliasCreate.Acknowledged {
		t.Errorf("expected AliasResult.Acknowledged %v; got %v", true, aliasCreate.Acknowledged)
	}

	// Alias should now exist
	aliasesResult2, err := client.Aliases().
		Index(testIndexName, testIndexName2).
		//Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if len(aliasesResult2.Indices) != 2 {
		t.Errorf("expected len(AliasesResult.Indices) = %d; got %d", 2, len(aliasesResult2.Indices))
	}
	for indexName, indexDetails := range aliasesResult2.Indices {
		if len(indexDetails.Aliases) != 1 {
			t.Errorf("expected len(AliasesResult.Indices[%s].Aliases) = %d; got %d", indexName, 1, len(indexDetails.Aliases))
		}
	}
	indicesResult, err := client.Aliases().
		Alias(testAliasName).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if len(indicesResult.Indices) != 2 {
		t.Errorf("expected len(indicesResult.Indices) = %d; got %d", 2, len(indicesResult.Indices))
	}
	for indexName, indexDetails := range indicesResult.Indices {
		if len(indexDetails.Aliases) != 1 {
			t.Errorf("expected len(indicesResult.Indices[%s].Aliases) = %d; got %d", indexName, 1, len(indexDetails.Aliases))
		}
		if indexName == testIndexName2 {
			if !indexDetails.Aliases[0].IsWriteIndex {
				t.Errorf("expected alias on %s to be a write index", testIndexName2)
			}
		}
		if indexName == testIndexName {
			if indexDetails.Aliases[0].IsWriteIndex {
				t.Errorf("expected alias on %s not to be a write index", testIndexName2)
			}
		}
	}

	// Check the reverse function:
	indexInfo1, found := aliasesResult2.Indices[testIndexName]
	if !found {
		t.Errorf("expected info about index %s = %v; got %v", testIndexName, true, found)
	}
	aliasFound := indexInfo1.HasAlias(testAliasName)
	if !aliasFound {
		t.Errorf("expected alias %s to include index %s; got %v", testAliasName, testIndexName, aliasFound)
	}

	// Check the reverse function:
	indexInfo2, found := aliasesResult2.Indices[testIndexName2]
	if !found {
		t.Errorf("expected info about index %s = %v; got %v", testIndexName, true, found)
	}
	aliasFound = indexInfo2.HasAlias(testAliasName)
	if !aliasFound {
		t.Errorf("expected alias %s to include index %s; got %v", testAliasName, testIndexName2, aliasFound)
	}

	// Remove first index should remove two tweets, so should only yield 1
	aliasRemove1, err := client.Alias().
		Remove(testIndexName, testAliasName).
		//Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !aliasRemove1.Acknowledged {
		t.Errorf("expected AliasResult.Acknowledged %v; got %v", true, aliasRemove1.Acknowledged)
	}

	// Alias should now exist only for index 2
	aliasesResult3, err := client.Aliases().Index(testIndexName, testIndexName2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if len(aliasesResult3.Indices) != 2 {
		t.Errorf("expected len(AliasesResult.Indices) = %d; got %d", 2, len(aliasesResult3.Indices))
	}
	for indexName, indexDetails := range aliasesResult3.Indices {
		if indexName == testIndexName {
			if len(indexDetails.Aliases) != 0 {
				t.Errorf("expected len(AliasesResult.Indices[%s].Aliases) = %d; got %d", indexName, 0, len(indexDetails.Aliases))
			}
		} else if indexName == testIndexName2 {
			if len(indexDetails.Aliases) != 1 {
				t.Errorf("expected len(AliasesResult.Indices[%s].Aliases) = %d; got %d", indexName, 1, len(indexDetails.Aliases))
			}
		} else {
			t.Errorf("got index %s", indexName)
		}
	}
}

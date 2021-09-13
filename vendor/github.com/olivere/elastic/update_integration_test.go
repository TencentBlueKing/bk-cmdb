// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestUpdateWithDoc(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	// Get original
	getRes, err := client.Get().Index(testIndexName).Id("1").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	var original tweet
	if err := json.Unmarshal(getRes.Source, &original); err != nil {
		t.Fatal(err)
	}

	// Partial update
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	updRes, err := client.Update().
		Index(testIndexName).
		Id("1").
		Doc(map[string]interface{}{
			"message": "Updated message text.",
		}).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if updRes == nil {
		t.Fatal("response is nil")
	}
	if want, have := "updated", updRes.Result; want != have {
		t.Fatalf("want Result = %q, have %v", want, have)
	}

	// Get new version
	getRes, err = client.Get().Index(testIndexName).Id("1").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	var updated tweet
	if err := json.Unmarshal(getRes.Source, &updated); err != nil {
		t.Fatal(err)
	}

	if want, have := original.User, updated.User; want != have {
		t.Fatalf("want User = %q, have %v", want, have)
	}
	if want, have := "Updated message text.", updated.Message; want != have {
		t.Fatalf("want Message = %q, have %v", want, have)
	}
}

func TestUpdateWithScript(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	// Get original
	getRes, err := client.Get().Index(testIndexName).Id("1").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	var original tweet
	if err := json.Unmarshal(getRes.Source, &original); err != nil {
		t.Fatal(err)
	}

	// Update with script
	updRes, err := client.Update().Index(testIndexName).Id("1").
		Script(
			NewScript(`ctx._source.message = "Updated message text."`).Lang("painless"),
		).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if updRes == nil {
		t.Fatal("response is nil")
	}
	if want, have := "updated", updRes.Result; want != have {
		t.Fatalf("want Result = %q, have %v", want, have)
	}

	// Get new version
	getRes, err = client.Get().Index(testIndexName).Id("1").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	var updated tweet
	if err := json.Unmarshal(getRes.Source, &updated); err != nil {
		t.Fatal(err)
	}

	if want, have := original.User, updated.User; want != have {
		t.Fatalf("want User = %q, have %v", want, have)
	}
	if want, have := "Updated message text.", updated.Message; want != have {
		t.Fatalf("want Message = %q, have %v", want, have)
	}
}

func TestUpdateWithScriptID(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Get original
	getRes, err := client.Get().Index(testIndexName).Id("1").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	var original tweet
	if err := json.Unmarshal(getRes.Source, &original); err != nil {
		t.Fatal(err)
	}

	// Set script with ID
	scriptID := "example-script-id"
	_, err = client.DeleteScript().Id(scriptID).Do(context.Background())
	if err != nil && !IsNotFound(err) {
		t.Fatal(err)
	}
	_, err = client.PutScript().
		Id(scriptID).
		BodyString(`{
			"script": {
				"lang": "painless",
				"source": "ctx._source.message = params.new_message"
			}
		}`).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Update with script
	updRes, err := client.Update().Index(testIndexName).Id("1").
		Script(
			NewScriptStored(scriptID).Param("new_message", "Updated message text."),
		).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if updRes == nil {
		t.Fatal("response is nil")
	}
	if want, have := "updated", updRes.Result; want != have {
		t.Fatalf("want Result = %q, have %v", want, have)
	}

	// Get new version
	getRes, err = client.Get().Index(testIndexName).Id("1").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	var updated tweet
	if err := json.Unmarshal(getRes.Source, &updated); err != nil {
		t.Fatal(err)
	}

	if want, have := original.User, updated.User; want != have {
		t.Fatalf("want User = %q, have %v", want, have)
	}
	if want, have := "Updated message text.", updated.Message; want != have {
		t.Fatalf("want Message = %q, have %v", want, have)
	}
}

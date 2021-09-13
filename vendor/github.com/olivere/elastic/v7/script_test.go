// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestScriptingDefault(t *testing.T) {
	builder := NewScript("doc['field'].value * 2")
	src, err := builder.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"source":"doc['field'].value * 2"}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestScriptingInline(t *testing.T) {
	builder := NewScriptInline("doc['field'].value * factor").Param("factor", 2.0)
	src, err := builder.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"params":{"factor":2},"source":"doc['field'].value * factor"}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestScriptingStored(t *testing.T) {
	builder := NewScriptStored("script-with-id").Param("factor", 2.0)
	src, err := builder.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"id":"script-with-id","params":{"factor":2}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestScriptingSource(t *testing.T) {
	tests := []struct {
		Input  string
		Source string
	}{
		{
			Input:  ``,
			Source: `{"source":""}`,
		},
		{
			Input:  `doc['field'].value * factor`,
			Source: `{"source":"doc['field'].value * factor"}`,
		},
		{
			Input:  `"doc['field'].value * factor"`,
			Source: `{"source":"doc['field'].value * factor"}`,
		},
		{
			Input:  `{"bool":{"filter":{"term":{"field1":"f"}}}}`,
			Source: `{"source":{"bool":{"filter":{"term":{"field1":"f"}}}}}`,
		},
	}
	for _, tt := range tests {
		b := NewScriptInline(tt.Input)
		src, err := b.Source()
		if err != nil {
			t.Fatalf("unable to generate source for %s: %v", tt.Input, err)
		}
		out, err := json.Marshal(src)
		if err != nil {
			t.Fatalf("unable to generate JSON for %s: %v", tt.Input, err)
		}
		if want, have := tt.Source, string(out); want != have {
			t.Fatalf("Input=%q: want %s, have %s", tt.Input, want, have)
		}
	}
}

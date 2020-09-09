package query

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
)

func assertArrayContainsInOrder(t *testing.T, array []interface{}, objects ...interface{}) {
	if len(array) != len(objects) {
		t.Fatalf("array contains %d objects but %d are expected", len(array), len(objects))
	}

	for i := 0; i < len(array); i++ {
		if array[i] != objects[i] {
			t.Fatalf("wanted '%s', have '%s'", objects[i], array[i])
		}
	}
}

func checkQuery(t *testing.T, tree *toml.Tree, query string, objects ...interface{}) {
	results, err := CompileAndExecute(query, tree)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	assertArrayContainsInOrder(t, results.Values(), objects...)
}

func TestQueryExample(t *testing.T) {
	config, _ := toml.Load(`
      [[book]]
      title = "The Stand"
      author = "Stephen King"
      [[book]]
      title = "For Whom the Bell Tolls"
      author = "Ernest Hemmingway"
      [[book]]
      title = "Neuromancer"
      author = "William Gibson"
	`)

	checkQuery(t, config, "$.book.author", "Stephen King", "Ernest Hemmingway", "William Gibson")

	checkQuery(t, config, "$.book[0].author", "Stephen King")
	checkQuery(t, config, "$.book[-1].author", "William Gibson")
	checkQuery(t, config, "$.book[1:].author", "Ernest Hemmingway", "William Gibson")
	checkQuery(t, config, "$.book[-1:].author", "William Gibson")
	checkQuery(t, config, "$.book[::2].author", "Stephen King", "William Gibson")
	checkQuery(t, config, "$.book[::-1].author", "William Gibson", "Ernest Hemmingway", "Stephen King")
	checkQuery(t, config, "$.book[:].author", "Stephen King", "Ernest Hemmingway", "William Gibson")
	checkQuery(t, config, "$.book[::].author", "Stephen King", "Ernest Hemmingway", "William Gibson")
}

func TestQueryReadmeExample(t *testing.T) {
	config, _ := toml.Load(`
[postgres]
user = "pelletier"
password = "mypassword"
`)

	checkQuery(t, config, "$..[user,password]", "pelletier", "mypassword")
}

func TestQueryPathNotPresent(t *testing.T) {
	config, _ := toml.Load(`a = "hello"`)
	query, err := Compile("$.foo.bar")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	results := query.Execute(config)
	if err != nil {
		t.Fatalf("err should be nil. got %s instead", err)
	}
	if len(results.items) != 0 {
		t.Fatalf("no items should be matched. %d matched instead", len(results.items))
	}
}

func ExampleNodeFilterFn_filterExample() {
	tree, _ := toml.Load(`
      [struct_one]
      foo = "foo"
      bar = "bar"

      [struct_two]
      baz = "baz"
      gorf = "gorf"
    `)

	// create a query that references a user-defined-filter
	query, _ := Compile("$[?(bazOnly)]")

	// define the filter, and assign it to the query
	query.SetFilter("bazOnly", func(node interface{}) bool {
		if tree, ok := node.(*toml.Tree); ok {
			return tree.Has("baz")
		}
		return false // reject all other node types
	})

	// results contain only the 'struct_two' Tree
	query.Execute(tree)
}

func ExampleQuery_queryExample() {
	config, _ := toml.Load(`
      [[book]]
      title = "The Stand"
      author = "Stephen King"
      [[book]]
      title = "For Whom the Bell Tolls"
      author = "Ernest Hemmingway"
      [[book]]
      title = "Neuromancer"
      author = "William Gibson"
    `)

	// find and print all the authors in the document
	query, _ := Compile("$.book.author")
	authors := query.Execute(config)
	for _, name := range authors.Values() {
		fmt.Println(name)
	}
}

func TestTomlQuery(t *testing.T) {
	tree, err := toml.Load("[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6")
	if err != nil {
		t.Error(err)
		return
	}
	query, err := Compile("$.foo.bar")
	if err != nil {
		t.Error(err)
		return
	}
	result := query.Execute(tree)
	values := result.Values()
	if len(values) != 1 {
		t.Errorf("Expected resultset of 1, got %d instead: %v", len(values), values)
	}

	if tt, ok := values[0].(*toml.Tree); !ok {
		t.Errorf("Expected type of Tree: %T", values[0])
	} else if tt.Get("a") != int64(1) {
		t.Errorf("Expected 'a' with a value 1: %v", tt.Get("a"))
	} else if tt.Get("b") != int64(2) {
		t.Errorf("Expected 'b' with a value 2: %v", tt.Get("b"))
	}
}

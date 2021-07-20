// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect creates an index with a mapping with different data types.
//
// Example
//
//     mapping -url=http://127.0.0.1:9200 -index=twitter
//
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/olivere/elastic/v7"
)

const (
	mapping = `
	{
		"settings":{
			"number_of_shards": {{.NumberOfShards}},
			"number_of_replicas": {{.NumberOfReplicas}}
		},
		"mappings":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text"
				},
				"retweets":{
					"type":"integer"
				},
				"created":{
					"type":"date"
				},
				"attributes":{
					"type":"object"
				}
			}
		}
	}
	`
)

// Tweet is just an example document.
type Tweet struct {
	User     string                 `json:"user"`
	Message  string                 `json:"message"`
	Retweets int                    `json:"retweets"`
	Created  time.Time              `json:"created"`
	Attrs    map[string]interface{} `json:"attributes,omitempty"`
}

func main() {
	var (
		url      = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff    = flag.Bool("sniff", true, "Enable or disable sniffing")
		index    = flag.String("index", "", "Index name")
		shards   = flag.Int("shards", 1, "Number of shards")
		replicas = flag.Int("replicas", 0, "Number of replicas")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}
	if *index == "" {
		log.Fatal("please specify an index name -index")
	}

	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff))
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	// Check if index already exists. We'll drop it then.
	// Next, we create a fresh index/mapping.
	ctx := context.Background()
	exists, err := client.IndexExists(*index).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		_, err := client.DeleteIndex(*index).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Dynamically create the index with the specified number of shards/replicas
	tmpl, err := template.New("T").Parse(mapping)
	if err != nil {
		log.Fatal(err)
	}
	var body bytes.Buffer
	err = tmpl.ExecuteTemplate(&body, "T", struct {
		NumberOfShards   int
		NumberOfReplicas int
	}{
		NumberOfShards:   *shards,
		NumberOfReplicas: *replicas,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.CreateIndex(*index).BodyString(body.String()).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add a tweet
	{
		tweet := Tweet{
			User:     "olivere",
			Message:  "Welcome to Go and Elasticsearch.",
			Retweets: 0,
			Created:  time.Now(),
			Attrs: map[string]interface{}{
				"views": 17,
				"vip":   true,
			},
		}
		_, err := client.Index().
			Index(*index).
			Id("1").
			BodyJson(&tweet).
			Refresh("true").
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read the tweet
	{
		doc, err := client.Get().
			Index(*index).
			Id("1").
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		var tweet Tweet
		if err = json.Unmarshal(doc.Source, &tweet); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s at %s: %s (%d retweets)\n",
			tweet.User,
			tweet.Created,
			tweet.Message,
			tweet.Retweets,
		)
		fmt.Printf("  %v\n", tweet.Attrs)
	}
}

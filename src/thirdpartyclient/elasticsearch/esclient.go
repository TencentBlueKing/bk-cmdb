package elasticsearch

import (
	"configcenter/src/common/blog"
	"context"
	"github.com/olivere/elastic"
)

type EsSrv struct {
	Client *elastic.Client
}

func NewEsClient(esurl string) (*elastic.Client, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient(elastic.SetURL(esurl))
	if err != nil {
		blog.Errorf("create new es client error, err: %v", err)
		return nil, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(esurl).Do(ctx)
	if err != nil {
		// Handle error
		blog.Errorf("esclient connect ping error, err: %v", err)
		return nil, err
	}
	blog.Debug("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return client, nil
}

func (es *EsSrv) Search(query, index string) ([]*elastic.SearchHit, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Search with a term query
	allQuery := elastic.NewQueryStringQuery(query)
	highlight := elastic.NewHighlight()
	highlight.Field("*")
	highlight.RequireFieldMatch(false)
	searchResult, err := es.Client.Search().
		Index(index).                         // search in index like "cmdb"
		Query(allQuery).Highlight(highlight). // specify the query
		Pretty(true).                         // pretty print request and response JSON
		Do(ctx)                               // execute
	if err != nil {
		// Handle error
		blog.Errorf("es search [%s] error, err: %v", query, err)
		return nil, err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	blog.Debug("Query took %d milliseconds\n", searchResult.TookInMillis)
	blog.Debug("Query hits %s\n", searchResult.Hits.Hits)
	return searchResult.Hits.Hits, nil
}

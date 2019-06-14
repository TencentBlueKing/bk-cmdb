package elasticsearch

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"

	"github.com/olivere/elastic"
)

type EsSrv struct {
	Client *elastic.Client
}

func NewEsClient(esurl string, tlsConfig *tls.Config) (*elastic.Client, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	httpClient := &http.Client{}
	client := &elastic.Client{}
	var err error
	if strings.HasPrefix(esurl, "https://") {
		// if use https tls or else, config httpClient first
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient.Transport = tr
		client, err = elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(esurl),
			elastic.SetScheme("https"),
			elastic.SetSniff(false))
		if err != nil {
			blog.Errorf("create new https es client error, err: %v", err)
			return nil, err
		}
	} else {
		client, err = elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(esurl))
		if err != nil {
			blog.Errorf("create new http es client error, err: %v", err)
			return nil, err
		}
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

func (es *EsSrv) CmdbSearch(query elastic.Query, types []string, from, size int) (*elastic.SearchResult, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// search highlight
	highlight := elastic.NewHighlight()
	highlight.Field("*")
	highlight.RequireFieldMatch(false)

	// search for paging
	searchSource := elastic.NewSearchSource()
	searchSource.From(from)
	searchSource.Size(size)

	// search for aggregations value count
	bkObjIdAgg := elastic.NewTermsAggregation().Field(common.BkObjIdAggField)
	typeAgg := elastic.NewTermsAggregation().Field(common.TypeAggField)
	searchResult, err := es.Client.Search().
		// search from es indexes
		Index(common.CMDBINDEX).

		// search from es types of index
		Type(types...).
		SearchSource(searchSource).        // search in index like "cmdb" and paging
		Query(query).Highlight(highlight). // specify the query and highlight
		Pretty(true).                      // pretty print request and response JSON
		// search result with aggregations
		Aggregation(common.BkObjIdAggName, bkObjIdAgg).Aggregation(common.TypeAggName, typeAgg).
		Do(ctx) // execute
	if err != nil {
		// Handle error
		blog.Errorf("es search [%s] error, err: %v", query, err)
		return nil, err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	blog.Debug("Query cmdb took %d milliseconds\n", searchResult.TookInMillis)
	blog.Debug("Query cmdb hits %s\n", searchResult.Hits.Hits)
	return searchResult, nil
}

package elasticsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"github.com/olivere/elastic"
)

type EsSrv struct {
	Client *elastic.Client
}

func NewEsClient(esAddr string, tlsConfig *tls.Config) (*elastic.Client, error) {

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	httpClient := &http.Client{}
	client := &elastic.Client{}
	var err error
	if strings.HasPrefix(esAddr, "https://") {
		// if use https tls or else, config httpClient first
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient.Transport = tr
		client, err = elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(esAddr),
			elastic.SetScheme("https"),
			elastic.SetSniff(false))
		if err != nil {
			blog.Errorf("create new es https es client error, err: %v", err)
			return nil, err
		}
	} else {
		client, err = elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(esAddr))
		if err != nil {
			blog.Errorf("create new http es client error, err: %v", err)
			return nil, err
		}
	}

	// it's amazing that we found new client result success with value nil once a time.
	if client == nil {
		return nil, errors.New("create es client, but it's is nil")
	}
	return client, nil
}

func (es *EsSrv) Search(ctx context.Context, query elastic.Query, types []string, from, size int) (*elastic.SearchResult, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	rid := util.ExtractRequestIDFromContext(ctx)

	// search highlight
	highlight := elastic.NewHighlight()
	highlight.Field("*")
	highlight.RequireFieldMatch(false)
	highlight.Field(common.BKInstIDField)
	highlight.Field(common.BKHostIDField)
	highlight.Field(common.BKAppIDField)

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
		blog.Errorf("es search cond[%v] failed, err: %v, rid: %s", query, err, rid)
		return nil, err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	blog.V(5).Infof("Query cmdb took %d milliseconds\n, rid: %s", searchResult.TookInMillis, rid)
	blog.V(5).Infof("Query cmdb hits %s\n, rid: %s", searchResult.Hits.Hits, rid)
	return searchResult, nil
}

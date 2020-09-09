package elasticsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strings"

	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/ssl"
	"configcenter/src/common/util"

	"github.com/olivere/elastic/v7"
)

type EsSrv struct {
	Client *elastic.Client
}

func NewEsClient(esConf EsConfig) (*elastic.Client, error) {
	// Obtain a client and connect to the default ElasticSearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	httpClient := &http.Client{}
	client := &elastic.Client{}
	var err error
	if strings.HasPrefix(esConf.EsUrl, "https://") {
		tlsConfig := new(tls.Config)
		tlsConfig.InsecureSkipVerify = esConf.TLSClientConfig.InsecureSkipVerify
		if !tlsConfig.InsecureSkipVerify && len(esConf.TLSClientConfig.CAFile) != 0 && len(esConf.TLSClientConfig.CertFile) != 0 && len(esConf.TLSClientConfig.KeyFile) != 0 {
			var err error
			tlsConfig, err = ssl.ClientTLSConfVerity(esConf.TLSClientConfig.CAFile, esConf.TLSClientConfig.CertFile, esConf.TLSClientConfig.KeyFile, esConf.TLSClientConfig.Password)
			if err != nil {
				return nil, err
			}
		}
		// if use https tls or else, config httpClient first
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient.Transport = tr
		client, err = elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(esConf.EsUrl),
			elastic.SetScheme("https"),
			elastic.SetSniff(false),
			elastic.SetBasicAuth(esConf.EsUser, esConf.EsPassword))
		if err != nil {
			blog.Errorf("create new es https es client error, err: %v", err)
			return nil, err
		}
	} else {
		client, err = elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(esConf.EsUrl),
			elastic.SetSniff(false),
			elastic.SetBasicAuth(esConf.EsUser, esConf.EsPassword))
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

func (es *EsSrv) Search(ctx context.Context, query elastic.Query, indexs []string, from, size int) (*elastic.SearchResult, error) {
	// Starting with elastic.v7, you must pass a context to execute each service
	rid := util.ExtractRequestIDFromContext(ctx)

	// search highlight
	highlight := elastic.NewHighlight()
	highlight.Field("*")
	highlight.RequireFieldMatch(false)

	// search for paging
	searchSource := elastic.NewSearchSource()
	searchSource.From(from)
	searchSource.Size(size)

	// search for aggregations value count
	typeAgg := elastic.NewTermsAggregation().Field(common.IndexAggField)

	searchResult, err := es.Client.Search().
		// search from es indexes
		Index(indexs...).
		SearchSource(searchSource).        // search in index like "cmdb" and paging
		Query(query).Highlight(highlight). // specify the query and highlight
		Pretty(true).                      // pretty print request and response JSON
		// search result with aggregations
		Aggregation(common.IndexAggName, typeAgg).
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

type EsConfig struct {
	FullTextSearch  string
	EsUrl           string
	EsUser          string
	EsPassword      string
	TLSClientConfig apiutil.TLSClientConfig
}

// ParseConfigFromKV returns a new config
func ParseConfigFromKV(prefix string, configMap map[string]string) (EsConfig, error) {
	fullTextSearch, _ := cc.String(prefix + ".fullTextSearch")
	url, _ := cc.String(prefix + ".url")
	usr, _ := cc.String(prefix + ".usr")
	pwd, _ := cc.String(prefix + ".pwd")

	conf := EsConfig{
		FullTextSearch: fullTextSearch,
		EsUrl:          url,
		EsUser:         usr,
		EsPassword:     pwd,
	}
	var err error
	conf.TLSClientConfig, err = apiutil.NewTLSClientConfigFromConfig(prefix, nil)
	return conf, err
}

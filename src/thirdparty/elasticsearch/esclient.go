package elasticsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	apiutil "configcenter/src/apimachinery/util"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/ssl"

	"github.com/olivere/elastic/v7"
)

// EsSrv TODO
type EsSrv struct {
	Client *elastic.Client
}

// NewEsClient TODO
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
			tlsConfig, err = ssl.ClientTLSConfVerity(esConf.TLSClientConfig.CAFile, esConf.TLSClientConfig.CertFile,
				esConf.TLSClientConfig.KeyFile, esConf.TLSClientConfig.Password)
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

// Search search elastic with target conditions.
func (es *EsSrv) Search(ctx context.Context, query elastic.Query, indexes []string,
	from, size int) (*elastic.SearchResult, error) {

	// search highlight
	highlight := elastic.NewHighlight()
	// NOTE: 文档高亮同时支持属性和表格属性的keyword高亮
	highlight.Field(metadata.IndexPropertyKeywords)
	highlight.Field(fmt.Sprintf("%s.*.*.%s", metadata.TablePropertyName, metadata.IndexPropertyTypeKeyword))

	highlight.RequireFieldMatch(false)

	searchSource := elastic.NewSearchSource()
	// searchSource.TrackScores(true)
	searchSource.From(from)
	searchSource.Size(size)
	// searchSource.Sort("_score", false)

	searchResult, err := es.Client.Search().
		Index(indexes...).
		SearchSource(searchSource).
		Query(query).Highlight(highlight). // specify the query and highlight
		Pretty(true).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

// Count count data in elastic with target conditions.
func (es *EsSrv) Count(ctx context.Context, query elastic.Query, indexes []string) (int64, error) {

	count, err := es.Client.Count().
		Index(indexes...).
		Query(query).
		Pretty(true).
		Do(ctx)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// EsConfig TODO
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
	conf.TLSClientConfig, err = apiutil.NewTLSClientConfigFromConfig(prefix)
	return conf, err
}

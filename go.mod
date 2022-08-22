module configcenter

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/FZambia/sentinel v1.1.0
	github.com/Shopify/sarama v1.33.0
	github.com/alexmullins/zip v0.0.0-20180717182244-4affb64b04d0
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/apache/thrift v0.12.0
	github.com/aws/aws-sdk-go v1.44.14
	github.com/boj/redistore v0.0.0-20180917114910-cd5dcc76aeff
	github.com/coccyx/timeparser v0.0.0-20161029180942-5644122b3667
	github.com/emicklei/go-restful/v3 v3.7.4
	github.com/fsnotify/fsnotify v1.5.4
	github.com/ghodss/yaml v1.0.0
	github.com/gin-contrib/sessions v0.0.4
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v7 v7.4.1
	github.com/go-zookeeper/zk v1.0.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.8
	github.com/google/uuid v1.3.0
	github.com/gorilla/sessions v1.2.1
	github.com/joyt/godate v0.0.0-20150226210126-7151572574a7 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/juju/ratelimit v1.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/mssola/user_agent v0.5.3
	github.com/olivere/elastic/v7 v7.0.32
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
	github.com/prometheus/client_golang v1.12.2
	github.com/prometheus/client_model v0.2.0
	github.com/rentiansheng/xlsx v1.0.3-r1
	github.com/robfig/cron v1.2.0
	github.com/rs/xid v1.4.0
	github.com/rwynn/monstache v4.12.3+incompatible
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.398
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.398
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc v1.0.398
	github.com/tidwall/gjson v1.14.1
	github.com/xdg-go/scram v1.1.1
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	go.mongodb.org/mongo-driver v1.9.1
	go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful v0.32.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.32.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	golang.org/x/text v0.3.7
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/client-go v0.24.2
	stathat.com/c/consistent v1.0.0
)

replace github.com/rwynn/monstache v4.12.3+incompatible => github.com/ZQHcode/monstache v1.0.0

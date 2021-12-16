蓝鲸CMDB全文索引Monstcache插件
==============================

## 概述

基于特定的版本包进行Monstcache和插件的部署安装;

```shell
monstache/
├── CHANGELOG.md
├── Makefile
├── README.md
├── build
│   └── monstache
│       ├── CHANGELOG.md
│       ├── README.md
│       ├── etc
│       │   ├── config.toml
│       │   └── monstache-plugin.so
│       ├── monstache
│       └── monstache.sh
├── etc
│   ├── config.toml
│   └── extra.toml
├── monstache.sh
└── plugin.go
```

## 配置

**Monstache config.toml配置解释**

| 参数                              | 说明                                                                                                                                                                                                                                                                                               |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| mongo-url                         | MongoDB实例的主节点访问地址。详情请参见。[mongo-url](https://rwynn.github.io/monstache-site/config/#mongo-url)                                                                                                                                                                                     |
| elasticsearch-urls                | Elasticsearch的访问地址。详情请参见 [elasticsearch-urls](https://rwynn.github.io/monstache-site/config/#elasticsearch-urls)                                                                                                                                                                        |
| direct-read-namespaces            | 指定待同步的集合，详情请参见[direct-read-namespaces](https://rwynn.github.io/monstache-site/config/#direct-read-namespaces)。                                                                                                                                                                      |
| direct-read-dynamic-include-regex | 通过正则表达式指定需要监听的集合。此设置可以用来监控符合正则表达式的集合中数据，注意：该功能是在2021-03-18日才合入rel6分支，请使用最新的rel6分支或者2021-03-18之后的release编译最新的Monstache                                                                                                                       |
| change-stream-namespaces          | 如果要使用MongoDB变更流功能，需要指定此参数。启用此参数后，oplog追踪会被设置为无效，详情请参见[change-stream-namespaces](https://rwynn.github.io/monstache-site/config/#change-stream-namespaces)。                                                                                                         |
| namespace-regex                   | 通过正则表达式指定需要监听的集合。此设置可以用来监控符合正则表达式的集合中数据的变化。                                                                                                                                                                                                                       |
| elasticsearch-user                | 访问Elasticsearch的用户名。                                                                                                                                                                                                                                                                        |
| elasticsearch-password            | 访问Elasticsearch的用户密码。                                                                                                                                                                                                                                                                      |
| elasticsearch-max-conns           | 定义连接ES的线程数。默认为4，即使用4个Go线程同时将数据同步到ES。                                                                                                                                                                                                                                         |
| mapper-plugin-path                | 启动插件相对于monstache的路径                                                                                                                                                                                                                                                                        |
| resume                            | 默认为false。设置为true，Monstache会将已成功同步到ES的MongoDB操作的时间戳写入monstache.monstache集合中。当Monstache因为意外停止时，可通过该时间戳恢复同步任务，避免数据丢失。如果指定了cluster-name，该参数将自动开启，详情请参见[resume](https://rwynn.github.io/monstache-site/config/#resume)。                      |


**Monstache plugin extra.toml配置解释**

| 参数                              | 说明                                                                                                                                                                                                                                                                                               |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| elasticsearch-shard-num           | 采用插件场景下必须指定ES的Sharding number。详情请参见。[elasticsearch-shard-num](https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html)                                                                                                                                   |
| elasticsearch-replica-num         | 采用插件场景下必须制定ES的Replica number。详情请参见。[elasticsearch-replica-num](https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html)                                                                                                                                     |

阅读官方文档[monstache doc](https://rwynn.github.io/monstache-site/config/) 可以根据自己的需求进行`高级配置`
## 编译

进入源码根目录执行`make`或`make server`编译指令时，默认会编译后端服务涉及到的所有组件，包括monstache及其对应的monstache-plugin.so插件。您也可以进入到monstache目录，执行`make`命令单独进行monstache及其插件的编译。
## 配置

monstache 涉及到的配置同样需要执行`init.py`执行,主要涉及到elasticsearch-shard-num ，elasticsearch-replica-num两个配置，其余告警配置如: `direct-read-dynamic-include-regex`、`namespace-regex`和`mapper-plugin-path`等如需变更，需要用户手动进行指定。
## 部署安装

整体打包cmdb.tgz时会将插件monstache-plugin.so及对应的配置文件进行打包，monstache二进制需要您按照本文概述中的目录结构示意图进行部署。之后修改 etc/config.toml和etc/extra.toml配置内容，其中配置文件的路径是相对于二进制 `monstache`的路径，如需改动此路径请注意需要同步修改启动脚本`monstache.sh`中的配置文件启动路径。上述步骤完成后您可以通过以下方式运行:

```shell
sh monstache.sh start
```

当然，也可以通过`systemd`或者简单的`nohup`方式运行, 例如 `monstache -f config.toml -mapper-plugin-path monstache-plugin.so`

## 索引管理

插件将会创建附带特定版本后缀的真实ES索引，如`bk_cmdb.biz_20210701`, 并且只会在索引不存在时创建，特定版本索引的结构信息在插件代码中固定，在索引结构发生变化时插件中版本后缀也会发生变化。
在成功创建索引后，插件会为每一个索引创建系统别名，如`bk_cmdb.biz` `bk_cmdb.set` `bk_cmdb.module` `bk_cmdb.host` `bk_cmdb.model` `bk_cmdb.object_instance`, 这些别名为蓝鲸CMDB内部索引、查询等操作所使用的别名。

索引分词器的指定是在plugin.go中完成,日常维护，如做reindex文档迁移后，需要将最终的真实索引和蓝鲸CMDB系统别名关联, 以保证系统能够正确处理文档数据。

阅读官方文档 [elastic reindex doc](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html) 了解reindex操作。
阅读官方文档 [elastic alias doc](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-aliases.html) 了解索引别名机制。

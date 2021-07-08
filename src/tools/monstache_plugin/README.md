蓝鲸CMDB全文索引Monstcache插件
==============================

## 概述

基于特定的版本包进行Monstcache和插件的部署安装;

```shell
.
└── monstache_plugin
    ├── etc
    │   ├── monstache-plugin.so
    │   └── config.toml
    ├── monstache
    ├── monstache.sh
    └── CHANGELOG.md
    └── README.md
```

## 部署安装

将版本包内容放到指定的安装目录，修改etc/config.toml配置内容后，通过以下方式运行:

```shell
sh monstache.sh start
```

当然，也可以通过`systemd`或者简单的`nohup`方式运行, 例如 `monstache -f config.toml -mapper-plugin-path monstache-plugin.so`

## 配置

阅读官方文档[monstache doc](https://rwynn.github.io/monstache-site/config/) 可以根据自己的需求进行`高级配置`

## 索引管理

插件将会创建附带特定版本后缀的真实ES索引，如`bk_cmdb.biz_20210701`, 并且只会在索引不存在时创建，特定版本索引的结构信息在插件代码中固定，在索引结构发生变化时插件中版本后缀也会发生变化。
在成功创建索引后，插件会为每一个索引创建系统别名，如`bk_cmdb.biz` `bk_cmdb.set` `bk_cmdb.module` `bk_cmdb.host` `bk_cmdb.model` `bk_cmdb.object_instance`, 这些别名为蓝鲸CMDB内部索引、查询等操作所使用的别名。

日常维护，如做reindex文档迁移后，需要将最终的真实索引和蓝鲸CMDB系统别名关联, 以保证系统能够正确处理文档数据。

阅读官方文档 [elastic reindex doc](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html) 了解reindex操作。
阅读官方文档 [elastic alias doc](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-aliases.html) 了解索引别名机制。

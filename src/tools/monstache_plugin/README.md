BK-CMDB elastic monstache plugin
================================

## Overview

Install base on the release package:

```shell
.
└── bk-cmdb-monstache-plugin-517da48-21.06.29
    ├── etc
    │   ├── bk-cmdb-monstache-plugin.so
    │   └── config.toml
    ├── monstache
    ├── monstache.sh
    └── README.md
```

## Install

Put the release package into your own INSTALL-PATH and run monstache with plugin like:

```shell
sh monstache.sh start
```

Also, you could run it in `systemd` or `nohup` mode base on `monstache -f config.toml -mapper-plugin-path bk-cmdb-monstache-plugin.so`

## Configuration

You could just use the `config.toml` in release package directly, it already includes the correct base configurations.
And read [monstache doc](https://rwynn.github.io/monstache-site/config/) for `Advanced Configuration`

## Indexes

The elastic plugin would create bk-cmdb target indexes of the default names `bk_cmdb.biz` `bk_cmdb.set` `bk_cmdb.module` `bk_cmdb.host` `bk_cmdb.model` `bk_cmdb.object_instance` with the version postfix.
You could create custom names and reindex index documents, make it alias to bk-cmdb default name.

Read [elastic reindex doc](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html) for elastic reindex operation details.
And read [elastic alias doc](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-aliases.html) for elastic alias details.

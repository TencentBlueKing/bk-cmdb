多租户版本数据升级工具
==================

多租户版本数据升级工具支持将数据从非多租户版本(3.14)的最新版本升级到不开启多租户功能的多租户版本(3.15)的第一个版本，支持原地升级和数据迁移两种升级方式

## 升级流程

### 前置准备
将cmdb升级到3.14的最新版本，如果没有升级到3.14的最新版本的话升级会失败或者产生脏数据

### 升级方式一：原地升级

**注意：因为非多租户版本的数据和多租户版本的数据不兼容，所以原地升级期间需要停服发布，否则可能升级失败或者产生脏数据**

升级流程；
1. 停止cmdb服务
2. 备份cmdb数据库
3. 升级cmdb的mongodb到7.0版本
4. 使用升级工具原地升级全部cmdb数据，升级命令：
    - 使用方式
      ```
      ./tool_ctl migrate-tenant in-place-upgrade [flags]
      ```
    - 命令行参数
      ```
      --mongo-uri string               the mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI
      --mongo-rs-name string           mongodb replica set name (default "rs0")
      --watch-mongo-uri string         watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --watch-mongo-rs-name string     watch db replica set name (default "rs0")
      --skip-remove-supplier-account   skip removing all bk_supplier_account fields, upgrade will be faster, but data will have redundant fields
      ```
    - 示例
      ```
      ./tool_ctl migrate-tenant in-place-upgrade --mongo-uri="mongodb://127.0.0.1:27017/cmdb" --mongo-rs-name=rs0 --watch-mongo-uri="mongodb://127.0.0.1:27017/cmdb_events" --watch-mongo-rs-name=rs0
      ```
    - **注意：**
      1. 原地升级操作默认会清理所有数据的bk_supplier_account字段，数据量大的情况下耗时非常长，如果可以接受数据上保留冗余的bk_supplier_account字段的话可以指定`--skip-remove-supplier-account`选项跳过这个清理操作，加速升级过程
      2. mongo-uri指定的mongodb的连接用户需要具备dbOwner权限，否则更新部分数据表的配置时会操作失败
5. 启动新版本的cmdb服务，等待所有migrate相关的job执行完毕
6. 执行以下命令刷新db索引
    ``` shell
    curl -X POST -H 'Content-Type:application/json' -H 'X-Bk-Tenant-Id:default' -H 'X-Bkcmdb-User:migrate'  http://{adminserver的访问地址}/migrate/v3/migrate/sync/db/index
    ``` 

### 升级方式二：数据迁移
1. 准备一个新的7.0版本的mongodb
2. 备份旧版本的cmdb数据库
3. 使用升级工具将全部cmdb数据迁移到新db，迁移命令：
    - 使用方式
      ```
      ./tool_ctl migrate-tenant copy-to-new-db --full-sync [flags]
      ```
    - 命令行参数
      ```
      --old-mongo-uri string               old mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --old-mongo-rs-name string           old mongodb replica set name (default "rs0")
      --mongo-uri string                   new mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI
      --mongo-rs-name string               new mongodb replica set name (default "rs0")
      --old-watch-mongo-uri string         old watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --old-watch-mongo-rs-name string     old watch db replica set name (default "rs0")
      --watch-mongo-uri string             new watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --watch-mongo-rs-name string         new watch db replica set name (default "rs0")
      ```
    - 示例
      ```
      ./tool_ctl migrate-tenant copy-to-new-db --full-sync --old-mongo-uri="mongodb://127.0.0.1:27017/old_cmdb" --old-mongo-rs-name=rs0 --mongo-uri="mongodb://127.0.0.1:27017/cmdb" --mongo-rs-name=rs0 --old-watch-mongo-uri="mongodb://127.0.0.1:27017/old_watch_db" --old-watch-mongo-rs-name=rs0 --watch-mongo-uri="mongodb://127.0.0.1:27017/cmdb_events" --watch-mongo-rs-name=rs0
      ```
4. 使用升级工具将上一次全量迁移的开始时间之后变更的旧db数据迁移到新db，迁移命令：
    - 使用方式
      ```
      ./tool_ctl migrate-tenant copy-to-new-db [flags]
      ```
    - 命令行参数
      ```
      --old-mongo-uri string               old mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --old-mongo-rs-name string           old mongodb replica set name (default "rs0")
      --mongo-uri string                   new mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI
      --mongo-rs-name string               new mongodb replica set name (default "rs0")
      --old-watch-mongo-uri string         old watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --old-watch-mongo-rs-name string     old watch db replica set name (default "rs0")
      --watch-mongo-uri string             new watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb
      --watch-mongo-rs-name string         new watch db replica set name (default "rs0")
      --start-from uint32                  unix timestamp to start incremental sync from, would use full sync start time if not set
      ```
    - 示例
      ```
      ./tool_ctl migrate-tenant copy-to-new-db --old-mongo-uri="mongodb://127.0.0.1:27017/old_cmdb" --old-mongo-rs-name=rs0 --mongo-uri="mongodb://127.0.0.1:27017/cmdb" --mongo-rs-name=rs0 --old-watch-mongo-uri="mongodb://127.0.0.1:27017/old_watch_db" --old-watch-mongo-rs-name=rs0 --watch-mongo-uri="mongodb://127.0.0.1:27017/cmdb_events" --watch-mongo-rs-name=rs0
      ```
    - **注意：升级工具会持续运行，将增量数据从旧db迁移到新db，不会主动退出**
5. 等待一段时间，观察升级工具输出的事件信息，在增量数据都迁移到新环境后停止旧环境的写操作，无事件后关闭增量数据迁移工具，准备切换到新环境
6. 在新db上部署一套新版本的cmdb服务，等待所有migrate相关的job执行完毕后，将流量切到新环境
7. 停止旧环境的服务，备份旧版本的cmdb数据库，回收旧版本db

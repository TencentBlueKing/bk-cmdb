upgrader
=========

## 概述

upgrader包用于全文检索相关的数据初始化和升级，包括Elasticsearch索引的初始化和升级，数据升级等功能。

## 目录结构

```
.
├── current.go      用于生成并创建最新版本的Elasticsearch索引
├── index.go        包含Elasticsearch索引相关的工具函数
├── readme.md       帮助文档
├── upgrader.go     通过运行upgrader进行索引和数据的升级
└── v{version}.go   指定版本的upgrader逻辑，例如v1.go文件存放版本1的upgrader
```

## 升级方式

### 前置准备

- 将最新版本的Elasticsearch索引信息写入到`current.go`文件中，如果索引有更新，则需要更新索引版本号
- 将当前版本的索引信息和与上一个版本对比需要进行的升级逻辑写入到新版本的`v{version}.go`文件中，并通过`RegisterUpgrader`方法注册到upgrader池中
- 如果涉及到索引的删除操作，则直接从新增该索引的版本开始删除掉索引的相关操作逻辑

### upgrader执行流程

- 获取当前版本信息，从当前版本开始执行upgrader进行升级，如果已经是最新版本则不需要进行升级
- 因为全文检索数据同步需要依赖于对应版本的Elasticsearch索引，所以优先创建最新版本的Elasticsearch索引，此时可以开始进行数据同步
- 按版本顺序执行每一个upgrader，其中如果涉及到数据迁移则直接将数据迁移到最新版本的Elasticsearch索引中

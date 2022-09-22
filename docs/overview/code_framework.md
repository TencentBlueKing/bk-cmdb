# 蓝鲸智云配置平台的代码结构

## 整体结构

```
├── api
├── cmd
│   ├── apiserver
│   ├── scene_server
│   │   ├── admin_server
│   │   ├── auth_server
│   │   ├── cloud_server
│   │   ├── datacollection
│   │   ├── event_server
│   │   ├── host_server
│   │   ├── operation_server
│   │   ├── proc_server
│   │   ├── synchronize_server
│   │   ├── task_server
│   │   └── topo_server
│   ├── source_controller
│   │   ├── cacheservice
│   │   └── coreservice
│   └── web_server
├── configs
├── docs
├── framework
├── gse                                             
├── pkg
├── resources
├── scripts
├── test
├── thirdparty
├── tools
└── ui
```

## api

api调用相关目录

## cmd

CMDB微服务目录，以下为划分的目录：

### web-server

web-server是基于gin框架构建的web服务

### api_server

api-server是基于go-restful框架构建的API服务

### scene_server

scene_server是基于go-restful框架构建的场景层服务，以下为划分的微服务目录：
- admin_server
- auth_server
- cloud_server
- datacollection
- event_server
- host_server
- operation_server
- proc_server
- synchronize_server
- task_server
- topo_server

### source_controller

source_controlle是基于go-restful框架构建的资源层服务，以下为划分的微服务目录：
- cacheservice
- coreservice

## configs

CMDB配置文件模板

## docs

CMDB文档

## framework

CMDB3.0二次开发框架

## gse

gse文件

## pkg

用于pkg外部引用的包，pkg内的包不可以引用pkg外的包

## resources

CMDB依赖的资源文件，包含i18n目录存放国际化相关文件

## scripts

CMDB脚本文件

## test

CMDB测试文件

## thirdparty

CMDB和第三方接入相关的文件

## tools

客户端管理工具和辅助脚本工具代码

## ui

前端代码

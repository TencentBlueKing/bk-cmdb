# BK-CMDB

蓝鲸配置平台（蓝鲸CMDB）是一个面向资产及应用的企业级配置管理平台。

本文档内容为如何在 Kubernetes 集群上部署 BK-CMDB web服务。

说明：内置的mongodb、redis等组件仅用于测试环境，正式环境部署必须配置为外部组件。

## BK-CMDB 部署项目

### 架构设计

* [点击这里](https://github.com/TencentBlueKing/bk-cmdb/blob/master/docs/overview/architecture.md)

## 部署

### 环境要求

- Kubernetes 1.12+
- Helm 3+

快速部署k8s集群可参考：[蓝鲸官方白皮书](https://bk.tencent.com/docs/markdown/ZH/DeploymentGuides/7.1/get-k8s-create-bcssh.md)

### 镜像制作

制作cmdb镜像请参考：[CMDB 编译指南](../../../overview/source_compile.md)

### 安装Chart

使用以下命令安装名称为`bkcmdb`的release, 其中`<bkcmdb helm repo url>`代表helm仓库地址, password为自己设置的任意密码:

```shell
# 添加helm仓库
$ helm repo add bitnami https://charts.bitnami.com/bitnami
$ helm repo add bkee <bkcmdb helm repo url>
# 更新并拉取依赖
$ helm dependency update
# 执行部署
$ helm install cmdb-web bkee/cmdb-web --set mongodb.auth.password=${password} --set redis.auth.password=${password}
```

上述命令将使用默认配置在Kubernetes集群中部署BK-CMDB web服务, 并输出访问指引。
注：执行部署前请检查values.yaml镜像配置中的镜像地址

### 卸载Chart

使用以下命令卸载`cmdb-web`:

```shell
$ helm uninstall cmdb-web
```

上述命令将移除所有和cmdb-web相关的Kubernetes组件。

## Chart依赖

- [bitnami/mongodb](https://github.com/bitnami/charts/tree/master/bitnami/mongodb)
- [bitnami/redis](https://github.com/bitnami/charts/tree/master/bitnami/redis)
- [bitnami/zookeeper](https://github.com/bitnami/charts/tree/master/bitnami/zookeeper)
- [bitnami/elasticsearch](https://github.com/bitnami/charts/tree/master/bitnami/elasticsearch)

## 配置说明

各项配置集中在仓库的一个values.yaml文件之中

### 镜像配置

|      参数       |     描述     |    默认值    |
| :-------------: | :----------: | :----------: |
| image.registry | 镜像源域名 | mirrors.tencent.com |
| image.pullPolicy | 镜像拉取策略 | IfNotPresent |

### 蓝鲸产品URL配置

|   参数   |   描述   |         默认值          |
| :------: | :------: | :---------------------: |
| bkPaasUrl | paas地址 | http://paas.example.com |
| bkIamApiUrl | bkiam后端地址 | http://bkiam-web |
| bkComponentApiUrl | 蓝鲸ESB地址 | http://bkapi.paas.example.com |
| bkLoginApiUrl | 蓝鲸登录地址 | http://bk-login-web |
| bkNodemanUrl | 节点管理地址 | http://apps.paas.example.com/bk--nodeman |

### webserver服务配置说明

|              参数               |              描述               |              默认值              |
|:-----------------------------:| :-----------------------------: | :------------------------------: |
|       webserver.enabled       | 是否在执行helm时启动 |               true               |
|  webserver.image.repository   |        服务镜像名        | cmdb_webserver |
|      webserver.image.tag      |          服务镜像版本           | {TAG_NAME} |
|      webserver.replicas       |           pod副本数量           |                1                 |
|        webserver.port         |            服务端口             |                80                |
|   webserver.ingress.enabled   | 开启ingress访问 | true |
|    webserver.ingress.hosts    | ingress代理访问的域名 |cmdb.example.com|
|    webserver.service.type     | 服务类型 | ClusterIP |
| webserver.service.targetPort  | 代理的目标端口 | 80 |
|  webserver.service.nodePort   | 访问端口 |  |
|   webserver.command.logDir    |          日志存放路径           |              /data/cmdb/cmdb_webserver/logs              |
|  webserver.command.logLevel   |            日志等级             |                3                 |
| webserver.command.logToStdErr |     是否把日志输出到stderr      |              false               |
|       webserver.workDir       |            工作目录             |      /data/cmdb/cmdb_webserver      |

### mongodb配置
|                 参数                 |              描述               |              默认值              |
| :----------------------------------: | :-----------------------------: | :------------------------------: |
|      mongodb.enabled      | 是否部署mognodb，如果需要使用外部数据库，设置为`false`并配置`mongodb.externalMongodb`和`mongodb.watch`下关于外部mongodb的配置 |               true               |

`mongodb.externalMongodb` 和 `mongodb.watch` 开头的配置，可根据原`mongodb.yaml`中的配置进行修改

### redis配置
|                 参数                 |              描述               |              默认值              |
| :----------------------------------: | :-----------------------------: | :------------------------------: |
|      redis.enabled      | 是否部署redis，如果需要使用外部数据库，设置为`false`并配置`redis.redis`、`redis.snapshotRedis`、`redis.discoverRedis`、`redis.netCollectRedis`下关于外部redis的配置 |               true               |

`redis.redis`、`redis.snapshotRedis`、`redis.discoverRedis`、`redis.netCollectRedis` 开头的配置，可根据原`redis.yaml`中的配置进行修改

### zookeeper配置
|                 参数                 |              描述               |              默认值              |
| :----------------------------------: | :-----------------------------: | :------------------------------: |
|      zookeeper.enabled      | 是否部署zookeeper作为配置发现中心、服务发现中心，如果需要使用外部zookeeper组件，设置为`false`并配置`configAndServiceCenter.addr` |               true               |

### 配置发现中心、服务发现中心配置

|            参数             |                             描述                             | 默认值 |
| :-------------------------: | :----------------------------------------------------------: | :----: |
| configAndServiceCenter.addr | 外部配置发现中心、服务发现中心地址，当zookeeper.enabled配置为`false`时，使用此参数连接外部组件 |        |

### elasticsearch配置

|                 参数                  |                             描述                             | 默认值 |
|:-----------------------------------:| :----------------------------------------------------------: | :----: |
|        web.es.fullTextSearch        | 开启全文索引开关，可选值为`on` 和 `off`, 默认关闭 | off       |

### monstache配置
monstache是一个用于将mongodb的数据同步到es去创建索引的一个组件

|            参数             |                             描述                             | 默认值 |
| :-------------------------: | :----------------------------------------------------------: | :----: |
| monstache.enabled | 是否启动内部部署的monstache，如果需要使用外部monstache组件，设置为`false` | false       |
| monstache.image.repository | 服务镜像名 |   cmdb_monstache     |
| monstache.image.tag | 服务镜像版本 |   {TAG_NAME}     |
|         monstache.replicas         |           pod副本数量           |                1                 |
|           monstache.port           |            服务端口             |                80                |
|        monstache.workDir         |       工作路径        | /data/cmdb/monstache |
|        monstache.configDir         |       需要的配置文件路径        | /data/cmdb/monstache/etc |
|        monstache.directReadDynamicIncludeRegex         | monstache配置内容 |内容过长请查看原value.yaml文件|
|        monstache.mapperPluginPath         | monstache配置内容 |/data/cmdb/monstache/monstache-plugin.so|
|        monstache.elasticsearchShardNum         | monstache配置内容 | 1 |
|        monstache.elasticsearchReplicaNum         | monstache配置内容 | 1 |

### bkLogConfig配置
- bkLogConfig配置用于配置接入蓝鲸日志平台功能

|            参数             |                             描述                             | 默认值 |
| :-------------------------: | :----------------------------------------------------------: | :----: |
| bkLogConfig.file.enabled | 是否采集容器内落地文件日志 | false       |
| bkLogConfig.file.dataId | 采集容器内落地文件日志的dataid，dataid在日志平台上申请分配 | 1       |
| bkLogConfig.std.enabled | 是否采集容器标准输出日志 | false       |
| bkLogConfig.std.dataId | 采集容器标准输出日志的dataid，dataid在日志平台上申请分配 | 1       |

### serviceMonitor配置
- serviceMonitor配置用于配置服务监控功能

|            参数             |                             描述                             | 默认值 |
| :-------------------------: | :----------------------------------------------------------: | :----: |
| serviceMonitor.enabled | 是否开启服务监控，采集cmdb业务指标数据 | false       |
| serviceMonitor.interval | cmdb业务指标数据采集间隔时间 | 15s |

## 配置案例

### 1. 使用外接mongodb
```yaml
mongodb:
  enabled: false
  ...
  # external mongo configuration
  externalMongodb:
    enabled: xxx
    usr: xxx
    pwd: xxx
    database: xxx
    host: 127.0.0.1:27017
    maxOpenConns: xxx
    maxIdleConns: xxx
    mechanism: xxx
    rsName: xxx
    socketTimeoutSeconds: xxx
  watch:
    usr: xxx
    pwd: xxx
    database: xxx
    host: 127.0.0.1:27017
    maxOpenConns: xxx
    maxIdleConns: xxx
    mechanism: xxx
    rsName: xxx
    socketTimeoutSeconds: xxx
```

### 2. 使用外接redis

```yaml
redis:
  enabled: false
  ...
  # external redis configuration
  redis:
    host: 127.0.0.1:6379
    pwd: xxx
    database: xxx
    maxOpenConns: xxx
    maxIdleConns: xxx
    sentinelPwd: xxx
    masterName: xxx

  snapshotRedis:
    host: 127.0.0.1:6379
    pwd: xxx
    database: xxx
    maxOpenConns: xxx
    maxIdleConns: xxx
    sentinelPwd: xxx
    masterName: xxx

  discoverRedis:
    host: 127.0.0.1:6379
    pwd: xxx
    database: xxx
    maxOpenConns: xxx
    maxIdleConns: xxx
    sentinelPwd: xxx
    masterName: xxx

  netCollectRedis:
    host: 127.0.0.1:6379
    pwd: xxx
    database: xxx
    maxOpenConns: xxx
    maxIdleConns: xxx
    sentinelPwd: xxx
    masterName: xxx
```

### 3. 使用外接zookeeper作为配置发现中心和服务发现中心

```yaml
zookeeper:
  enabled: false
configAndServiceCenter:
  addr: 127.0.0.1:2181

```

### 4. 使用elasticsearch相关操作

- 如果在cmdb中使用es，首先需要开始es的开关

```yaml
web:
  es:
    fullTextSearch: "on"
```

- 在cmdb中使用elasticsearch需要依赖两个组件，一个是elasticsearch本身，一个monstache（用于将mongodb数据同步到elasticsearch）

  （1）使用内置组件

  ​	helm chart中有内置的elasticsearch和monstache，可通过下面操作打开：

  ```yaml
  elasticsearch:
    enabled: true
  
  ··
  
  monstache:
    enabled: true
  ```

  ​	将elasticsearch和monstache的enabled变为true即可

  

  （2）使用外接组件

  这里仅需配置连接外置的elasticsearch，这时外置的monstach已经与cmdb没有配置上的联系

  ```yaml
  common:
    es: 
      url: xxx
      usr: xxx
      pwd: xxx
  ```

  配置上外部es的url，账户密码的信息即可

- 当然也可以使用内置的monstache，连接外部的elasticsearch

  ```yaml
  monstache:
    enabled: true
  	
  ···
  
  common:
    es: 
      url: xxx
      usr: xxx
      pwd: xxx
  ```

  将monstache.enabled变为true，配置好外部的elasticsearch的配置即可

### 5. 配置webserver不同的服务暴露方式

默认通过Ingress暴露服务，也可以使用以下方式：

- 使用NodePort直接访问

  ```yaml
  webserver:
    ingress:
      enabled: false
      ...
    service:
      type: "NodePort"   
      ports:
        ...
          nodePort: 32033 # 端口可以自定义
  ```

  ```yaml
  common:
    ...
    site:
      domainUrl: http://127.0.0.1:32033/ # ip需要根据实际情况的进行配置，端口为上面配置的同一端口
  ```

  修改上述配置后，即可通过`ip:32033`的方式访问

 ### 6. 开启权限验证
 通过进行下面的配置：
```yaml
开启权限
iam:
  auth:
    enabled: true

// 配置权限中心和esb地址、app code、app secret，开启前端的auth
bkIamApiUrl: xxx
bkComponentApiUrl: xxx

web:
 ...
  auth:
    appCode: xxx
    appSecret: xxx
  esb:
    appCode: xxx
    appSecret: xxx
 ...
  webServer:
    site:
      authScheme: iam
```

### 7. blueking方式登陆
```yaml
通过将登陆方式设置为蓝鲸登陆方式和配置蓝鲸登陆地址等信息：

# pass地址
bkPaasUrl: xxx
# bk-login地址
bkLoginApiUrl: xxx

web:
  ...
  webServer:
    site:
      appCode: bk_cmdb
    ...
    login:
      version: blueking
```

## 常见问题

### 1. cmdb的helm chart启动后如何访问

答：因为默认的访问方式是通过ingress访问，域名为 cmdb.example.com，所以需要配置 cmdb.example.com 的dns解析，例如在机器的/usr/hosts文件中配置：

```yaml
127.0.0.1 cmdb.example.com 
```

在minikube环境通过下面指令启用` Ingress `控制器 
```yaml
  minikube addons enable ingress
```
配置完后，通过访问`cmdb.example.com/login`地址进行登陆，默认 的账号为`cc`，密码为`cc`

### 2. 想要配置多个外置zookeeper地址作为服务中心怎么办？

答：通过,(逗号)分隔，如下：
```
configAndServiceCenter:
  addr: 127.0.0.1:2181,127.0.0.2:2181

```

### 3. 想要配置多个外置redis地址怎么办？

答：通过,(逗号)分隔，如下：
```
redis:
  ...
  # external redis configuration
  redis:
    host: 127.0.0.1:6379,128.0.0.1:6379

```

### 4. 想要配置多个外置mongo地址怎么办？

答：通过,(逗号)分隔，如下：
```
mongodb:
  # external mongo configuration
  externalMongodb:
    enabled: xxx
    usr: xxx
    pwd: xxx
    database: xxx
    host: 127.0.0.1:27017,127.0.0.1:27018
```





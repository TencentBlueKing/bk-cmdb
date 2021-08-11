# CMDB 部署文档

## 依赖第三方组件

* ZooKeeper >= 3.4.11
* Redis   >= 3.2.11
* MongoDB >= 4.2
* Elasticsearch >= 7.0.0 (用于全文检索功能)
* Monstache >= 6.0.0 (用于全文检索功能)

## CMDB 微服务进程清单

### 1. web层服务进程
* cmdb_webserver

### 2. 服务网关进程
* cmdb_apiserver


### 3. 场景层服务进程
* cmdb_adminserver
* cmdb_eventserver
* cmdb_hostserver
* cmdb_procserver
* cmdb_toposerver
* cmdb_datacollection
* cmdb_operationserver
* cmdb_synchronizeserver
* cmdb_taskserver

### 4. 资源管理服务进程
* cmdb_coreservice

---

## 部署介绍
### 1. 部署ZooKeeper

请参看官方资料 [ZooKeeper](https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#ch_deployment)

推荐版本下载： [ZooKeeper 3.4.12](https://mirrors.tuna.tsinghua.edu.cn/apache/zookeeper/zookeeper-3.4.12/zookeeper-3.4.12.tar.gz)

### 2. 部署Redis

请参看官方资料 [Redis](https://redis.io/download)

推荐版本下载： [Redis 3.2.11](http://download.redis.io/releases/redis-3.2.11.tar.gz)


### 3. 部署MongoDB

请参考官方资料 [MongoDB](https://docs.mongodb.com/manual/installation/)

推荐版本下载：[MongoDB 4.2.8](https://www.mongodb.com/dr/fastdl.mongodb.org/linux/mongodb-linux-x86_64-rhel70-4.2.8.tgz/download)

### 4. Release包下载

官方发布的 **Linux Release** 包下载地址见[这里](https://github.com/Tencent/bk-cmdb/releases)。如果你想自已编译，具体的编译方法见[这里](source_compile.md)。

### 5. 配置数据库

#### 1. Redis需要打开auth认证的功能，并为其配置密码

##### a. 修改配置文件
redis的配置文件默认在/etc/redis.conf，找到如下行：
``` json
#requirepass foobared
``` 
去掉前面的注释，并修改为所需要的密码：
``` json
 requirepass myPassword （其中myPassword就是要设置的密码）
``` 
##### b. 重启Redis
如果Redis已经配置为service服务，可以通过以下方式重启：
```json
service redis restart
```
##### c. 登录验证
设置Redis认证密码后，客户端登录时需要使用-a参数输入认证密码,举例如下：
```json
$ ./redis-cli -h 127.0.0.1 -p 6379 -a myPassword
127.0.0.1:6379> config get requirepass
1) "requirepass"
2) "myPassword"
```
看到类似上面的输出，说明Reids密码认证配置成功。

#### 2. 安装MongoDB后，配置集群，创建数据库 cmdb

#### 3. 为新创建的数据库设置用户名和密码

> MongoDB 示例:

mongodb以集群的方式启动，需加入参数--replSet,如--replSet=rs0

进入mongodb后，在members中配置集群ip和端口

如: 配置集群中只有单台机器
```json
 >rs.initiate({ _id : "rs0",members: [{ _id: 0, host: "ip:port" }]})
```

注:rs0为集群名字，仅作展示，用户使用中可以根据实际情况自行配置

接下来登陆MongoDB后，根据需求执行以下命令:

- 未开启ES情况(用于全文检索, 可选, 控制开关见第9步的full_text_search)
``` json
 > use cmdb
 > db.createUser({user: "cc",pwd: "cc",roles: [ { role: "readWrite", db: "cmdb" } ]})
```
- 开启ES情况(用于全文检索, 可选, 控制开关见第9步的full_text_search)
``` json
 > use cmdb
 > db.createUser({user: "cc",pwd: "cc",roles: [ { role: "readWrite", db: "cmdb" },{ role: "readWrite", db: "monstache" } ]})
```

**注：以上用户名、密码、数据库名仅作示例展示，用户使用中可以根据实际情况自行配置。如果安装的MongoDB的版本大于等于3.6，需要手动修改init.py自动生成的配置文件，详细步骤参看init.py相关小节。**

详细手册请参考官方资料 [MongoDB](https://docs.mongodb.com/manual/reference/method/db.createUser/)

### 6. 部署Elasticsearch (用于全文检索, 可选, 控制开关见第9步的full_text_search)

官方下载 [ElasticSearch](https://www.elastic.co/cn/downloads/past-releases)
搜索7.x的版本下载，推荐下载7.0.0
下载后解压即可，解压后找到配置文件config/elasticsearch.yml，可以配置指定network.host为
具体的host的地址
然后到目录的bin目录下运行(注意，不能使用root权限运行，**要普通用户**)：
```shell
./elasticsearch
```

如果想部署高可用可扩展的ES，可参考官方文档[ES-Guide](https://www.elastic.co/guide/index.html)

### 7.  部署Monstache (用于全文检索, 可选, 控制开关见第9步的full_text_search)

 [下载 Monstache 可执行文件](https://github.com/rentiansheng/monstache/tree/rel6_release/build)

**Monstache-Mongodb-Es 版本关系:**

| Monstache version | Git branch (used to build plugin) | Docker tag                         | Description             | Elasticsearch    | MongoDB   |
| ----------------- | --------------------------------- | ---------------------------------- | ----------------------- | ---------------- | --------- |
| 3                 | rel3                              | rel3                               | mgo community go driver | Versions 2 and 5 | Version 3 |
| 4                 | master                            | rel4 (note this used to be latest) | mgo community go driver | Version 6        | Version 3 |
| 5                 | rel5                              | rel5                               | MongoDB, Inc. go driver | Version 6        | Version 4 |
| 6                 | rel6                              | rel6, latest                       | MongoDB, Inc. go driver | Version 7        | Version 4 |

**Monstache配置解释**

| 参数                     | 说明                                                                                                                                                                                                                                                                                               |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| mongo-url                | MongoDB实例的主节点访问地址。详情请参见。[mongo-url](https://rwynn.github.io/monstache-site/config/#mongo-url)                                                                                                                                                                                     |
| elasticsearch-urls       | Elasticsearch的访问地址。详情请参见 [elasticsearch-urls](https://rwynn.github.io/monstache-site/config/#elasticsearch-urls)                                                                                                                                                                        |
| direct-read-namespaces   | 指定待同步的集合，详情请参见[direct-read-namespaces](https://rwynn.github.io/monstache-site/config/#direct-read-namespaces)。                                                                                                                                                                      |
| direct-read-dynamic-include-regex   | 通过正则表达式指定需要监听的集合。此设置可以用来监控符合正则表达式的集合中数据，注意：该功能是在2021-03-18日才合入rel6分支，请使用最新的rel6分支或者2021-03-18之后的release编译最新的Monstache                        |
| change-stream-namespaces | 如果要使用MongoDB变更流功能，需要指定此参数。启用此参数后，oplog追踪会被设置为无效，详情请参见[change-stream-namespaces](https://rwynn.github.io/monstache-site/config/#change-stream-namespaces)。                                                                                                |
| namespace-regex          | 通过正则表达式指定需要监听的集合。此设置可以用来监控符合正则表达式的集合中数据的变化。                                                                                                                                                                                                             |
| elasticsearch-user       | 访问Elasticsearch的用户名。                                                                                                                                                                                                                                                                        |
| elasticsearch-password   | 访问Elasticsearch的用户密码。                                                                                                                                                                                                                                                                      |
| elasticsearch-max-conns  | 定义连接ES的线程数。默认为4，即使用4个Go线程同时将数据同步到ES。                                                                                                                                                                                                                                   |
| dropped-collections      | 默认为true，表示当删除MongoDB集合时，会同时删除ES中对应的索引。                                                                                                                                                                                                                                    |
| dropped-databases        | 默认为true，表示当删除MongoDB数据库时，会同时删除ES中对应的索引。                                                                                                                                                                                                                                  |
| resume                   | 默认为false。设置为true，Monstache会将已成功同步到ES的MongoDB操作的时间戳写入monstache.monstache集合中。当Monstache因为意外停止时，可通过该时间戳恢复同步任务，避免数据丢失。如果指定了cluster-name，该参数将自动开启，详情请参见[resume](https://rwynn.github.io/monstache-site/config/#resume)。 |
| resume-strategy          | 指定恢复策略。仅当resume为true时生效，详情请参见[resume-strategy](https://rwynn.github.io/monstache-site/config/#resume-strategy)。                                                                                                                                                                |
| verbose                  | 默认为false，表示不启用调试日志。                                                                                                                                                                                                                                                                  |
| cluster-name             | 指定集群名称。指定后，Monstache将进入高可用模式，集群名称相同的进程将进行协调，详情请参见[cluster-name](https://rwynn.github.io/monstache-site/config/#cluster-name)。                                                                                                                             |
| mapping                  | 指定ES索引映射。默认情况下，数据从MongoDB同步到ES时，索引会自动映射为`数据库名.集合名`。如果需要修改索引名称，可通过该参数设置，详情请参见[Index Mapping](https://rwynn.github.io/monstache-site/advanced/#index-mapping)。
|
 

**config.toml 内容举例如下：**

```shell
# cmdb connection settings

# connect to MongoDB using the following URL
mongo-url =  "mongodb://localhost:27017/cmdb"
# connect to the Elasticsearch REST API at the following node URLs
elasticsearch-urls = ["http://localhost:9200"]

# frequently required settings

# if you need to seed an index from a collection and not just listen and sync changes events
# you can copy entire collections or views from MongoDB to Elasticsearch
direct-read-namespaces = [""]
direct-read-dynamic-include-regex = "cmdb.cc_ObjectBase_(.*)_pub_|cmdb.cc_ApplicationBase$|cmdb.cc_HostBase$|cmdb.cc_ObjDes$"

# if you want to use MongoDB change streams instead of legacy oplog tailing use change-stream-namespaces
# change streams require at least MongoDB API 3.6+
# if you have MongoDB 4+ you can listen for changes to an entire database or entire deployment
# in this case you usually don't need regexes in your config to filter collections unless you target the deployment.
# to listen to an entire db use only the database name.  For a deployment use an empty string.
change-stream-namespaces = [""]
namespace-regex = "cmdb.cc_ObjectBase_(.*)_pub_|cmdb.cc_ApplicationBase$|cmdb.cc_HostBase$|cmdb.cc_ObjDes$"

# additional settings

# compress requests to Elasticsearch
gzip = true
# use the following user name for Elasticsearch basic auth
elasticsearch-user = ""
# use the following password for Elasticsearch basic auth
elasticsearch-password = ""
# use 4 go routines concurrently pushing documents to Elasticsearch
elasticsearch-max-conns = 4 
# propagate dropped collections in MongoDB as index deletes in Elasticsearch
dropped-collections = true
# propagate dropped databases in MongoDB as index deletes in Elasticsearch
dropped-databases = true
# resume processing from a timestamp saved in a previous run
resume = true
# do not validate that progress timestamps have been saved
resume-write-unsafe = false
# override the name under which resume state is saved
resume-name = "default"
# use a custom resume strategy (tokens) instead of the default strategy (timestamps)
# tokens work with MongoDB API 3.6+ while timestamps work only with MongoDB API 4.0+
resume-strategy = 0
# print detailed information including request traces
verbose = true


enable-patches = true


# mapping settings

[[mapping]]
namespace = "cmdb.cc_ApplicationBase"
index = "cmdb.cc_applicationbase"

[[mapping]]
namespace = "cmdb.cc_HostBase"
index = "cmdb.cc_hostbase"

[[mapping]]
namespace = "cmdb.cc_ObjDes"
index = "cmdb.cc_objdes"

[[script]]
#namespace = "cmdb_oa.cc_ObjectBase_pub_"
# this script does not declare a namespace
# it is global to all collections
routing = true
script = """
var re = new RegExp("cmdb.cc_ObjectBase_(.*)_pub_")
module.exports = function(doc, ns, updateDesc) {
    if(re.test(ns)) {
        doc["_meta_monstache"] = {"index":"cmdb.cc_objectbase"};
    }
    // the doc namespace e.g. test.test is passed as the 2nd arg
        // if available, an object containing the update description is passed as the 3rd arg
    return  doc
}
"""

```
添加新的 direct-read-namespaces，change-stream-namespaces 需要添加对应的 mapping。

**启动：**

```shell
nohup ./monstache -f config.toml &
```

**检查：**

```shell
> curl 'localhost:9200/_cat/indices?v'
health status index                   uuid                   pri rep docs.count docs.deleted store.size pri.store.size
yellow open   cmdb.cc_objdes          nIPMWSqsRN6Y4RlZIUZyKw   1   1         10            0       12kb           12kb
yellow open   cmdb.cc_hostbase        R3uXSNHbR4iFNI0YOl_X3Q   1   1         39            0     17.9kb         17.9kb
yellow open   cmdb.cc_applicationbase aFjTbeiTQMKcqIyDMqBtUA   1   1        749            0    158.5kb        158.5kb
yellow open   cmdb.cc_objectbase      c_G-N4_XTp--uqqRzQ4PJQ   1   1          2            0     10.4kb         10.4kb
```

如果 MongoDB 与上述 ES index 对应集合中数据为空，无法自行创建 index，请自行创建空 index。
```shell
# 例：
> curl -XPUT http://localhost:9200/cmdb.cc_objectbase
```

### 8. 部署CMDB

编译后下载 **cmdb.tar.gz**

在目标机上解压包解**cmdb.tar.gz**，解压后根目录结构如下：

``` shell
drwxr-xr-x 5 root root  4096 Jun 18 15:24 cmdb_adminserver
drwxr-xr-x 4 root root  4096 Jun 18 15:24 cmdb_apiserver
drwxr-xr-x 4 root root  4096 Jun 18 15:24 cmdb_coreservice
drwxr-xr-x 5 root root  4096 Jun 18 15:24 cmdb_datacollection
drwxr-xr-x 4 root root  4096 Jun 18 15:24 cmdb_eventserver
drwxr-xr-x 5 root root  4096 Jun 18 15:24 cmdb_hostserver
drwxr-xr-x 4 root root  4096 Jun 18 15:24 cmdb_operationserver
drwxr-xr-x 5 root root  4096 Jun 18 15:24 cmdb_procserver
drwxr-xr-x 3 root root  4096 Jun 18 10:33 cmdb_synchronizeserver
drwxr-xr-x 5 root root  4096 Jun 18 15:24 cmdb_taskserver
drwxr-xr-x 4 root root  4096 Jun 18 15:24 cmdb_toposerver
drwxr-xr-x 4 root root  4096 Jun 18 15:24 cmdb_webserver
drwxr-xr-x 2 root root  4096 Jun 18 10:33 docker
-rwxr--r-- 1 root root   913 Jun 18 10:33 image.sh
-rwxr-xr-x 1 root root 19311 Jun 18 10:33 init.py
-rwxr--r-- 1 root root   372 Jun 18 15:20 init_db.sh
-rwxr--r-- 1 root root   211 Jun 18 10:33 ip.py
-rwxr-xr-x 1 root root    34 Jun 18 10:33 restart.sh
-rwxr-xr-x 1 root root  1231 Jun 18 10:33 start.sh
-rwxr-xr-x 1 root root   810 Jun 18 10:33 stop.sh
drwxr-xr-x 2 root root  4096 Jun 18 10:33 tool_ctl
-rwxr-xr-x 1 root root  2144 Jun 18 10:33 upgrade.sh
drwxr-xr-x 7 root root  4096 Jun 18 10:33 web
```

各目录代表的服务及职责：

| 目标                   | 类型       | 用途描述                                                                               |
| ---------------------- | ---------- | -------------------------------------------------------------------------------------- |
| cmdb_adminserver       | server     | 负责系统数据的初始化以及配置管理工作                                                   |
| cmdb_apiserver         | server     | 场景层服务，api 服务                                                                   |
| cmdb_coreservice       | server     | 资源管理层，提供原子接口服务                                                           |
| cmdb_datacollection    | server     | 场景层服务，数据采集服务                                                               |
| cmdb_eventserver       | server     | 场景层服务，事件推送服务                                                               |
| cmdb_hostserver        | server     | 场景层服务，主机数据维护                                                               |
| cmdb_operationserver   | server     | 场景层服务，提供与运营统计相关功能服务                                                 |
| cmdb_procserver        | server     | 场景层服务，负责进程数据的维护                                                         |
| cmdb_synchronizeserver | server     | 场景层服务，数据同步服务                                                               |
| cmdb_taskserver        | server     | 场景层服务，异步任务管理服务                                                           |
| cmdb_toposerver        | server     | 场景层服务，负责模型的定义以及主机、业务、模块及进程等实例数据的维护                   |
| cmdb_webserver         | server     | web server 服务子目录                                                                  |
| docker                 | Dockerfile | 各服务的Dockerfile模板                                                                 |
| image.sh               | script     | 用于制作Docker镜像                                                                     |
| init.py                | script     | 用于初始化服务及配置项，在需要重置服务配置的时候也可以运行此脚本，按照提示输入配置参数 |
| init_db.sh             | script     | 初始化数据库的数据                                                                     |
| ip.py                  | script     | 查询主机真实的IP脚本                                                                   |
| restart.sh             | script     | 用于重启所有服务                                                                       |
| start.sh               | script     | 用于启动所有服务                                                                       |
| stop.sh                | script     | 用于停止所有服务                                                                       |
| tool_ctl               | ctl        | 管理小工具                                                                             |
| upgrade.sh             | script     | 用于全量升级服务进程                                                                   |
| web                    | ui         | CMDB UI 页面                                                                           |

### 9. 初始化

假定安装目录是 **/data/cmdb/**

进入安装目录并执行初始化脚本，**按照提示输入参数**。

``` shell
[root@SWEBVM000229 /data/cmdb]# python init.py

	usage:
	--discovery           <discovery>           the ZooKeeper server address, eg:127.0.0.1:2181
	--database           <database>             the database name, default cmdb
	--redis_ip           <redis_ip>             the redis ip, eg:127.0.0.1
	--redis_port         <redis_port>           the redis port, default:6379
	--redis_pass         <redis_pass>           the redis user password
	--mongo_ip           <mongo_ip>             the mongo ip ,eg:127.0.0.1
	--mongo_port         <mongo_port>           the mongo port, eg:27017
	--mongo_user         <mongo_user>           the mongo user name, default:cc
	--mongo_pass         <mongo_pass>           the mongo password
	--blueking_cmdb_url  <blueking_cmdb_url>    the cmdb site url, eg: http://127.0.0.1:8088 or http://bk.tencent.com
	--blueking_paas_url  <blueking_paas_url>    the blueking paas url, eg: http://127.0.0.1:8088 or http://bk.tencent.com
	--listen_port        <listen_port>          the cmdb_webserver listen port, should be the port as same as -c <cc_url> specified, default:8083
	--full_text_search   <full_text_search>     full text search function, off or on, default off
	--es_url             <es_url>               the elasticsearch listen url
 	--user_info          <user_info>            the system user info, user and password are combined by semicolon, multiple users are separated by comma. eg: user1:password1,user2:password2
```

**init.py 参数详解：**

| ZooKeeper地址       | 用途说明                                                                                                                                                                             | 必填                    | 默认值                  |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------- | ----------------------- |
| --discovery         | 服务发现组件，ZooKeeper 服务地址                                                                                                                                                     | 是                      | 无                      |
| --database          | 数据库名字                                                                                                                                                                           | mongodb 中数据库名      | 否                      | cmdb |
| --redis_ip          | Redis监听的IP                                                                                                                                                                        | 是                      | 无                      |
| --redis_port        | Redis监听的端口                                                                                                                                                                      | 否                      | 6379                    |
| --redis_pass        | Redis登陆密码                                                                                                                                                                        | 是                      | 无                      |
| --mongo_ip          | MongoDB服务监听的IP                                                                                                                                                                  | 是                      | 无                      |
| --mongo_port        | MongoDB端口                                                                                                                                                                          | 否                      | 27017                   |
| --mongo_user        | MongoDB中CMDB数据库用户名                                                                                                                                                            | 是                      | 无                      |
| --mongo_pass        | MongoDB中CMDB数据库用户名密码                                                                                                                                                        | 是                      | 无                      |
| --blueking_cmdb_url | 该值表示部署完成后,输入到浏览器中访问的cmdb 网址, 格式: http://xx.xxx.com:80, 用户自定义填写;在没有配置 DNS 解析的情况下, 填写服务器的 IP:PORT。端口为当前cmdb_webserver监听的端口。 | 是                      | 无                      |
| --blueking_paas_url | 蓝鲸PAAS 平台的地址，对于独立部署的CC版本可以不配置                                                                                                                                  | 否                      | 无                      |
| --listen_port       | cmdb_webserver服务监听的端口，默认是8083                                                                                                                                             | 是                      | 8083                    |
| --full_text_search  | 全文检索功能开关(取值：off/on)，默认是off，开启是on                                                                                                                                  | 否                      | off                     |
| --es_url            | elasticsearch服务监听url，默认是http://127.0.0.1:9200                                                                                                                                | 否                      | http://127.0.0.1:9200   |
| --auth_scheme       | 权限模式，web页面使用，可选值: internal, iam                                                                                                                                         | 否                      | internal                |
| --auth_enabled      | 是否采用蓝鲸权限中心鉴权                                                                                                                                                             | 否                      | false                   |
| --auth_address      | 蓝鲸权限中心地址                                                                                                                                                                     | auth_enabled 为真时必填 | https://iam.domain.com/ |
| --auth_app_code     | cmdb项目在蓝鲸权限中心的应用编码                                                                                                                                                     | auth_enabled 为真时必填 | bk_cmdb                 |
| --auth_app_secret   | cmdb项目在蓝鲸权限中心的应用密钥                                                                                                                                                     | auth_enabled 为真时必填 | xxxxxxx                 |
| --log_level         | 日志级别0-9, 9日志最详细                                                                                                                                                             | 否                      | 3                       |
| --register_ip       | 进程注册到zookeeper上的IP地址，可以是域名                                                                                                                                            | 否                      | 无                      |
| --user_info         | 登陆 web 页面的账号密码                                                                                                                                                              | 否                      | 无                      |

**注:init.py 执行成功后会自动生成cmdb各服务进程所需要的配置。**

**示例(示例中的参数需要用真实的值替换)：**

如果部署了用于全文检索的第6和第7步，如要开启全文检索功能把full_text_search的值置为on
``` shell
python init.py  \
  --discovery          127.0.0.1:2181 \
  --database           cmdb \
  --redis_ip           127.0.0.1 \
  --redis_port         6379 \
  --redis_pass         1111 \
  --mongo_ip           127.0.0.1 \
  --mongo_port         27017 \
  --mongo_user         cc \
  --mongo_pass         cc \
  --blueking_cmdb_url  http://127.0.0.1:8080 \
  --blueking_paas_url  http://paas.domain.com \
  --listen_port        8080 \
  --auth_scheme        internal \
  --auth_enabled       false \
  --auth_address       https://iam.domain.com/ \
  --auth_app_code      bk_cmdb \
  --auth_app_secret    xxxxxxx \
  --full_text_search   off \
  --es_url             http://127.0.0.1:9200 \
  --log_level          3 \
  --register_ip         cmdb.domain.com \
  --user_info admin:admin
```

### 10. init.py 生成的配置如下

配置文件的存储路径：{安装目录}/cmdb_adminserver/configures/

``` shell
-rw-r--r-- 1 root root 873 Jun 18 17:25 common.yaml
-rw-r--r-- 1 root root   0 Jun 18 15:20 extra.yaml
-rw-r--r-- 1 root root 580 Jun 18 15:20 migrate.yaml
-rw-r--r-- 1 root root 155 Jun 18 15:20 mongodb.yaml
-rw-r--r-- 1 root root 321 Jun 18 15:20 redis.yaml
``` 

配置文件目录：{安装目录}/cmdb_adminserver/configures

**注：由于MongoDB 从3.6开始更改了默认加密方式，所以如果安装的MongoDB的版本大于等于3.6，需要手动将以上配置文件中MongoDB的配置项中增加 mechanism=SCRAM-SHA-1**

> 配置文件mongodb小节增加mechanism 配置项示例如下

``` toml
[mongodb]
host=127.0.0.1
usr=cc
pwd=cc
database=cmdb
port=27017
maxOpenConns=3000
maxIDleConns=1000
mechanism=SCRAM-SHA-1
```

---

## 运行效果

### 1. 启动服务

``` shell
[root@SWEBVM000229 /data/cmdb]#  ./start.sh 
starting: cmdb_adminserver
starting: cmdb_apiserver
starting: cmdb_coreservice
starting: cmdb_datacollection
starting: cmdb_eventserver
starting: cmdb_hostserver
starting: cmdb_operationserver
starting: cmdb_procserver
starting: cmdb_taskserver
starting: cmdb_toposerver
starting: cmdb_webserver
root       209     1  5 08:27 pts/0    00:00:00 ./cmdb_adminserver --addrport=127.0.0.1:60004 --logtostderr=false --log-dir=./logs --v=3 --config=configures/migrate.conf
root       230     1  1 08:27 pts/0    00:00:00 ./cmdb_apiserver --addrport=127.0.0.1:8080 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       263     1  0 08:27 pts/0    00:00:00 ./cmdb_coreservice --addrport=127.0.0.1:50009 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181
root       284     1  1 08:27 pts/0    00:00:00 ./cmdb_datacollection --addrport=127.0.0.1:60005 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       305     1  4 08:27 pts/0    00:00:00 ./cmdb_eventserver --addrport=127.0.0.1:60009 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       326     1  3 08:27 pts/0    00:00:00 ./cmdb_hostserver --addrport=127.0.0.1:60001 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       445     1  4 08:27 pts/0    00:00:00 ./cmdb_operationserver --addrport=127.0.0.1:60011 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       642     1  7 08:27 pts/0    00:00:00 ./cmdb_procserver --addrport=127.0.0.1:60003 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       661     1 11 08:27 pts/0    00:00:00 ./cmdb_taskserver --addrport=127.0.0.1:60012 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181
root       724     1  6 08:27 pts/0    00:00:00 ./cmdb_toposerver --addrport=127.0.0.1:60002 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181 --enable-auth=false
root       937     1  0 08:27 pts/0    00:00:00 ./cmdb_webserver --addrport=127.0.0.1:8090 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181
process count should be: 11 , now: 11
```

**注：此处cmdb_test仅用作效果展示，非有效进程。**


### 2. 服务启动之后初始化数据库

``` shell
[root@SWEBVM000229 /data/cmdb]# bash ./init_db.sh
{"result":true,"bk_error_code":0,"bk_error_msg":"success","data":"migrate success"}
```
**注：以上输出表示初始化数据库成功，此步骤必需要所有cmdb进程成功启动后执行。**



### 3. 系统运行页面

**打开浏览器:** 数据cmdb_webserver 监听的地址，如本文档中示例服务监听的地址: http://127.0.0.1:8083

![image](../resource/img/page.png)



### 4. 停止服务

``` shell
[root@SWEBVM000229 /data/cmdb]# ./stop.sh
Running process count: 0
```

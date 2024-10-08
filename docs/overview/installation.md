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

官方发布的 **Linux Release** 包下载地址见[这里](https://github.com/TencentBlueKing/bk-cmdb/releases)。如果你想自已编译，具体的编译方法见[这里](source_compile.md)。

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
由于redis版本不同，若配置文件中没有注释直接添加密码行即可
##### b. 重启Redis
如果Redis已经配置为service服务，可以通过以下方式重启：
```json
service redis restart
```
若通过brew方式安装，可以通过以下方式重启：
```json
brew services start redis@redis版本
```
##### c. 登录验证

设置Redis认证密码后，客户端登录时需要使用-a参数输入认证密码,举例如下：

```shell
$ ./redis-cli -h 127.0.0.1 -p 6379 -a myPassword
127.0.0.1:6379> config get requirepass
1) "requirepass"
2) "myPassword"
```

看到类似上面的输出，说明Reids密码认证配置成功。

#### 2. 安装MongoDB后数据库配置

##### a. MongoDB 集群搭建（供参考,可按实际要求搭建MongoDB 集群）

1. 创建数据存放目录

```shell
mkdir -p ~/data/mongodb/cmdb
```

2. 创建配置文件

   主节点（Primary）

   ```
   # vim /etc/mongodb_cmdb.conf  //写入：
   
   # mongodb_cmdb.conf
   dbpath=/root/app/data/mongodb/mongodb_cmdb
   logpath=/root/app/data/mongodb/mongodb_cmdb.log
   pidfilepath=/root/app/data/mongodb/mongodb_cmdb.pid
   directoryperdb=true
   logappend=true
   replSet=rs0
   bind_ip=0.0.0.0
   port=27017
   oplogSize=100
   fork=true

   # 备注：以上配置信息仅供参考
   # 参数说明：
   # dbpath：存放数据目录
   # logpath：日志数据目录
   # pidfilepath：pid文件
   # directoryperdb：数据库是否分目录存放
   # logappend：日志追加方式存放
   # replSet：Replica Set的名字
   # bind_ip：mongodb绑定的ip地址
   # port：端口
   # fork：守护进程运行，创建进程
   ```
MongoDB官方资料：
   [配置文件选项](https://www.mongodb.com/docs/manual/reference/configuration-options/)
   [ReplicaSet配置](https://www.mongodb.com/docs/manual/reference/replica-configuration/)

3. 启动 mongodb，进入 mongodb 的 bin 目录，执行：
使用配置文件启动
   ```
   mongod -f /etc/mongodb_cmdb.conf
   ```
使用命令行方式启动
   ```
   mongod --dbpath /root/app/data/mongodb/mongodb_cmdb \
      --logpath /root/app/data/mongodb/mongodb_cmdb.log \
      --pidfilepath /root/app/data/mongodb/mongodb_cmdb.pid \
      --directoryperdb \
      --logappend \
      --replSet rs0 \
      --bind_ip 0.0.0.0 \
      --port 27017 \
      --oplogSize 100 \
      --fork
   ```
   ps：启动失败可查看对应日志文件排查问题

4. 配置集群后，执行`mongo`命令连接mongodb服务

   ```
   > cfg={ _id:"rs0", members:[ {_id:0,host:'IP:27017',priority:2}] };
   > rs.initiate(cfg)
   ```

   说明：
   cfg 名字可选，只要跟mongodb参数不冲突，rs0为集群名字，仅作展示，用户使用中可以根据实际情况自行配置，_id 为 Replica Set 名字，
   members 里面的优先级 priority 值高的为主节点，对于仲裁点一定要加上`arbiterOnly:true`，否则主备模式不生效，
   使集群cfg配置生效：`rs.initiate(cfg)`
   查看集群状态：`rs.status()`
   IP更改为实际ip地址

##### b. 创建数据库 cmdb 设置用户名和密码

接下来连接MongoDB服务后，根据需求执行以下命令:

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

蓝鲸CMDB针对需求场景采用定制化的Monstache组件，组件以及其插件SO请从指定的Release Package中获取。

插件基于Monstache v6.0.0+, 需要依赖Elasticsearch v7+和MongoDB v4.2+。

阅读[蓝鲸CMDB全文检索插件文档](../../src/tools/monstache/README.md), 按照指引进行安装部署。

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

| 目标                   | 类型       | 用途描述                                                     |
| ---------------------- | ---------- | ------------------------------------------------------------ |
| cmdb_adminserver       | server     | 负责系统数据的初始化以及配置管理工作                         |
| cmdb_apiserver         | server     | 场景层服务，api 服务                                         |
| cmdb_coreservice       | server     | 资源管理层，提供原子接口服务                                 |
| cmdb_datacollection    | server     | 场景层服务，数据采集服务                                     |
| cmdb_eventserver       | server     | 场景层服务，事件推送服务                                     |
| cmdb_hostserver        | server     | 场景层服务，主机数据维护                                     |
| cmdb_operationserver   | server     | 场景层服务，提供与运营统计相关功能服务                       |
| cmdb_procserver        | server     | 场景层服务，负责进程数据的维护                               |
| cmdb_synchronizeserver | server     | 场景层服务，数据同步服务                                     |
| cmdb_taskserver        | server     | 场景层服务，异步任务管理服务                                 |
| cmdb_toposerver        | server     | 场景层服务，负责模型的定义以及主机、业务、模块及进程等实例数据的维护 |
| cmdb_webserver         | server     | web server 服务子目录                                        |
| docker                 | Dockerfile | 各服务的Dockerfile模板                                       |
| image.sh               | script     | 用于制作Docker镜像                                           |
| init.py                | script     | 用于初始化服务及配置项，在需要重置服务配置的时候也可以运行此脚本，按照提示输入配置参数 |
| init_db.sh             | script     | 初始化数据库的数据                                           |
| ip.py                  | script     | 查询主机真实的IP脚本                                         |
| restart.sh             | script     | 用于重启所有服务                                             |
| start.sh               | script     | 用于启动所有服务                                             |
| stop.sh                | script     | 用于停止所有服务                                             |
| tool_ctl               | ctl        | 管理小工具                                                   |
| upgrade.sh             | script     | 用于全量升级服务进程                                         |
| web                    | ui         | CMDB UI 页面                                                 |

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

| ZooKeeper地址       | 用途说明                                                     | 必填                    | 默认值                  |
| ------------------- | ------------------------------------------------------------ | ----------------------- | ----------------------- |
| --discovery         | 服务发现组件，ZooKeeper 服务地址                             | 是                      | 无                      |
| --database          | 数据库名字                                                   | mongodb 中数据库名      | 否                      |
| --redis_ip          | Redis服务的IP                                                | 是                      | 无                      |
| --redis_port        | Redis服务的端口                                              | 否                      | 6379                    |
| --redis_pass        | Redis登陆密码                                                | 是                      | 无                      |
| --mongo_ip          | MongoDB服务监听的IP                                          | 是                      | 无                      |
| --mongo_port        | MongoDB端口                                                  | 否                      | 27017                   |
| --mongo_user        | MongoDB中CMDB数据库用户名                                    | 是                      | 无                      |
| --mongo_pass        | MongoDB中CMDB数据库用户名密码                                | 是                      | 无                      |
| --blueking_cmdb_url | 该值表示部署完成后,输入到浏览器中访问的cmdb 网址, 格式: http://xx.xxx.com:80, 用户自定义填写;在没有配置 DNS 解析的情况下, 填写服务器的 IP:PORT。端口为当前cmdb_webserver监听的端口。 | 是                      | 无                      |
| --blueking_paas_url | 蓝鲸PAAS 平台的地址，对于独立部署的CC版本可以不配置          | 否                      | 无                      |
| --listen_port       | cmdb_webserver服务监听的端口，默认是8083                     | 是                      | 8083                    |
| --full_text_search  | 全文检索功能开关(取值：off/on)，默认是off，开启是on          | 否                      | off                     |
| --es_url            | elasticsearch服务监听url，默认是http://127.0.0.1:9200        | 否                      | http://127.0.0.1:9200   |
| --auth_scheme       | 权限模式，web页面使用，可选值: internal, iam                 | 否                      | internal                |
| --auth_enabled      | 是否采用蓝鲸权限中心鉴权                                     | 否                      | false                   |
| --auth_address      | 蓝鲸权限中心地址                                             | auth_enabled 为真时必填 | https://iam.domain.com/ |
| --auth_app_code     | cmdb项目在蓝鲸权限中心的应用编码                             | auth_enabled 为真时必填 | bk_cmdb                 |
| --auth_app_secret   | cmdb项目在蓝鲸权限中心的应用密钥                             | auth_enabled 为真时必填 | xxxxxxx                 |
| --log_level         | 日志级别0-9, 9日志最详细                                     | 否                      | 3                       |
| --register_ip       | 进程注册到zookeeper上的IP地址，可以是域名                    | 否                      | 无                      |
| --user_info         | 登陆 web 页面的账号密码                                      | 否                      | 无                      |

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
  --blueking_cmdb_url  http://127.0.0.1:80 \
  --blueking_paas_url  http://127.0.0.1:80 \
  --listen_port        80 \
  --auth_scheme        internal \
  --auth_enabled       false \
  --full_text_search   off \
  --es_url             http://127.0.0.1:9200 \
  --log_level          3 \
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

启动前检查配置，`vim cmdb_adminserver/configures/common.yaml`命令进入 common.yaml ，如下：
1. 检查输入到浏览器访问的cmdb地址和登录地址这两项是否正确。 
```yaml
webServer:
   site :
      #该值表示部署完成后,输入到浏览器中访问的cmdb 网址
      domainUrl: http://127.0.0.1:80
      #登录地址
      bkLoginUr: http://127.0.0.1/login/?app_id=%s&c_url=%s
      appCode: cc
```
2. 检查登录模式
```yaml
webServer:
  login:
     #登录模式
     version: opensource
``` 
登录模式可选值：
- `opensource` 代表跳转到 CMDB登录页面进行登录，需要对账户密码进行配置，可在使用 init.py 初始化配置文件时指定参数 " --user_info 账号:密码 " ，或者找到以下配置项进行配置：
```yaml
webServer:
   session:
      #账号密码，以 : 分割
      userInfo: 账号:密码
```
- `skip-login` 代表不需要进行登陆操作
- `blueking` 代表通过「蓝鲸统一登录」进行登录


`mongodb.yaml 和 redis.yaml`等配置也要确保与实际部署的 mongodb 和 redis 服务配置相同，不同处手动修改，以下为示例配置：
mongodb.yaml:
```yaml
mongodb:
   host: 127.0.0.1:27017
   port: 27017
   usr: cc
   pwd: "cc"
   database: cmdb
   maxOpenConns: 3000
   maxIdleConns: 100
   mechanism: SCRAM-SHA-1
   rsName: rs0
   #mongo的socket连接的超时时间，以秒为单位，默认10s，最小5s，最大30s。
   socketTimeoutSeconds: 10
   # mongodb事件监听存储事件链的mongodb配置
watch:
   host: 127.0.0.1:27017
   port: 27017
   usr: cc
   pwd: "cc"
   database: cmdb
   maxOpenConns: 10
   maxIdleConns: 5
   mechanism: SCRAM-SHA-1
   rsName: rs0
   socketTimeoutSeconds: 10
```
redis.yaml：
```yaml
redis:
   #公共redis配置信息,用于存取缓存，用户信息等数据
   host: 127.0.0.1:6379
   pwd: "cc"
   sentinelPwd: ""
   database: "0"
   maxOpenConns: 3000
   maxIDleConns: 1000
   #以下几个redis配置为datacollection模块所需的配置,用于接收第三方提供的数据
   #接收主机信息数据的redis
   snap:
      host: 127.0.0.1:6379
      pwd: "cc"
      sentinelPwd: ""
      database: "0"
   #接收模型实例数据的redis
   discover:
      host: 127.0.0.1:6379
      pwd: "cc"
      sentinelPwd: ""
      database: "0"
   #接受硬件数据的redis
   netcollect:
      host: 127.0.0.1:6379
      pwd: "cc"
      sentinelPwd: ""
      database: "0"
```
确认配置无误后启动服务：

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
root       937     1  0 08:27 pts/0    00:00:00 ./cmdb_webserver --addrport=127.0.0.1:80 --logtostderr=false --log-dir=./logs --v=3 --regdiscv=127.0.0.1:2181
process count should be: 11 , now: 11
```

**注：cmdb_authserver需要依赖蓝鲸体系中的权限中心平台，如果启动失败属于正常现象。此处cmdb_test仅用作效果展示，非有效进程。**


### 2. 服务启动之后初始化数据库

```shell
[root@SWEBVM000229 /data/cmdb]# bash ./init_db.sh
{"result":true,"bk_error_code":0,"bk_error_msg":"success","data":"migrate success"}
```

**注：以上输出表示初始化数据库成功，此步骤必需要所有cmdb进程成功启动后执行。**



### 3. 系统运行页面

**打开浏览器:** 数据cmdb_webserver 监听的地址，如本文档中示例服务监听的地址: http://127.0.0.1:80

![image](../resource/img/page.png)



### 4. 停止服务

```shell
[root@SWEBVM000229 /data/cmdb]# ./stop.sh
Running process count: 0
```

## 常见问题

### 1. authserver服务无法正常启动

- authserver启动需要依赖于第三方[paas](https://github.com/TencentBlueKing/legacy-bk-paas)的权限中心系统，除了没有鉴权逻辑外， 不影响cmdb的单独部署使用

### 2. 运行./init_db.sh出现 ReplicaSetNoPrimary 错误

- 可参考 [issue](https://github.com/TencentBlueKing/bk-cmdb/issues/6155) 方法解决该问题，上述文档也有参考配置步骤，可对照查看是否遗漏某些步骤
- 如果 mongdb 服务运行正常，但无法连接，可检查`cmdb_adminserver/configures/mongodb.yaml`的配置与实际部署的 mongodb 服务配置是否相同，不同处手动修改，redis 服务同理

### 3. 服务启动失败了，如何排查

- 查看服务logs目录下的std.log以及xx.ERROR日志文件，根据里面的日志，定位到无法正常启动原因

### 4. 服务启动成功，但无法访问

- 启动前检查配置，`vim cmdb_adminserver/configures/common.yaml`命令进入 common.yaml ，检查输入到浏览器访问的cmdb地址和登录地址这两项是否正确。
- 查看ip.py文件具体生成的访问ip,由于电脑版本不同可能导致获取到一些无法访问的ip,可以将ip.py末尾的print(localhost)改为print("127.0.0.1")

### 其他问题

- 查看 [cmdb项目issues地址](https://github.com/TencentBlueKing/bk-cmdb/issues) ,寻找相同或类似问题的解决办法
- 创建issue, 带上版本号+错误日志文件+配置文件等信息，我们看到后，会第一时间为您解答；
- 同时我们也鼓励，有问题互相解答，提PR参与开源贡献，共建开源社区。

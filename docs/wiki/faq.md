---
### cmdb 对于操作系统是否限定为centos？
> 由于我们的开发、测试、部署都是在centos环境，推荐开发者在centos环境进行开发与测试，但go语言是跨平台的语言，也可以尝试在其它系统上进行开发测试。

### zookeeper对系统是否是必须的？
> 是的，我们的配置中心服务、服务发现服务都依赖zookeeper 实现，后续我们也会考虑以etcd 为载体。

### cmdb支持独立部署吗？
> 常规情况下cmdb 是为整个蓝鲸体系服务的，但是我们也支持独立部署版本，只需要在web_server中修改配置项即可。独立部署版本是为了方便大家体验。

### cmdb 用到了哪些golang的开发框架
> cmdb 主要以go http 框架为主，用到了gin和go-restful框架，关于这两款框架的使用，可以加入我们的技术交流群大家交流

### 社区版功能是不是比开源版的丰富呢？

**答疑：** 

> 1. 两个CMDB都是一样的版本功能。
> 2. 蓝鲸社区版软件包包含了蓝鲸各个产品以及官方saas，CMDB是集成到社区版的。
> 3. github上的开源版的是可以独立部署，脱离蓝鲸社区版的，可以自定义开发的。
> 4. 开源版CMDB的release版本与社区版内置的CMDB保持一致并同步更新。


### init_db.sh 执行失败要怎么排查？

**答疑：** 
> 1. 检查进程数量，ps -ef| grep cmdb | grep -v grep | wc -l， 如果所有的cmdb的进程都在运行，那么进程数量应该是12个；
> 2. 检查cmdb_adminserver 进程是否存在；
> 3. 检查cmdb_adminserver 的日志；
> 4. 检查init.py 做初始化时配置的 zk、redis、mongodb 的地址是否是其真正监听的地址。**特别说明：即使这些组件都是与cmdb进程同机部署的，也不可以在做init 的时候将地址配置为 127.0.0.1。**


### 是否支持容器部署？

**答疑：** 
> 1. 当前版本没有提供官方的容器镜像
> 2. 为了满足容器部署的需要，当前仅提供了 image.sh 脚本方便有需要的用户自行快速构建镜像。
> 3. image.sh 此文件在生成的安装包的根目录可以找到。

### 是否有用户管理？

**答疑：** 
> 当前开源版本是免登陆版本，因此对用户管理功能做了屏蔽。


### nodejs 编译通不过？

**答疑：** 
> 1. 对于大陆用户需要经过特殊配置才可以使用nodejs 进行编译。
> 2. 源码编译请参考：https://github.com/Tencent/bk-cmdb/blob/master/docs/overview/source_compile.md


### 可以在虚拟机里安装吗？

**答疑：** 
> 1. 这个是没有做限制的，但是对于虚拟机里的网络环境需要进行正确的设置，否则可能会出现无法正确配置编译环境或者部署后无法在虚拟机外部访问bk-cmdb的问题。
> 2. 对于最小化部署单机即可。

### bk-cmdb 是采用什么技术架构开发的？

**答疑：** 
> 1. bk-cmdb是采用微服务架构设计实现。
> 2. 后端代码采用golang语言进行开发。
> 3. 前段是采用vue.js 框架构建的。

### bk-cmdb 有没有在线体验地址？

**答疑：** 
> 当前没有提供体验地址。

### init.py 进行初始化，为什么没有生成配置文件？

**答疑：** 
> 1. 当init.py 正确执行后会看到 如下的输出

``` text
initial configurations success, configs could be found at cmdb_adminserver/configures
```
> 2. 如果没有看到请查看屏显错误提示。
> 3. 请确认执行脚本的目录是否在安装目录的根目录。

### 日志里看到如下的内容是什么含义？

``` text
fail to watch children for path(/cc/services/endpoints/XXXX), err:zk: node does not exist
```

**答疑：** 
> 1. 请确认ZooKeeper监听的IP及端口于执行./init.py为cmdb服务配置的是否一致；
> 2. 检查XXX服务进程是否处于运行状态；

### 进程在运行，但是页面打不开，应该如何排查问题？

**答疑：** 
> 1. 执行./init.py为cmdb服务配置的listen_port与--blueking_cmdb_url 指定的地址（或域名所映射的地址）所包含的端口是否一致；
> 2. 如果第一检查发现端口不一致，那么需要重新进行正确的初始化配置，初始化后要重启所有服务进程。

### ZooKeeper 在cmdb 系统中起到什么作用？

**答疑：** 
> 1. 用于bk-cmdb内部的服务发现。
> 2. 用于bk-cmdb的系统配置存储。

### 重新执行 init.py 后，新刷入的配置为什么没有生效？

**答疑：** 
> 1. 执行init.py之后需要重启bk-cmdb的所有服务进程。

### 执行init.py的时候出现以下输出是什么意思？

``` text
	option --listen_portd not recognized

	usage: 
	-discovery           <discovery>            the ZooKeeper server address, eg:127.0.0.1:2181 
	--database           <database>             the database name, default cmdb 
	--redis_ip           <redis_ip>             the redis ip, eg:127.0.0.1 
	--redis_port         <redis_port>           the redis port, default:6379 
	--redis_pass         <redis_pass>           the redis user password 
	--mongo_ip           <mongo_ip>             the mongo ip ,eg:127.0.0.1 
	--mongo_port         <mongo_port>           the mongo port, eg:27017 
	--mongo_user         <mongo_user>           the mongo user name, default:cc 
	--mongo_pass         <mongo_pass>           the mongo password 
	--blueking_cmdb_url  <blueking_cmdb_url>    the cmdb site url, eg: http://127.0.0.1:8088 or http://bk.tencent.com 
	--blueking_paas_url  <blueking_paas_url>    the blueking pass url, eg: http://127.0.0.1:8088 or http://bk.tencent.com 
	--listen_port        <listen_port>          the cmdb_webserver listen port, should be the port as same as -c <blueking_cmdb_url> specified, default:8083
```

**答疑：** 
> 1. option --listen_portd not recognized 参数为非指定参数，不能被识别
> 2. usage 是正确的参数列表及参数含义
> 3. 详细的使用请参考：https://github.com/Tencent/bk-cmdb/blob/master/docs/overview/installation.md

### 日志里看到以下内容如何处理？

``` text
fail to get configure, will get again
```

**参考：** https://github.com/Tencent/bk-cmdb/issues/67

### 除了chrome 浏览器之外，其他浏览器会支持吗？

**答疑：** 
> 1. 后续的迭代中会陆续支持更多主流浏览器。

### 系统启动成功，日志无错误， 在浏览器中中访问显示空白

**答疑：** 
由于使用前端页面使用vue框架, 不可使用chrome 版本小于51.00.00

**答疑：** 
> 1. 后续的迭代中会陆续支持更多主流浏览器。

### bk-cmdb 2.0 的数据要如何升级到bk-cmdb 3.0？

**答疑：** 
> 1. 目前已经有2.0版本迁移到3.0版本的工具，可以加我们的开源版本QQ群获取。
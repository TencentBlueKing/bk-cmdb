# V0.03R060D54.0.1
## 修改
* 调整ZK路径使用，只使用1种路径
* 原来有2个log，每次启动程序会打印2份，现在修改为只使用1份log，log配置维持在本地
* 去除dataop写数据失败，保存本地文件功能
## 修复
* 修复free object的问题
* 修复dataop功能，企业版中使用V2路径，内部版继续使用V1路径，统一ZK配置解析
* 修复dataid type和storage type对应转换的问题

# V0.03R060D53.0.2
## 修复
* 转发某一存储失败的时候不影响转发其余存储

# V0.03R060D53.0.1
## 修复
* 内存泄漏问题

# V0.03R060D53
## 新增
* 新增redis publish功能，直接连接redis写数据，默认宏关闭
* 新增多存储转发功能，新增ZK上dataid节点格式
## 修改
* 修改日志
* 不再转发redis的实时数据
* 代码调整，提取一些公共函数，变量
* 调整配置和解析。去除basecfg配置项
* ops发送使用127.0.0.1的ip

# V0.03R060D52.0.1
## 修复
* 修复连接ZK超时后无法感知zk的问题

# V0.03R060D49.0.2
### 修复
* 无内网IP，proxy启动失败的问题

# V0.03R060D49.0.1
### 修改
* 获取zk的DB连接端口修改

# V0.03R060D48.0.2
### 修改
* 去除灯塔转发模式
* 清理代码
### 新增
* 写ZK新增ACL验证，新节点根路径在gsev2

# V0.03R060D43.0.1
### 修改
* 发送心跳版本信息由上报redis调整为上报到ops
* 去除连接心跳dbporxy
### 修复
* 加强Json解析检验
* dbproxy死锁问题（去除heart dbproxy）
* 修复动态协议问题

# V0.03R060D40.0.0
### 新增
* 合入ts功能，proxy，beacon，tglog模式
* 接收和发送端可以分别设置SSL
### 修改
* 调整glog日志大小为20*10MB
* beacon, tglog use timecenter
* 修改了部分运营数据，并且其中的累计值修改为每分钟的数值

# V0.03R060D38.0.4
### 修复
* redis节点初始化问题，会导致实时数据无法写入

# V0.03R060D38.0.3
### 新增
* 证书，运行时等路径可配
* 新增 http 工作模式，提供 restapi 用于上报数据，runmode 为 3
### 修改
* 发送数据大小上限修改为3000000B

# V0.03R060D38.0.2
### 修改
* 完善监控打点
* 去除UdpClient
* dataid to dataId
### 修复
* logical id 设置默认值为0，保持和agent行为一致 

# V0.03R060D38.0.1
### 修改
* ops 模式去除向ops server发送运营数据

# V0.03R060D38.0.0
### 修复
* ops 模式修复读取的dsip错误的问题
# V0.03R060D37.0.1
### 修复
* 在SSL接收数据的时候，增加保护代码
* 修改断链时间，从1分钟调整为5分钟

# V0.03R060D37.0.0
### 新增
* DataServer 增加监控打点数据
* DataServer OPS模式增加监控打点数据上报
### 修复
* 修复时间戳计算错误的bug

# V0.03R060D35.1.6
### 修复
* DataServer ssl 读数据不完整问题
* 支持配置外网IP
 
# V0.03R060D35.1.5
### 修复
* DataServerOps内存泄漏问题（很早之前引入）

# V0.03R060D35.1.4
### 修复
* 程序退出时资源重重复释放的问题

### 修改
* 移除server_id字段
* DataServer模式修改发送DataServerOps流水运营数据格式，移除timeout字段
* DataServerOps模式移除延时统计功能
* DataServer模式对发送失败的返回值分类，事件上报
* 优化高频时间获取的方案，提高整体性能

### 新增
* 自动创建data目录

# V0.03R060D35.1.0
### 修复
* 日志配置解析错误

### 新增
* 支持区分kafka集群和dataid类型
* rdkafka统计信息周期从10秒改为1分钟
* OpsServer模式支持上报数据到Kafka
* OpsServer模式发送运营数据的目标dataid可以在zk上配置
* OpsServer模式发送运营数据失败的内容会保存到本地
* DataServer模式rdkafka统计信息上报到远端
* DataServer模式增加0x09类型指定的DynamicalProtocol 相关协议的解析，以及与原数据相关的时间戳的解析
* DataServer模式原数据流经各个节点的时间戳入库kafka，并写入kafka 的key中


# V0.03R059D41
### 修复
* 更新zk上的存储信息逻辑错误的问题。该问题可能会造成后续数据转发失败，ds发生coredump。

# V0.03R059D40
### 修复
* rdkafka统计时，未调用poll导致的内存泄漏问题

# V0.03R059D39
### 修复
* 返回值错误 
* fd泄漏 
* 内存越界

# V0.03R059D38
### 修复
* 修复错误topic问题


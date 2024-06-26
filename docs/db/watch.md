# 资源变更事件相关表

## cc_{资源类型}WatchChain

#### 作用

存放某类资源变更事件信息，变更事件相关的表有统一的结构，此处进行归类，包含的表如下：

| 表名称                                  | 作用                   |
|--------------------------------------|----------------------|
| cc_ApplicationBaseWatchChain         | 存放业务变更事件信息           |
| cc_BizSetBaseWatchChain              | 存放业务集变更事件            |
| cc_bizSetRelationMixedWatchChain     | 存放业务集关系变更事件          |
| cc_ClusterBaseWatchChain             | 容器数据纳管——存放集群变更事件     |
| cc_HostBaseWatchChain                | 存放主机变更事件信息           |
| cc_HostIdentityMixedWatchChain       | 存放主机身份变更事件信息         |
| cc_InstAsstWatchChain                | 存放实例关联变更事件信息         |
| cc_MainlineInstanceWatchChain        | 存放主线实例变更事件信息         |
| cc_ModuleBaseWatchChain              | 存放模块变更事件信息           |
| cc_ModuleHostConfigWatchChain        | 存放模块主机关联关系变更事件信息     |
| cc_NamespaceBaseWatchChain           | 容器数据纳管——存放命名空间变更事件信息 |
| cc_NodeBaseWatchChain                | 容器数据纳管——存放节点变更事件信息   |
| cc_ObjectBaseWatchChain              | 存放模型变更事件             |
| cc_PlatBaseWatchChain                | 存放管控区域变更事件信息         |
| cc_ProcessInstanceRelationWatchChain | 存放进程实例关联关系变更事件信息     |
| cc_ProcessWatchChain                 | 存放进程变更事件信息           |
| cc_ProjectBaseWatchChain             | 存放项目变更事件信息           |
| cc_SetBaseWatchChain                 | 存放集群变更事件信息           |
| cc_SetTemplateWatchChain             | 存放集群模板变更事件信息         |
| cc_WorkloadBaseWatchChain            | 容器数据纳管——工作负载变更事件信息   |
| cc_PodBaseWatchChain                 | 容器数据纳管——Pod变更事件信息    |

#### 表结构

| 字段                  | 类型         | 描述      |
|---------------------|------------|---------|
| _id                 | ObjectId   | 数据唯一ID  |
| id                  | NumberLong | 自增ID    |
| cluster_time        | ISODate    | 事件时间    |
| oid                 | String     | 事件ID    |
| type                | String     | 操作类型    |
| token               | String     | 数据库恢复令牌 |
| cursor              | String     | cc的事件游标 |
| inst_id             | NumberLong | 实例ID    |
| bk_supplier_account | String     | 开发商ID   |

## cc_WatchToken

#### 作用

数据库表的监控令牌信息，监控令牌用于跟踪指定数据集合的变化情况

#### 表结构

| 字段            | 类型         | 描述   |
|---------------|------------|------|
| _id           | String     | 数据库表 |
| id            | NumberLong | 自增id |
| token         | String     | 令牌   |
| cursor        | String     | 事件游标 |
| start_at_time | ISODate    | 开始时间 |
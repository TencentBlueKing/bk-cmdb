# 容器数据纳管功能相关表

## cc_ClusterBase

#### 作用

容器数据纳管——集群表

#### 表结构

| 字段                   | 类型         | 描述                       |
|----------------------|------------|--------------------------|
| _id                  | ObjectId   | 数据唯一ID                   |
| bk_biz_id            | NumberLong | 业务id                     |
| bk_cluster_name      | String     | 集群名称                     |
| bk_scheduling_engine | String     | 调度引擎                     |
| bk_uid               | String     | 集群ID                     |
| bk_tke_cluster_id    | String     | 集群在TKE服务中的唯一标识符          |
| bk_cluster_version   | String     | 集群版本                     |
| bk_network_type      | String     | 集群网络类型， overlay或underlay |
| bk_region            | String     | 所属地域                     |
| bk_vpc               | String     | vpc网络                    |
| bk_cluster_network   | String     | 集群网络                     |
| bk_cluster_type      | String     | 集群类型                     |
| create_time          | ISODate    | 创建时间                     |
| last_time            | ISODate    | 最后更新时间                   |

## cc_ContainerBase

#### 作用

容器数据纳管——容器表

#### 表结构

| 字段                  | 类型         | 描述         |
|---------------------|------------|------------|
| _id                 | ObjectId   | 数据唯一ID     |
| bk_pod_id           | NumberLong | pod在cc中的id |
| container_uid       | String     | 容器ID       |
| image               | String     | 镜像信息       |
| name                | String     | 容器名称       |
| ports               | Array      | 容器端口       |
| host_ports          | Array      | 主机端口映射     |
| args                | Array      | 启动参数       |
| started             | NumberLong | 启动时间       |
| requests            | Object     | 申请资源大小     |
| limits              | Object     | 资源限制       |
| liveness            | Object     | 存活探针       |
| environment         | Array      | 环境变量       |
| mounts              | Array      | 挂载卷        |
| bk_supplier_account | String     | 开发商ID      |
| create_time         | ISODate    | 创建时间       |
| last_time           | ISODate    | 最后更新时间     |

#### requests 和 limits 字段结构示例

| 字段  | 类型     | 描述                          |
|-----|--------|-----------------------------|
| cpu | String | 资源名称，字段名称由用户自定义，此处 cpu 仅为例子 |

#### requests.cpu 和 limits.cpu 字段结构示例

| 字段     | 类型     | 描述         |
|--------|--------|------------|
| format | String | 申请或限制的资源大小 |

#### liveness 字段结构示例

| 字段                            | 类型         | 描述                    |
|-------------------------------|------------|-----------------------|
| exec                          | Object     | 要在容器内执行的命令行操作         |
| httpGet                       | Object     | 基于 HTTP GET 请求的存活探针配置 |
| tcpSocket                     | Object     | 基于 TCP Socket 的存活探针配置 |
| grpc                          | Object     | 基于 gRPC 的存活探针配置       |
| initialDelaySeconds           | NumberLong | 容器启动后等待多少秒后开始执行探测     |
| timeoutSeconds                | NumberLong | 探测的超时时间               |
| periodSeconds                 | NumberLong | 连续探测之间的时间间隔           |
| successThreshold              | NumberLong | 连续多少次成功的探测才认为容器健康     |
| failureThreshold              | NumberLong | 连续多少次失败的探测才认为容器不健康    |
| terminationGracePeriodSeconds | NumberLong | 探测失败后容器需要终止的等待时间      |

#### liveness.exec 字段结构示例

| 字段      | 类型           | 描述              |
|---------|--------------|-----------------|
| command | String Array | 要在容器内执行的命令行操作数组 |

#### liveness.httpGet 字段结构示例

| 字段          | 类型           | 描述          |
|-------------|--------------|-------------|
| path        | String       | 访问的路径       |
| port        | NumberLong   | 访问的端口       |
| host        | String       | 连接的主机名      |
| scheme      | String       | 连接主机时要使用的方案 |
| httpHeaders | Object Array | HTTP 请求头信息  |

#### liveness.httpGet.httpHeaders 字段结构示例

| 字段    | 类型     | 描述    |
|-------|--------|-------|
| name  | String | 头信息名称 |
| value | String | 头信息值  |

#### liveness.tcpSocket 字段结构示例

| 字段   | 类型         | 描述     |
|------|------------|--------|
| port | NumberLong | 访问的端口  |
| host | String     | 连接的主机名 |

#### liveness.grpc 字段结构示例

| 字段      | 类型         | 描述          |
|---------|------------|-------------|
| service | NumberLong | gRPC 服务     |
| port    | NumberLong | gRPC 服务的端口号 |

## cc_NamespaceBase

#### 作用

容器数据纳管——命名空间表

#### 表结构

| 字段                  | 类型           | 描述           |
|---------------------|--------------|--------------|
| _id                 | ObjectId     | 数据唯一ID       |
| bk_biz_id           | NumberLong   | 业务id         |
| cluster_uid         | String       | 集群在BCS表中的标识  |
| bk_cluster_id       | NumberLong   | 集群在cc表中的唯一标识 |
| name                | String       | 命名空间名称       |
| labels              | Object       | 标签           |
| resource_quotas     | Object Array | 命名空间资源配额信息   |
| bk_supplier_account | String       | 开发商ID        |
| create_time         | ISODate      | 创建时间         |
| last_time           | ISODate      | 最后更新时间       |

#### labels 字段结构示例

| 字段  | 类型     | 描述                      |
|-----|--------|-------------------------|
| env | String | 标签名称，由用户自定义，此处 env 仅为例子 |

#### resource_quotas 字段结构示例

| 字段             | 类型           | 描述           |
|----------------|--------------|--------------|
| hard           | Object       | 命名空间资源的请求与限制 |
| scopes         | String Array | 资源配额范围       |
| scope_selector | Object       | 范围选择器        |

#### resource_quotas.scope_selector 字段结构示例

| 字段                | 类型           | 描述     |
|-------------------|--------------|--------|
| match_expressions | Object Array | 匹配的表达式 |

#### resource_quotas.scope_selector.match_expressions 字段结构示例

| 字段         | 类型     | 描述             |
|------------|--------|----------------|
| scope_name | String | 选择器适用的作用域名称    |
| operator   | String | 匹配操作符          |
| values     | —      | 匹配的值，值类型取决于操作符 |

## cc_NodeBase

#### 作用

容器数据纳管——节点表

#### 表结构

| 字段                  | 类型         | 描述                           |
|---------------------|------------|------------------------------|
| _id                 | ObjectId   | 数据唯一ID                       |
| bk_biz_id           | NumberLong | 业务id                         |
| cluster_uid         | String     | 集群本身的id                      |
| bk_cluster_id       | NumberLong | 集群在cc表中的唯一标识                 |
| bk_host_id          | NumberLong | 主机在cc中的唯一标识                  |
| name                | String     | 节点名称                         |
| roles               | String     | 节点类型，master或者None            |
| labels              | Object     | 标签                           |
| unschedulable       | Boolean    | 是否关闭可调度，true为不可调度，false代表可调度 |
| internal_ip         | Array      | 内网IP                         |
| external_ip         | Array      | 外网IP                         |
| hostname            | String     | 主机名，与节点名一致                   |
| runtime_component   | String     | 运行时组件                        |
| kube_proxy_mode     | String     | 代理模式                         |
| pod_cidr            | String     | 此节点Pod地址的分配范围                |
| bk_supplier_account | String     | 开发商ID                        |
| create_time         | ISODate    | 创建时间                         |
| last_time           | ISODate    | 最后更新时间                       |

#### labels 字段结构示例

| 字段  | 类型     | 描述                      |
|-----|--------|-------------------------|
| env | String | 标签名称，由用户自定义，此处 env 仅为例子 |

## cc_NsSharedClusterRelation

#### 作用

容器数据纳管——共享集群关联关系表

#### 表结构

| 字段                  | 类型         | 描述               |
|---------------------|------------|------------------|
| _id                 | ObjectId   | 数据唯一ID           |
| bk_biz_id           | NumberLong | namespace所在的业务ID |
| bk_cluster_id       | NumberLong | 共享集群ID           |
| bk_namespace_id     | NumberLong | namespace ID     |
| bk_asst_biz_id      | NumberLong | 关联的平台业务ID        |
| bk_supplier_account | String     | 开发商ID            |

## cc_PodBase

#### 作用

容器数据纳管——Pod表

#### 表结构

| 字段                  | 类型         | 描述                                                                                                                                                          |
|:--------------------|:-----------|:------------------------------------------------------------------------------------------------------------------------------------------------------------|
| _id                 | ObjectId   | 数据唯一ID                                                                                                                                                      |
| id                  | NumberLong | 在cc表中Pod的唯一标识                                                                                                                                               |
| bk_biz_id           | NumberLong | 业务id                                                                                                                                                        |
| bk_supplier_account | String     | 供应商id                                                                                                                                                       |
| bk_cluster_id       | NumberLong | 集群在cc表中的唯一标识                                                                                                                                                |
| cluster_uid         | String     | bcs集群id                                                                                                                                                     |
| namespace           | String     | namespace name                                                                                                                                              |
| bk_namespace_id     | NumberLong | 所属的namespace在cc中的id                                                                                                                                         |
| bk_host_id          | NumberLong | 主机在cc中的唯一标识                                                                                                                                                 |
| node_name           | String     | 所属node名称                                                                                                                                                    |
| bk_reference_id     | NumberLong | 所属的namespace在cc中的id                                                                                                                                         |
| reference_kind      | String     | workload kind                                                                                                                                               |
| reference_name      | String     | 冗余的workload name                                                                                                                                            |
| name                | String     | pod 名称                                                                                                                                                      |
| priority            | NumberLong | 优先级                                                                                                                                                         |
| labels              | Object     | 标签，key和value均是string，官方文档：http://kubernetes.io/docs/user-guide/labels                                                                                       |
| ip                  | String     | 容器网络IP                                                                                                                                                      |
| ips                 | Array      | 容器网络IP数组，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#podip-v1-core                                                              |
| volumes             | Object     | 使用的卷信息，官方文档：https://kubernetes.io/zh/docs/concepts/storage/volumes/ ，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volume-v1-core |
| qos_class           | String     | 服务质量，官方文档：https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/quality-service-pod/                                                               |
| node_selectors      | Object     | 节点标签选择器，key和value均是string，官方文档：https://kubernetes.io/zh/docs/concepts/scheduling-eviction/assign-pod-node/                                                  |
| tolerations         | Object     | 容忍度，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#toleration-v1-core                                                              |
| create_time         | ISODate    | 创建时间                                                                                                                                                        |
| update_time         | ISODate    | 更新时间                                                                                                                                                        |

#### labels 字段结构示例

| 字段  | 类型     | 描述                      |
|-----|--------|-------------------------|
| env | String | 标签名称，由用户自定义，此处 env 仅为例子 |

## cc_{Workload类型}Base

#### 作用

Workload共分为八种类型：Deployment、DaemonSet、StatefulSet、GameStatefulSet、GameDeployments、CronJob、Job和Pod，覆盖各种容器应用场景。cmdb中对八种类型进行分类分表存放，共分为以下8个表：

| 表名称                    | 作用                                       |
|------------------------|------------------------------------------|
| cc_DeploymentBase      | 容器数据纳管——存放工作负载（Workload）类型为Deployment的数据 |
| cc_DaemonSetBase       | 容器数据纳管——存放工作负载（Workload）类型为DaemonSet的数据  |
| cc_StatefulSetBase     | 容器数据纳管——存放Workload类型为StatefulSet的数据      |
| cc_GameStatefulSetBase | 容器数据纳管——存放Workload类型为GameStatefulSet的数据  |
| cc_GameDeploymentBase  | 容器数据纳管——存放Workload类型为GameDeployments的数据  |
| cc_CronJobBase         | 容器数据纳管——存放Workload类型为CronJob的数据          |
| cc_JobBase             | 容器数据纳管——存放Workload类型为cc_JobBase的数据       |
| cc_PodWorkloadBase     | 容器数据纳管——存放Workload类型为Pod的数据              |

#### 表结构

| 字段                    | 类型         | 描述                                                                                                                                      |
|-----------------------|------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| _id                   | ObjectId   | 数据唯一ID                                                                                                                                  |
| id                    | NumberLong | 自增id                                                                                                                                    |
| name                  | String     | 工作负载名称                                                                                                                                  |
| namespace             | String     | 工作负载所属命名空间                                                                                                                              |
| type                  | String     | 用于区分不同的Workload在CC里同一个模型下的分类                                                                                                            |
| labels                | String     | 工作负载labels，官方文档：https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/labels/                                            |
| selector              | String     | 工作负载选择器，一般与labels同时使用，例如：Selector:   k8s-app=kube-dns，官方文档：https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/labels/ |
| replicas              | NumberLong | 工作负载实例个数，官方文档：https://kubernetes.io/zh/docs/concepts/workloads/controllers/replicaset/                                                  |
| strategyType          | String     | 升级策略，工作负载更新机制，详情请见：https://kubernetes.io/zh/docs/concepts/workloads/controllers/deployment/                                             |
| minReadySeconds       | NumberLong | 最小就绪时间，工作负载更新机制，详情请见：https://kubernetes.io/zh/docs/concepts/workloads/controllers/deployment/                                           |
| rollingUpdateStrategy | String     | 滚动更新策略，工作负载更新机制，详情请见：https://kubernetes.io/zh/docs/concepts/workloads/controllers/deployment/                                           |
| bk_supplier_account   | String     | 开发商ID                                                                                                                                   |
| create_time           | ISODate    | 创建时间                                                                                                                                    |
| last_time             | ISODate    | 最后更新时间                                                                                                                                  |
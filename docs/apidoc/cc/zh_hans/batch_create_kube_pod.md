### 功能描述

新建容器Pod及container(v3.12.1+，权限:容器pod的创建权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段   | 类型    | 必选  | 描述           |
|------|-------|-----|--------------|
| data | array | 是   | 所要创建pod的详细信息 |

#### data[x]

| 字段        | 类型    | 必选  | 描述              |
|-----------|-------|-----|-----------------|
| bk_biz_id | int   | 是   | 业务ID            |
| pods      | array | 是   | 此业务下要创建pod的详细信息 |

#### pods[x]

| 字段             | 类型           | 必选 | 描述           |
|----------------|--------------|----|--------------|
| spec           | object       | 是  | pod关联信息      |
| bk_host_id     | int          | 是  | pod关联host id |
| name           | string       | 是  | pod名称        |
| operator       | string array | 是  | pod负责人       |
| priority       | object       | 否  | 优先级          |
| labels         | object       | 否  | 标签           |
| ip             | string       | 否  | 容器网络IP       |
| ips            | array        | 否  | 容器网络IP数组     |
| volumes        | object       | 否  | 卷信息          |
| qos_class      | string       | 否  | 服务质量         |
| node_selectors | object       | 否  | 节点标签选择器      |
| tolerations    | object       | 否  | 容忍度          |
| containers     | array        | 否  | 容器信息         |

#### spec

| 字段              | 类型     | 必选  | 描述                 |
|-----------------|--------|-----|--------------------|
| bk_cluster_id   | int    | 是   | pod所在集群的ID         |
| bk_namespace_id | int    | 是   | pod所属于namespace的ID |
| bk_node_id      | int    | 是   | pod所在node的ID       |
| ref             | object | 是   | pod对应workload的相关信息 |

#### ref

| 字段   | 类型  | 必选  | 描述                        |
|------|-----|-----|---------------------------|
| kind | int | 是   | pod相关联的workload类别，具体类别见注意 |
| id   | int | 是   | pod相关联的workload的ID        |

#### containers[x]

| 字段            | 类型     | 必选  | 描述     |
|---------------|--------|-----|--------|
| name          | string | 是   | 容器名称   |
| container_uid | string | 是   | 容器ID   |
| image         | string | 否   | 镜像信息   |
| ports         | array  | 否   | 容器端口   |
| host_ports    | array  | 否   | 主机端口映射 |
| args          | array  | 否   | 启动参数   |
| started       | int    | 否   | 启动时间   |
| limits        | object | 否   | 资源限制   |
| requests      | object | 否   | 申请资源大小 |
| liveness      | object | 否   | 存活探针   |
| environment   | array  | 否   | 环境变量   |
| mounts        | array  | 否   | 挂载卷    |

#### ports[x]

| 字段            | 类型     | 必选  | 描述   |
|---------------|--------|-----|------|
| name          | string | 是   | 端口名称 |
| hostPort      | int    | 否   | 主机端口 |
| containerPort | int    | 否   | 容器端口 |
| protocol      | string | 否   | 协议名称 |
| hostIP        | string | 否   | 主机IP |

#### liveness

| 字段        | 类型     | 必选  | 描述         |
|-----------|--------|-----|------------|
| exec      | object | 是   | 执行动作       |
| httpGet   | object | 否   | Http Get动作 |
| tcpSocket | object | 否   | tcp socket |
| grpc      | object | 否   | grpc 协议    |

**注意：**

- 一次性创建pod数量不超过200个。
- 具体的workload类别:deployment、statefulSet、daemonSet、gameStatefulSet、gameDeployment、cronJob、job、pods。
- 此接口会同步创建pod和对应的container。

### 请求参数示例

```json
 {
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "data": [
    {
      "bk_biz_id": 1,
      "pods": [
        {
          "spec": {
            "bk_cluster_id": 1,
            "bk_namespace_id": 1,
            "ref": {
              "kind": "deployment",
              "id": 1
            },
            "bk_node_id": 1
          },
          "name": "name",
          "operator": [
            "user1",
            "user2"
          ],
          "bk_host_id": 1,
          "priority": 1,
          "labels": {
            "env": "test"
          },
          "ip": "127.0.0.1",
          "ips": [
            {
              "ip": "127.0.0.1"
            },
            {
              "ip": "127.0.0.2"
            }
          ],
          "containers": [
            {
              "name": "name",
              "container_uid": "uid",
              "image": "xxx",
              "started": 1
            }
          ]
        }
      ]
    }
  ]
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "ids": [
      1,
      2
    ]
  },
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

**注意：**

- 返回的data中的podID数组顺序与参数中的数组数据顺序保持一致。

### 返回结果参数说明

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |
| request_id | string | 请求链ID                      |

### data

| 名称  | 类型    | 描述           |
|-----|-------|--------------|
| ids | array | 创建的容器podID列表 |

### 功能描述

新建容器节点(v3.12.1+，权限:容器节点的创建权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型    | 必选  | 描述            |
|-----------|-------|-----|---------------|
| bk_biz_id | int   | 是   | 业务ID          |
| data      | array | 是   | 需要创建的node具体信息 |

#### data[x]

| 字段                | 类型     | 必选  | 描述                           |
|-------------------|--------|-----|------------------------------|
| bk_cluster_id     | int    | 是   | 容器集群在cmdb中的ID                |
| bk_host_id        | int    | 是   | 关联的主机ID                      |
| name              | string | 是   | 节点名称                         |
| roles             | string | 否   | 节点类型                         |
| labels            | object | 否   | 标签                           |
| taints            | object | 否   | 污点                           |
| unschedulable     | bool   | 否   | 是否关闭可调度，true为不可调度，false代表可调度 |
| internal_ip       | array  | 否   | 内网IP                         |
| external_ip       | array  | 否   | 外网IP                         |
| hostname          | string | 否   | 主机名                          |
| runtime_component | string | 否   | 运行时组件                        |
| kube_proxy_mode   | string | 否   | kube-proxy 代理模式              |
| pod_cidr          | string | 否   | 此节点Pod地址的分配范围                |

### 请求参数示例

```json
 {
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "data": [
    {
      "bk_host_id": 1,
      "bk_cluster_id": 1,
      "name": "k8s",
      "roles": "master",
      "labels": {
        "env": "test"
      },
      "taints": {
        "type": "gpu"
      },
      "unschedulable": false,
      "internal_ip": [
        "127.0.0.1"
      ],
      "external_ip": [
        "127.0.0.1"
      ],
      "hostname": "xxx",
      "runtime_component": "runtime_component",
      "kube_proxy_mode": "ipvs",
      "pod_cidr": "127.0.0.1/26"
    },
    {
      "bk_host_id": 2,
      "bk_cluster_id": 1,
      "name": "k8s-node",
      "roles": "master",
      "labels": {
        "env": "test"
      },
      "taints": {
        "type": "gpu"
      },
      "unschedulable": false,
      "internal_ip": [
        "127.0.0.1"
      ],
      "external_ip": [
        "127.0.0.2"
      ],
      "hostname": "xxx",
      "runtime_component": "runtime_component",
      "kube_proxy_mode": "ipvs",
      "pod_cidr": "127.0.0.1/26"
    }
  ]
}
```

**注意：**

- internal_ip 和external_ip不能同时为空。
- 一次性创建节点数量不超过100个。

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

- 返回的data中的节点ID数组顺序与参数中的数组数据顺序保持一致。

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

| 名称  | 类型    | 描述          |
|-----|-------|-------------|
| ids | array | 创建的容器节点ID列表 |

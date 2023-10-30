### 功能描述

更新容器节点属性字段(v3.12.1+，权限: 容器节点编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型     | 必选  | 描述             |
|-----------|--------|-----|----------------|
| bk_biz_id | int    | 是   | 所属业务ID         |
| ids       | array  | 是   | 需要更新的node id列表 |
| data      | object | 是   | 需要更新的节点属性字段    |

#### data

| 字段                | 类型          | 必选  | 描述                               |
|-------------------|-------------|-----|----------------------------------|
| labels            | json object | 否   | 标签                               |
| taints            | string      | 否   | cluster 名称                       |
| unschedulable     | bool        | 否   | 设置是否可调度                          |
| hostname          | string      | 否   | 主机名                              |
| runtime_component | string      | 否   | 运行时组件                            |
| kube_proxy_mode   | string      | 否   | Kube-proxy 代理模式                  |
| pod_cidr          | string      | 否   | 此节点Pod地址的分配范围，例如：172.17.0.128/26 |

**注意：**

- 其中labels、taints是需要整体更新的。
- 一次性更新集群数量不超过100个。

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "ids": [
    1,
    2
  ],
  "data": {
    "labels": {
      "env": "test"
    },
    "taints": {
      "type": "gpu"
    },
    "unschedulable": false,
    "hostname": "xxx",
    "runtime_component": "runtime_component",
    "kube_proxy_mode": "ipvs",
    "pod_cidr": "127.0.0.1/26"
  }
}
```

### 返回结果示例

```json
 {
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": null
}
```

### 返回结果参数说明

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 无数据返回                      |

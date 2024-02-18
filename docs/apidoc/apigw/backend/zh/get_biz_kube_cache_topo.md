### 描述

根据业务ID,查询该业务的全量容器拓扑树信息 (版本：v3.12.5+，权限：业务访问)
该业务容器拓扑树的全量信息，包含了从业务这个根节点开始，到Cluster、Namespace、Workload层级的所有拓扑层级树数据。

注意：

- 该接口为缓存接口，默认全量缓存刷新时间为15分钟。
- 如果业务的容器拓扑信息发生变化，会通过事件机制实时刷新该业务的容器拓扑数据缓存。

### 输入参数

| 参数名称      | 参数类型 | 必选 | 描述                 |
|-----------|------|----|--------------------|
| bk_biz_id | int  | 是  | 要查询的业务容器拓扑所属的业务的ID |

### 调用示例

```json
{
  "bk_biz_id": 3
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "biz": {
      "id": 3,
      "nm": "biz1",
      "cnt": 100
    },
    "nds": [
      {
        "kind": "cluster",
        "id": 22,
        "nm": "cluster1",
        "cnt": 100,
        "nds": [
          {
            "kind": "namespace",
            "id": 16,
            "nm": "namespace1",
            "cnt": 100,
            "nds": [
              {
                "kind": "deployment",
                "id": 48,
                "nm": "deployment1",
                "cnt": 11,
                "nds": null
              },
              {
                "kind": "daemonSet",
                "id": 49,
                "nm": "daemonSet1",
                "cnt": 89,
                "nds": null
              }
            ]
          }
        ]
      }
    ]
  }
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称 | 参数类型         | 描述         |
|------|--------------|------------|
| biz  | object       | 业务信息       |
| nds  | object array | 业务下的拓扑节点数据 |

#### data.biz

| 参数名称 | 参数类型   | 描述              |
|------|--------|-----------------|
| id   | int    | 业务ID            |
| nm   | string | 业务名             |
| cnt  | int    | 业务下的Container数量 |

#### data.nds[n]

| 参数名称 | 参数类型   | 描述                                                                                                                   |
|------|--------|----------------------------------------------------------------------------------------------------------------------|
| kind | string | 该节点的资源类型，目前支持的类型：cluster、namespace、deployment、daemonSet、statefulSet、gameStatefulSet、gameDeployment、cronJob、job、pods。 |
| id   | int    | 该节点的ID                                                                                                               |
| nm   | string | 该节点的名称                                                                                                               |
| cnt  | int    | 该节点下的Container数量                                                                                                     |
| nds  | object | 该节点下的子节点信息，按照拓扑层级逐级循环嵌套。该节点下若无子节点，则该值为空。                                                                             |
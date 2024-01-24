### 描述

新建容器集群(v3.12.1+，权限:容器集群的创建权限)

### 输入参数

| 参数名称              | 参数类型   | 必选 | 描述                                                     |
|-------------------|--------|----|--------------------------------------------------------|
| bk_biz_id         | int    | 是  | 业务ID                                                   |
| name              | string | 是  | 集群名称                                                   |
| scheduling_engine | string | 否  | 调度引擎                                                   |
| uid               | string | 是  | 集群自有ID                                                 |
| xid               | string | 否  | 关联集群ID                                                 |
| version           | string | 否  | 集群版本                                                   |
| network_type      | string | 否  | 网络类型                                                   |
| region            | string | 否  | 地域                                                     |
| vpc               | string | 否  | vpc网络                                                  |
| network           | array  | 否  | 集群网络                                                   |
| type              | string | 是  | 集群类型。枚举值：INDEPENDENT_CLUSTER（独立集群）、SHARE_CLUSTER（共享集群） |
| environment       | string | 否  | 环境                                                     |
| bk_project_id     | string | 否  | 项目ID                                                   |
| bk_project_name   | string | 否  | 项目名称                                                   |
| bk_project_code   | string | 否  | 项目英文名                                                  |

### 调用示例

```json
{
  "bk_biz_id": 2,
  "name": "cluster",
  "scheduling_engine": "k8s",
  "uid": "xxx",
  "xid": "xxx",
  "version": "1.1.0",
  "network_type": "underlay",
  "region": "xxx",
  "vpc": "xxx",
  "network": [
    "127.0.0.0/21"
  ],
  "type": "INDEPENDENT_CLUSTER",
  "environment": "xxx",
  "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
  "bk_project_name": "test",
  "bk_project_code": "test"
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
    "id": 1
  },
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

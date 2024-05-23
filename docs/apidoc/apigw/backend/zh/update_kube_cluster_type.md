### 描述

更新容器集群类型(v3.12.1+，权限:容器集群的编辑权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                                                     |
|-----------|--------|----|--------------------------------------------------------|
| bk_biz_id | int    | 是  | 业务ID                                                   |
| id        | int    | 是  | cluster在cmdb中的唯一ID列表                                   |
| type      | string | 是  | 集群类型。枚举值：INDEPENDENT_CLUSTER（独立集群）、SHARE_CLUSTER（共享集群） |

### 调用示例

```json
{
  "bk_biz_id": 2,
  "id": 1,
  "type": "INDEPENDENT_CLUSTER"
}
```

### 响应示例

```json
 {
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 无数据返回                      |

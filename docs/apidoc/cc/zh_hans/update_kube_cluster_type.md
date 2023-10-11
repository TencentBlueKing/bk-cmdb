### 功能描述

更新容器集群类型(v3.12.1+，权限:容器集群的编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型     | 必选  | 描述                                                     |
|-----------|--------|-----|--------------------------------------------------------|
| bk_biz_id | int    | 是   | 业务ID                                                   |
| id        | int    | 是   | cluster在cmdb中的唯一ID列表                                   |
| type      | string | 是   | 集群类型。枚举值：INDEPENDENT_CLUSTER（独立集群）、SHARE_CLUSTER（共享集群） |

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "id": 1,
  "type": "INDEPENDENT_CLUSTER"
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

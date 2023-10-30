### 功能描述

删除容器集群(v3.12.1+， 权限:容器集群的删除权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型    | 必选  | 描述              |
|-----------|-------|-----|-----------------|
| bk_biz_id | int   | 是   | 容器集群所属业务ID      |
| ids       | array | 是   | 容器集群在cmdb中的ID列表 |

**注意：**

- 用户需要保证所要删除集群下没有关联资源(如namespace、pod、node workload等)，否则会删除失败。
- 一次性删除集群的数量不能超过10个。

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
  ]
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 返回结果参数说明

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 无数据返回                      |
| request_id | string | 请求链id                      |

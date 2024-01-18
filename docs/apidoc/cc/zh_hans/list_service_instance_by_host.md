### 功能描述

根据主机id获取绑定到主机上的服务实例列表

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   | 描述                   |
|----------------------|------------|--------|----------------------|
| bk_biz_id            | int  | 是   | 业务id                 |
| bk_host_id            | int  | 是   | 主机ID,获取绑定到主机上的服务实例信息 |
| page       |  object    | 否     | 查询条件                 |

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大500 |

### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 1
  },
  "bk_host_id": 26
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "count": 1,
    "info": [
       {
          "bk_biz_id": 1,
          "id": 1,
          "name": "test",
          "labels": {
              "test1": "1"
          },
          "service_template_id": 32,
          "bk_host_id": 26,
          "bk_module_id": 12,
          "creator": "admin",
          "modifier": "admin",
          "create_time": "2021-12-31T03:11:54.992Z",
          "last_time": "2021-12-31T03:11:54.992Z",
          "bk_supplier_account": "0"
      }
    ]
  }
}
```

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | object | 请求返回的数据 |

#### data 字段说明

| 字段|类型|描述|
|---|---|---|
|count|int|总数|
|info|array|返回结果|

#### info 字段说明

| 字段|类型|说明|
|---|---|---|
|id|int|服务实例ID|
|name|string|服务实例名称|
|bk_module_id|int|模型id|
|service_template_id|int|服务模版ID|
| labels           | map  |标签信息 |
|bk_host_id|int|主机id|
| creator              | string             | 本条数据创建者                                                                                 |
| modifier             | string             | 本条数据的最后修改人员            |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string       | 开发商账号 |

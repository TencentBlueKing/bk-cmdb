### 功能描述

创建集群

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | 否     | 开发商账号 |
| bk_biz_id      | int     | 是     | 业务ID |
| data           | dict    | 是     | 集群数据 |

#### data

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_parent_id        |  int     | 是     | 父节点的ID |
| bk_set_name         |  string  | 是     | 集群名字 |
| default             |  int     | 否     | 0-普通集群，1-内置模块集合，默认为0 |
| set_template_id     |  int     | 否     | 集群模板ID，需要通过集群模板创建集群时必填 |
|bk_capacity        |   int      |  否   |   设计容量    |
| description           | string     | 否     | 数据的描述信息     |
|bk_set_desc|string|否|集群描述|

**注意：其它需要的字段由模型定义**

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 1,
    "data": {
        "bk_parent_id": 1,
        "bk_set_name": "test-set",
        "default": 0,
        "bk_set_desc": "test-set",
        "bk_capacity": 1000,
        "description": "description",
        "set_template_id": 1
    }
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
     "data": {
        "bk_biz_id": 11,
        "bk_capacity": 1000,
        "bk_parent_id": 11,
        "bk_service_status": "1",
        "bk_set_desc": "test-set",
        "bk_set_env": "3",
        "bk_set_id": 4780,
        "bk_set_name": "test-set",
        "bk_supplier_account": "0",
        "create_time": "2022-02-22T20:34:01.386+08:00",
        "default": 0,
        "description": "description",
        "last_time": "2022-02-22T20:34:01.386+08:00",
        "set_template_id": 11,
        "set_template_version": null
     }
}
```
### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| data    | object | 请求返回的数据                           |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |

#### data
| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
| bk_biz_id | int | 业务id |
| bk_capacity | int | 设计容量 |
|bk_parent_id|int|父节点的ID|
| bk_set_id | int | 集群id |
| bk_service_status | string   | 服务状态:1/2(1:开放,2:关闭)           |
|bk_set_desc|string|集群描述|
| bk_set_env        | string   | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
|bk_set_name|string|集群名称|
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
| default             |  int     | 0-普通集群，1-内置模块集合，默认为0 |
| description           | string     | 数据的描述信息     |
| set_template_version|  array |集群模板的当前版本 |
| set_template_id|  int |集群模板ID |
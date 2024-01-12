### 功能描述

查询集群

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | 否     | 开发商账号 |
| bk_biz_id      |  int     | 是     | 业务id |
| fields         |  array   | 是     | 查询字段，所有字段均为set定义的字段，这些字段包括预置字段，也包括用户自定义字段 |
| condition      |  dict    | 是     | 查询条件，所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段 |
| page           |  dict    | 是     | 分页条件 |

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大200 |
| sort     |  string | 否     | 排序字段 |

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 2,
    "fields": [
        "bk_set_name"
    ],
    "condition": {
        "bk_set_name": "test"
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_set_name"
    }
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
                "bk_set_name": "test",
                "default": 0
            }
        ]
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

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int       | 数据数量 |
| info      | array     | 结果集，其中，所有字段均为集群定义的属性字段 |

#### info
| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| bk_set_name     | string       | 集群名称 |
| default             |  int     | 0-普通集群，1-内置模块集合，默认为0 |
| bk_biz_id | int | 业务id |
| bk_capacity | int | 设计容量 |
|bk_parent_id|int|父节点的ID|
| bk_set_id | int | 集群id |
| bk_service_status | string   | 服务状态:1/2(1:开放,2:关闭)           |
|bk_set_desc|string|集群描述|
| bk_set_env        | string   | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
| description           | string     | 数据的描述信息     |
| set_template_version|  array |集群模板的当前版本 |
| set_template_id|  int |集群模板ID |
| bk_created_at      | string |  创建时间        |
| bk_updated_at      | string |  更新时间        |
| bk_created_by      | string |  创建人         |
**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
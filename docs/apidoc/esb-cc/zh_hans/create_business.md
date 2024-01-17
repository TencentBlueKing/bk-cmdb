### 功能描述

新建业务(权限：业务创建权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | 否     | 开发商账号 |
| data           | dict    | 是     | 业务数据 |

#### data

| 字段      |  类型      | 必选 | 描述                            |
|-----------|------------|----|-------------------------------|
| bk_biz_name       |  string  | 是  | 业务名                           |
| bk_biz_maintainer |  string  | 是  | 运维人员                          |
| bk_biz_productor  |  string  | 否  | 产品人员                          |
| bk_biz_developer  |  string  | 否  | 开发人员                          |
| bk_biz_tester     |  string  | 否  | 测试人员                          |
| operator     |  string  | 否  | 操作人员                          |
| bk_created_at      | string | 否     | 创建时间                          |
| bk_updated_at      | string | 否     | 更新时间                          |
| bk_created_by      | string | 否     | 创建人                           |
| bk_updated_by      | string | 否     | 更新人                           |
| life_cycle      | string | 否     | 生命周期：测试中(1)，已上线(2, 默认值)，停运(3) |
| time_zone         |  string  | 否  | 时区                            |
| language          |  string  | 否  | 语言, "1"代表中文, "2"代表英文          |

**注意：此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段**

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "data": {
        "bk_biz_name": "main-a1",
        "bk_biz_maintainer": "admin",
        "bk_biz_productor": "admin",
        "bk_biz_developer": "admin",
        "bk_biz_tester": "admin",
        "time_zone": "Asia/Shanghai",
        "bk_created_at": "",
        "bk_updated_at": "",
        "bk_created_by": "admin",
        "bk_updated_by": "admin",
        "language": "1",
        "operator": "admin",
        "life_cycle": "2"
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
        "bk_biz_developer": "admin",
        "bk_biz_id": 5,
        "bk_biz_maintainer": "admin",
        "bk_biz_name": "main-a1",
        "bk_biz_productor": "admin",
        "bk_biz_tester": "admin",
        "bk_created_at": "2023-11-14T16:51:02.168+08:00",
        "bk_created_by": "admin",
        "bk_supplier_account": "0",
        "bk_updated_at": "2023-11-14T16:51:02.168+08:00",
        "create_time": "2023-11-14T16:51:02.168+08:00",
        "default": 0,
        "language": "1",
        "last_time": "2023-11-14T16:51:02.168+08:00",
        "life_cycle": "2",
        "operator": "admin",
        "time_zone": "Asia/Shanghai"
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
| bk_biz_name       |  string       | 业务名 |
| bk_biz_maintainer |  string       | 运维人员 |
| bk_biz_productor  |  string      | 产品人员 |
| bk_biz_developer  |  string      | 开发人员 |
| bk_biz_tester     |  string       | 测试人员 |
| time_zone         |  string       | 时区 |
| language          |  string      | 语言, "1"代表中文, "2"代表英文 |
| bk_supplier_account | string       | 开发商账号   |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
|default | int | 表示业务类型 |
| operator | string | 主要维护人 |
|life_cycle|string|业务生命周期|
| bk_created_at      | string |  创建时间        |
| bk_updated_at      | string |  更新时间        |
| bk_created_by      | string |  创建人         |
### 功能描述

订阅事件

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型      | 必选   |  描述                                            |
|---------------------|------------|--------|--------------------------------------------------|
| subscription_name   | string     | 是     | 订阅的名字                                       |
| system_name         | string     | 是     | 订阅事件的系统的名字                             |
| callback_url        | string     | 是     | 回调函数                                         |
| confirm_mode        | string     | 是     | 事件发送成功校验模式,可选 1-httpstatus,2-regular |
| confirm_pattern     | string     | 是     | callback的httpstatus或正则                       |
| subscription_form   | string     | 是     | 订阅的事件,以逗号分隔                            |
| timeout             | int        | 是     | 发送事件超时时间                                 |

### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "subscription_name":"mysubscribe",
  "system_name":"SystemName",
  "callback_url":"http://127.0.0.1:8080/callback",
  "confirm_mode":"httpstatus",
  "confirm_pattern":"200",
  "subscription_form":"hostcreate",
  "timeout":10
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
    "data":{
        "subscription_id": 1
    }
}
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

#### data

| 字段            | 类型    | 描述             |
|-----------------|---------|------------------|
| subscription_id | int     | 新增订阅的订阅ID |

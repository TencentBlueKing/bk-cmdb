### 功能描述

修改订阅

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                   |  类型    | 必选   |  描述                                            |
|------------------------|----------|--------|--------------------------------------------------|
| bk_supplier_account    | string   | 是     | 开发商账号                                       |
| subscription_id        | int      | 是     | 订阅ID                                           |
| subscription_name      | string   | 是     | 订阅的名字                                       |
| system_name            | string   | 是     | 订阅事件的系统的名字                             |
| callback_url           | string   | 是     | 回调函数                                         |
| confirm_mode           | string   | 是     | 事件发送成功校验模式,可选 1-httpstatus,2-regular |
| confirm_pattern        | string   | 是     | callback的httpstatus或正则                       |
| subscription_form      | string   | 是     | 订阅的事件,以逗号分隔                            |
| timeout                | int      | 是     | 发送事件超时时间                                 |


### 请求参数示例

```python
{
  "bk_supplier_account": "0",
  "subscription_name":"mysubscribe",
  "subscription_id": 2,
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
    "data": "success"
}
```

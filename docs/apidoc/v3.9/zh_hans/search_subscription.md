### 功能描述

查询事件订阅

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型      | 必选   |  描述                       |
|---------------------|------------|--------|-----------------------------|
| page                | object     | 否     | 分页参数                    |
| condition           | object     | 否     | 查询条件                    |

#### page

| 字段      |  类型      | 必选   |  描述                |
|-----------|------------|--------|----------------------|
| start     |  int       | 是     | 记录开始位置         |
| limit     |  int       | 是     | 每页限制条数,最大200 |
| sort      |  string    | 否     | 排序字段             |

#### condition

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| subscription_name  |string      |是      | 此处仅为示例数据，需要被设置为查询的字段 |

### 请求参数示例

```python
{
    "bk_supplier_account":"0",
    "condition":{
        "subscription_name":"name"
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"HostName"
    }
}
```

### 返回结果示例

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "count": 1,
        "info": [
            {
                "subscription_id": 1,
                "subscription_name": "mysubscribe",
                "system_name": "SystemName",
                "callback_url": "http://127.0.0.1:8080/callback",
                "confirm_mode": "httpstatus",
                "confirm_pattern": "200",
                "time_out": 10,
                "subscription_form": "hostcreate",
                "operator": "user",
                "bk_supplier_account": "0",
                "last_time": "2017-09-19 16:57:07",
                "statistics": {
                    "total": 30,
                    "failure": 2
                }
            }
        ]
    }
}
```

### 返回结果参数说明

#### data
| 字段   | 类型         | 描述              |
|-------|--------------|------------------|
| count | int          | 记录条数          |
| info  | object array | 事件订阅的详情列表  |

#### info
| 字段                 | 类型      | 描述                                       |
|----------------------|-----------|--------------------------------------------|
| subscription_id      | int       | 订阅ID                                     |
| subscription_name    | string    | 订阅名                                     |
| system_name          | string    | 系统名称                                   |
| callback_url         | string    | 回调地址                                   |
| confirm_mode         | string    | 回调成功确认模式，可选:httpstatus，regular |
| confirm_pattern      | string    | 回调成功标志                               |
| subscription_form    | string    | 订阅单，用","分隔                          |
| timeout              | int       | 超时时间，单位：秒                         |
| operator             | int       | 本条数据的最后更新人员                     |
| last_time            | int       | 更新时间                                   |
| statistics.total     | int       | 推送总数                                   |
| statistics.failure   | int       | 推送失败数                                 |
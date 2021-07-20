#### 根据条件获取操作日志

* API:  POST /api/v3/audit/search
* API 名称：get_operation_log
* 功能说明：
	- 中文： 获取操作日志
	- English：get operation log
* input:
```
{
    "condition":{
        "audit_type": "business",
        "bk_supplier_account": "0",
        "user": "admin",
        "resource_type": "business",
        "action": "create",
        "operation_time": [
            "2017-12-25 10:10:10",
            "2017-12-25 10:10:11"
        ]
        "operate_from": "user",
        "operation_detail": {
            "bk_biz_id":99999
        },
    },
    "start":0,
    "limit":10,
    "sort":"-create_time"
}
```

* input 参数说明

| 名称                | 类型       | 必填 | 默认值 | 说明                                   | Description                                      |
| ------------------- | ---------- | ---- | ------ | -------------------------------------- | ------------------------------------------------ |
| audit_type          | string     | 否   | 无     | 审计资源大类                           | audit type                                       |
| bk_supplier_account | string     | 否   | "0"    | 供应商账号                             | supplier account                                 |
| user                | string     | 否   | 无     | 创建审计的用户                         | the user who created the audit log               |
| resource_type       | string     | 否   | 无     | 审计资源类型                           | audit resource type                              |
| action              | string     | 否   | 无     | 操作类型                               | audit action, it can be add, delete, update etc. |
| operation_time      | string数组 | 否   | 无     | 没有条件，为空, 开始和结束时间成对出现 | no condition, start time and end time is pair    |
| operate_from        | string     | 否   | 无     | 审计来源                               | where audit operation comes from                 |
| operation_detail    | object     | 否   | 无     | 审计详情，包括业务ID等                 | audit operation detail, contains biz id etc.     |
| start               | int        | 是   | 无     | 记录开始位置                           | start record                                     |
| limit               | int        | 是   | 无     | 每页限制条数,最大200                   | page limit, max is 200                           |
| sort                | string     | 否   | 无     | 排序字段                               | the field for sort                               |

ext_key 字段说明： 为根据ip的匹配搜索


* output

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
        "count":1,
        "info":[
            {
                "audit_type": "business",
                "bk_supplier_account": "0",
                "user": "admin",
                "resource_type": "business",
                "action": "update",
                "operation_time": [
                    "2017-12-25 10:10:10",
                    "2017-12-25 10:10:11"
                ]
                "operate_from": "user",
                "operation_detail": {
                    "bk_biz_id":100,
                    "resource_id": 100,
                    "resource_name": "test",
                    "details": {
                        "pre_data":{
                            "last_time":"2018-03-08T15:10:42.264+08:00",
                            "bk_biz_name": "test1",
                            "create_time":"2018-03-08T14:23:05.05+08:00",
                            "language" : "1",
                            "bk_supplier_account" : "0",
                            "bk_supplier_id" : 0,
                            "bk_biz_tester" : "",
                            "operator" : "",
                            "bk_biz_maintainer" : "admin",
                            "time_zone" : "Asia/Shanghai",
                            "life_cycle" : "2",
                            "default" : 0,
                            "bk_biz_productor" : "",
                            "bk_biz_developer" : "",
                            "bk_biz_id" : 100
                        },
                        "cur_data":{
                             "last_time":"2018-03-08T15:10:42.264+08:00",
                             "bk_biz_name": "test",
                             "create_time":"2018-03-08T14:23:05.05+08:00",
                             "language" : "1",
                             "bk_supplier_account" : "0",
                             "bk_supplier_id" : 0,
                             "bk_biz_tester" : "",
                             "operator" : "",
                             "bk_biz_maintainer" : "admin",
                             "time_zone" : "Asia/Shanghai",
                             "life_cycle" : "2",
                             "default" : 0,
                             "bk_biz_productor" : "",
                             "bk_biz_developer" : "",
                             "bk_biz_id" : 100
                        },
                        "properties":[
                            {
                                "bk_property_id": "bk_biz_name",
                                "bk_property_name": "业务名"
                            },
                            {
                                "bk_property_id": "life_cycle",
                                "bk_property_name": "生命周期"
                            },
                            {
                                "bk_property_id": "bk_biz_maintainer",
                                "bk_property_name": "运维人员"
                            },
                            {
                                "bk_property_id": "bk_biz_productor",
                                "bk_property_name": "产品人员"
                            },
                            {
                                "bk_property_id": "bk_biz_tester",
                                "bk_property_name": "测试人员"
                            },
                            {
                                "bk_property_id": "bk_biz_developer",
                                "bk_property_name": "开发人员"
                            },
                            {
                                "bk_property_id": "operator",
                                "bk_property_name": "操作人员"
                            },
                            {
                                "bk_property_id": "time_zone",
                                "bk_property_name": "时区"
                            },
                            {
                                "bk_property_id": "language",
                                "bk_property_name": "语言"
                            }
                        ],
                    }
                }
            }
        ]
    }
}
```

* output字段说明：

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | the data response                                          |

data 字段说明：

| 名称  | 类型         | 说明               | Description               |
| ----- | ------------ | ------------------ | ------------------------- |
| count | int          | 请求记录条数       | the count of record       |
| info  | object array | record information | the information of record |

info 字段说明：

| 名称                | 类型   | 说明                   | Description                                      |
| ------------------- | ------ | ---------------------- | ------------------------------------------------ |
| audit_type          | string | 审计资源大类           | audit type                                       |
| bk_supplier_account | string | 供应商账号             | supplier account                                 |
| user                | string | 创建审计的用户         | the user who created the audit log               |
| resource_type       | string | 审计资源类型           | audit resource type                              |
| action              | string | 操作类型               | audit action, it can be add, delete, update etc. |
| operation_time      | time   | 审计记录时间           | audit log create time                            |
| operate_from        | string | 审计来源               | where audit operation comes from                 |
| operation_detail    | object | 审计详情，包括业务ID等 | audit operation detail, contains biz id etc.     |

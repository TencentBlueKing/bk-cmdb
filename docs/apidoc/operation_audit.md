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
        "bk_biz_id":99999,
        "ext_key":{
            "$in":[
                "127.0.0.23",
                "127.0.0.22"
            ]
        },
        "op_target":"host",
        "op_type":"add",
        "op_time":[
            "2017-12-25 10:10:10",
            "2017-12-25 10:10:11"
        ]
    },
    "start":0,
    "limit":10,
    "sort":"-create_time"
}
```

* input 参数说明

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务ID |  business ID|
|ext_key|object|否|无|当bk_op_type为0， 填入多个IP地址| if bk_op_type=0, the input is ip array|
|op_target|string|否|无|操作对象，可以为biz host process set module object| op target, and it can be biz host process set module object|
|op_type|string|否|无|操作类型， add delete update | op type, and it can be add , delete ,update|
|op_time|string数组|否|无|没有条件，为空, 开始和结束时间成对出现 | no condition, start time and end time is pair|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

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
                "bk_supplier_account":"0",
                "bk_biz_id":1,
                "op_desc":"修改主机",
                "op_type":2,
                "op_target":"host",
                "operator":"admin",
                "content":{
                    "pre_data":{
                        "last_time":"2018-03-08T15:10:42.264+08:00",
                        "bk_cloud_id":[
                            {
                                "ref_id":1,
                                "ref_name":"Direct connecting area"
                            }
                        ],
                        "create_time":"2018-03-08T14:23:05.05+08:00",
                        "bk_host_id":1,
                        "bk_host_innerip":"127.0.01",
                        "bk_import_from":"1"
                    },
                    "cur_data":{
                        "last_time":"2018-03-08T15:10:42.264+08:00",
                        "bk_cloud_id":[
                            {
                                "ref_id":2,
                                "ref_name":"test connecting area"
                            }
                        ],
                        "create_time":"2018-03-08T14:23:05.05+08:00",
                        "bk_host_id":1,
                        "bk_host_innerip":"127.0.0.1",
                        "bk_import_from":"1"
                    },
                    "header":[
                        {
                            "bk_property_id":"bk_host_innerip",
                            "bk_property_name":"内网IP"
                        },
                        {
                            "bk_property_id":"bk_host_outerip",
                            "bk_property_name":"外网IP"
                        }
                    ],
                    "type":"map"
                },
                "ext_key":"127.0.0.1",
                "op_time":"2018-03-08T03:30:28.056Z",
                "inst_id":1
            }
        ]
    }
}
```

* output字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object | 请求返回的数据 |the data response|

data 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| count| int| 请求记录条数 |the count of record|
| info| object array | record information | the information of record  |

info 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| bk_supplier_account| string| 开发商ID |supplier account code|
| bk_biz_id| int | 业务ID | business ID  |
| op_type| string | 操作类型 | the type of record  |
| op_desc| string | 操作描述 | operation description  |
| op_target| string| 操作对象 | operation target  |
| operator| string | 操作者 | the man operate it  |
| content| object 对象 | 操作内容 | operate content  |
| ext_key| string  | 附加信息 | ext key  |
| op_time| string |  操作时间 | operation time  |
| inst_id| int | 实例ID | instantiation ID |

content  字段说明： content为实际的操作内容
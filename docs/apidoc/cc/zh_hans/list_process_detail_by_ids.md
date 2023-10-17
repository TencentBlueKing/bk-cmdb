### 功能描述

查询某业务下进程ID对应的进程详情 (v3.9.8)

### 请求参数

{{ common_args_desc }}

#### 接口参数

|字段|类型|必填|描述|
|---|---|---|---|
|bk_biz_id|int|是| 进程所在的业务ID |
|bk_process_ids|array|是|进程ID列表，最多传500个|
|fields|array|否|进程属性列表，控制返回结果的进程实例信息里有哪些字段，能够加速接口请求和减少网络流量传输<br>为空时返回进程所有字段,bk_process_id为必返回字段|


### 请求参数示例

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":1,
    "bk_process_ids": [
        43,
        44
    ],
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_id",
        "bk_func_name"
    ]
}
```

### 返回结果示例
``` json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_func_id": "",
            "bk_func_name": "pr1",
            "bk_process_id": 43,
            "bk_process_name": "pr1"
        },
        {
            "bk_func_id": "",
            "bk_func_name": "pr2",
            "bk_process_id": 44,
            "bk_process_name": "pr2"
        }
    ]
}
```

### 返回结果参数说明

| 名称  | 类型  | 描述 |
|---|---|--- |
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | array | 请求返回的数据 |

#### data
| 名称  | 类型  | 描述 |
|---|---|--- |
|bk_func_id|string|功能ID|
|bk_func_name|string|进程名称|
|bk_process_id|int|进程id|
|bk_process_name|string|进程别名|
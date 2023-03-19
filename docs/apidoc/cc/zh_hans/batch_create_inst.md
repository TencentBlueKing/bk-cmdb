### 功能描述

 批量创建通用模型实例(v3.10.2+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 参数      | 类型   | 必选 | 描述               |
| -------- | ------ | ---- | ------------------ |
| bk_obj_id | string | 是   | 用于创建的模型id，只允许创建通用模型的实例   |
| details   | array | 是   | 需要创建的实例内容，最多不能超过200个，内容为该模型实例的属性信息 |

#### details

| 参数            | 类型   | 必选 | 描述           |
| --------------- | ------ | ---- | -------------- |
| bk_inst_name      | string | 是   | 实例名   |
| bk_asset_id      | string | 是  | 固资编号      | 
| bk_sn | string |否 |设备SN  |
| bk_operator | string |否 |维护人  |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"bk_switch",
    "details":[
        {
            "bk_inst_name":"s1",
            "bk_asset_id":"test_001",
            "bk_sn":"00000001",
            "bk_operator":"admin"
        },
        {
            "bk_inst_name":"s2",
            "bk_asset_id":"test_002",
            "bk_sn":"00000002",
            "bk_operator":"admin"
        },
        {
            "bk_inst_name":"s3",
            "bk_asset_id":"test_003",
            "bk_sn":"00000003",
            "bk_operator":"admin"
        }
    ]
}
```

### 返回结果示例

```json
{
    "result":true,
    "code":0,
    "message":"",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data":{
        "success_created":{
            "1":1001,
            "2":1002
        },
        "error_msg":{
            "0":"数据唯一性校验失败， [bk_asset_id: test_001] 重复"
        }
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
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data

| 字段            | 类型 | 描述                                                     |
| -------------- | ---- | -------------------------------------------------------- |
| success_created | map | key为实例在参数details中的index，value为创建成功的实例id |
| error_msg       | map | key为实例在参数details中的index，value为失败信息          |
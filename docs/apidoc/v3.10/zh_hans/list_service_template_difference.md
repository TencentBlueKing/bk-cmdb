### 功能描述

列出服务模版和服务实例之间的差异 (v3.9.19)

- 该接口专供GSEKit使用，在ESB文档中为hidden状态

### 请求参数

{{ common_args_desc }}

#### 接口参数

|字段|类型|必填|描述|
|---|---|---|---|
| bk_biz_id  | int64       | Yes      | 业务ID |
|bk_module_ids|int64 array|No|模块ID列表，最多不能超过20个|
|service_template_ids|int64 array|No|服务模板ID列表，最多不能超过20个|
|is_partial|bool|Yes|为true时，使用service_template_ids参数，返回service_template的状态；为false时，使用bk_module_ids参数，返回module的状态|


### 请求参数示例

- 示例1
``` json
{
    "bk_biz_id": 3,
    "service_template_ids": [
        1,
        2
    ],
    "is_partial": true
}
```
- 示例2
```
{
    "bk_biz_id": 3,
    "bk_module_ids": [
        11,
        12
    ],
    "is_partial": false
}
```

### 返回结果示例
- 示例1
``` json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "service_templates": [
            {
                "service_template_id": 1,
                "need_sync": true
            },
            {
                "service_template_id": 2,
                "need_sync": false
            }
        ]
    }
}
```
- 示例2
```
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "modules": [
            {
                "bk_module_id": 11,
                "need_sync": false
            },
            {
                "bk_module_id": 12,
                "need_sync": true
            }
        ]
    }
}
```

### 返回结果参数说明

| 名称  | 类型  | 描述 |
|---|---|--- |
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |

- data 字段说明

| 名称  | 类型  | 描述 |
|---|---|--- |
|service_templates|object array|服务模板信息列表|
|modules|object array|模块信息列表|

- service_templates 字段说明

| 名称  | 类型  | 描述 |
|---|---|--- |
|service_template_id|int|服务模板ID|
|need_sync|bool|服务模版所应用的模块下的服务实例和服务模板是否有差异|

- modules 字段说明

| 名称  | 类型  | 描述 |
|---|---|--- |
|bk_module_id|int|模块ID|
|need_sync|bool|模块下的服务实例和服务模板是否有差异|

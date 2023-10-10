### 功能描述

获取推送主机身份到机器结果（只能获取到30分钟内推送的任务情况）
(版本：v3.10.23+，权限：当任务包含的主机在业务下时, 需要对应业务的访问权限；当主机在主机池下时，需要主机的更新权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段 | 类型 | 必选 | 描述       |
| ---- | ---- | ---- | ---------- |
|  task_id | string    |  是  |任务id |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "task_id": "GSETASK:F:202201251046313618521052:198"
}
```

### 返回结果示例
```json
{
    "result": true,
    "code": 0,
    "msg": "success",
    "permission": null,
    "request_id": "c11aasdadadadsadasdadasd1111ds",
    "data": {
            "success_list": [
                1,
                2
            ],
            "pending_list": [
                3,
                4
            ],
            "failed_list": [
                5,
                6
            ]
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

#### data 字段说明
| 名称    | 类型   | 描述                                       |
| ------- | ------ | ------------------------------------------ |
|  success_list |  array  |执行成功的主机id列表 |
|  failed_list |  array  |执行失败的主机id列表 |
|  pending_list |  array  |调用gse下发主机身份，还没有拿到结果的主机id列表 |
### 功能描述

推送主机身份信息到机器上，返回本次任务gse的任务id，可以根据此gse任务id去gse查询任务的推送结果(v3.10.18+，对于在业务中的主机，需要业务访问权限，对于在主机池的主机，需要主机池主机编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段 | 类型 | 必选 | 描述       |
| ---- | ---- | ---- | ---------- |
| bk_host_ids     |  array | 是    | 主机id数组，数量不能超过200 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_ids": [1,2]
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
        "task_id": "GSETASK:F:202206222053523618521052:393",
        "host_infos": [
            {
                "bk_host_id": 2,
                "identification": "0:127.0.0.1"
            }
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
|  task_id |  string  |任务id，此id为gse侧的task_id |
|  host_infos |  array  |任务中推送的主机信息，只包含成功推送的信息 |

#### host_infos 字段说明
| 名称    | 类型   | 描述                                       |
| ------- | ------ | ------------------------------------------ |
|  bk_host_id |  int  |主机id |
|  identification |  string  |推送的主机在任务中对应的标识 |
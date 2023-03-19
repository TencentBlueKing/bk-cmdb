### 功能描述

同步主机身份到主机上（注意，由于是调用gse的接口进行操作，为异步过程，可能会出现下发时间慢导致查询不到结果将主机设为下发失败主机）

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
##### 1、传入数据没有问题并且调用gse接口正常情况下：
全部成功时：
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
        "failed_list": [],
        "task_id": "GSETASK:F:202201251046313618521052:198"
    }
}
```
出现失败时：
```json
{
    "result": true,
    "code": 0,
    "msg": "success",
    "permission": null,
    "request_id": "c11aasdadadadsadasdadasd1111ds",
    "data": {
        "success_list": [
            1
        ],
        "failed_list": [
            2
        ],
        "task_id": "GSETASK:F:202201251046313618521052:198"
    }
}
```
##### 2、传入数据有问题或调用gse接口有问题情况下：
```json
{
    "result": false,
    "code": xxx,
    "msg": "xxx",
    "permission": null,
    "request_id": "c11aasdadadadsadasdadasd1111ds",
    "data": null
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
|  task_id |  string  |任务id |

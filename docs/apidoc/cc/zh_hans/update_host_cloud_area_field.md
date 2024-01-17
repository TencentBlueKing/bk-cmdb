### 功能描述

根据主机id列表和管控区域id,更新主机的管控区域字段(权限：业务主机编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            | int  | 否   | 业务ID |
| bk_cloud_id         | int  | 是   | 管控区域ID |
| bk_host_ids         | array  | 是   | 主机IDs, 最多2000个 |


### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_ids": [43, 44], 
    "bk_cloud_id": 27,
    "bk_biz_id": 1
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
}
```

### 返回结果实例 - 管控区域 + 内网IP 重复

```python
{
  "result": false,
  "code": 1199014,
  "message": "数据唯一性校验失败， bk_host_innerip 重复",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": null
}
```

### 返回结果实例 - 一次操作主机数太多
```python
{
  "result": false,
  "code": 1199077,
  "message": "一次操作记录数超过最大限制：2000",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": null
}
```

### 返回结果参数说明

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                          |

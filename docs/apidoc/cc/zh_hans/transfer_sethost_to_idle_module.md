### 功能描述

根据业务id,集群id,模块id,将指定业务集群模块下的主机上交到业务的空闲机模块

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段          |  类型      | 必选     |  描述    |
|---------------|------------|----------|----------|
| bk_biz_id     | int        | 是       | 业务id   |
| bk_set_id     | int        | 是       | 集群id   |
| bk_module_id  | int        | 是       | 模块id   |


### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":10,
    "bk_module_id":58,
    "bk_set_id":1
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "sucess"
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

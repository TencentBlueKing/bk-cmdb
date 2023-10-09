### 功能描述

删除对象模型属性，可以删除业务自定义字段

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段  |  类型       | 必选   |  描述                         |
|-------|-------------|--------|-------------------------------|
| id    | int         | 否     | 被删除的数据记录的唯一标识ID  |


### 请求参数示例

```python

{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id" : 0
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
    "data": "success"
}
```

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

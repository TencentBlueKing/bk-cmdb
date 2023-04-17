### 功能描述

通过模型分类ID删除模型分类

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段  |  类型       | 必选   |  描述                            |
|-------|-------------|--------|----------------------------------|
| delete      | object | 是    |  删除  |

#### delete
| 字段                |  类型       | 必选   |  描述                            |
|---------------------|-------------|--------|----------------------------------|
|id     | int         | 是     | 分类数据记录ID                   |


### 请求参数示例

```python

{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "delete":{
    "id" : 0
    }
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

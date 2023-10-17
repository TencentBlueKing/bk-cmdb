### 功能描述

通过模型分类ID删除模型分类，只能删除空分类，当该分类下存在模型时无法删除

### 请求参数

{{ common_args_desc }}

#### 接口参数

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
    "id": 10
}
```


### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "3ceb0b81e6f04e52925dc6d9f419f058"
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

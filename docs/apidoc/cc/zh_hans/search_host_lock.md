### 功能描述

根据主机id列表查询主机锁(v3.8.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型       | 必选   |  描述                            |
|---------------------|-------------|--------|----------------------------------|
|id_list| array| 是|主机ID列表|


### 请求参数示例

```python
{
   "bk_app_code": "esb_test",
   "bk_app_secret": "xxx",
   "bk_username": "xxx",
   "bk_token": "xxx",
   "id_list":[1, 2]
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        1: true,
        2: false
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

#### data
| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
| data |object| 请求返回的数据, key 是 ID，value 是否上锁 |

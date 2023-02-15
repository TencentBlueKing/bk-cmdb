### 功能描述

查询模型的实例之间的关联关系。

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 是否必填	   |  描述 |
|----------------------|------------|--------|-----------------------------|
| condition | string map     | Yes   | 查询条件 |


condition params

| 字段                 |  类型      | 是否必填	   |  描述 |
|---------------------|------------|--------|-----------------------------|
| bk_asst_id           | string     | Yes     | 模型的关联类型唯一id|
| bk_obj_id           | string     | Yes     | 源模型id|
| bk_asst_id           | string     | Yes     | 目标模型id|


### 请求参数示例

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_asst_id": "belong",
        "bk_obj_id": "bk_switch",
        "bk_asst_obj_id": "bk_host"
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
           "id": 27,
           "bk_supplier_account": "0",
           "bk_obj_asst_id": "test1_belong_biz",
           "bk_obj_asst_name": "1",
           "bk_obj_id": "test1",
           "bk_asst_obj_id": "biz",
           "bk_asst_id": "belong",
           "mapping": "n:n",
           "on_delete": "none",
           "ispre": null
        }
    ]
}

```


### 返回结果参数说明
#### response
| 名称    | 类型   | 说明                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

#### data

| 字段       | 类型     | 描述 |
|------------|----------|--------------|
| id|int64|模型关联关系的身份id|
| bk_obj_asst_id| string|  模型关联关系的唯一id.|
| bk_obj_asst_name| string| 关联关系的别名. |
| bk_asst_id| string| 关联类型id|
| bk_obj_id| string| 源模型id |
| bk_asst_obj_id| string| 目标模型id|
| mapping| string|  源模型与目标模型关联关系实例的映身关系，可以是以下中的一种[1:1, 1:n, n:n] |
| on_delete| string| 删除关联关系时的动作, 取值为以下其中的一种 [none, delete_src, delete_dest], "none" 什么也不做, "delete_src" 删除源模型的实例, "delete_dest" 删除目标模型的实例.|
| bk_supplier_account | string | 开发商账号   |
| ispre               | bool         | true:预置字段,false:非内置字段                             |
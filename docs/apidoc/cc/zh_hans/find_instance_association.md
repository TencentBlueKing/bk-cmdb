### 功能描述

查询模型的实例关联关系。

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 是否必填	   |  描述          |
|----------------------|------------|--------|-----------------------------|
| condition | object     | Yes   | 查询条件 |
| bk_obj_id           | string     | YES     | 源模型id(v3.10+)|


#### condition

| 字段                 |  类型      | 是否必填	   |  描述         |
|---------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           | string     | Yes     | 模型关联关系的唯一id|
| bk_asst_id           | string     | NO     | 关联类型的唯一id|
| bk_asst_obj_id           | string     | NO     | 目标模型id|


### 请求参数示例

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_asst_obj_id": ""
    },
    "bk_obj_id": "xxx"
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
    "data": [{
        "id": 481,
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
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

#### data

| 字段       | 类型     | 描述         |
|------------|----------|--------------|
|id|int|the association's unique id|
| bk_obj_asst_id| string|  自动生成的模型关联关系id.|
| bk_obj_id| string| 关联关系源模型id |
| bk_asst_obj_id| string| 关联关系目标模型id|
| bk_inst_id| int| 源模型实例id|
| bk_asst_inst_id| int| 目标模型实例id|


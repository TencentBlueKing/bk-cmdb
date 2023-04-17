### 功能描述

获取主线模型的业务拓扑

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx"
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
  "data": [
    {
      "bk_obj_id": "biz",
      "bk_obj_name": "business",
      "bk_supplier_account": "0",
      "bk_next_obj": "set",
      "bk_next_name": "set",
      "bk_pre_obj_id": "",
      "bk_pre_obj_name": ""
    },
    {
      "bk_obj_id": "set",
      "bk_obj_name": "set",
      "bk_supplier_account": "0",
      "bk_next_obj": "module",
      "bk_next_name": "module",
      "bk_pre_obj_id": "biz",
      "bk_pre_obj_name": "business"
    },
    {
      "bk_obj_id": "module",
      "bk_obj_name": "module",
      "bk_supplier_account": "0",
      "bk_next_obj": "host",
      "bk_next_name": "host",
      "bk_pre_obj_id": "set",
      "bk_pre_obj_name": "set"
    },
    {
      "bk_obj_id": "host",
      "bk_obj_name": "host",
      "bk_supplier_account": "0",
      "bk_next_obj": "",
      "bk_next_name": "",
      "bk_pre_obj_id": "module",
      "bk_pre_obj_name": "module"
    }
  ]
}
```

### 返回结果参数说明

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_obj_id | string | 模型的唯一ID |
|bk_obj_name | string |模型名称|
|bk_supplier_account | string |开发商帐户名称|
|bk_next_obj | string |当前模型的下一个模型唯一ID|
|bk_next_name | string |当前模型的下一个模型名称|
|bk_pre_obj_id | string |当前模型的前一个模型的唯一ID|
|bk_pre_obj_name | string |当前模型的前一个模型的名称|

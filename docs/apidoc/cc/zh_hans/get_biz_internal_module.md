### 功能描述

根据业务ID获取业务空闲机, 故障机和待回收模块

### 请求参数

{{ common_args_desc }}


#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | 是     | 业务ID     |

### 请求参数示例

```python

{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":0
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
    "bk_set_id": 2,
    "bk_set_name": "空闲机池",
    "module": [
      {
        "bk_module_id": 3,
        "bk_module_name": "空闲机",
        "default": 1,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 4,
        "bk_module_name": "故障机",
        "default": 2,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 5,
        "bk_module_name": "待回收",
        "default": 3,
        "host_apply_enabled": false
      }
    ]
  }
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


#### data说明
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_set_id | int64 | 空闲机, 故障机和待回收模块所属的set的实例ID |
|bk_set_name | string |空闲机, 故障机和待回收模块所属的set的实例名称|

#### module说明
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_module_id | int | 空闲机, 故障机或待回收模块的实例ID |
|bk_module_name | string |空闲机, 故障机或待回收模块的实例名称|
|default | int | 表示模块类型 |
| host_apply_enabled|bool|是否启用主机属性自动应用|

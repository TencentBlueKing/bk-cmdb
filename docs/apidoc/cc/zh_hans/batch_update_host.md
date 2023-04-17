### 功能描述

根据主机id和属性批量更新主机属性（不能用于更新主机属性中的云区域字段）

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型         | 必选   |  描述                           |
|---------------------|--------------|--------|---------------------------------|
| update              | array | 是     | 主机被更新的属性和值，最多500条   |

#### update
| 字段        | 类型    | 必选   | 描述                                                |
|-------------|--------|--------|----------------------------------------------------|
| properties  | object | 是     | 主机被更新的属性和值，不能用于更新主机属性中的云区域字段 |
| bk_host_id  | int    | 是     | 用于更新的主机ID                                     |

#### properties
| 字段         | 类型   | 必选   | 描述                                                      |
|--------------|--------|-------|-----------------------------------------------------------|
| bk_host_name | string | 否    | 主机名，也可以为其它属性，不能用于更新主机属性中的云区域字段    |
| operator     | string | 否    | 主要维护人，也可以为其它属性，不能用于更新主机属性中的云区域字段 |
| bk_comment   | string | 否    | 备注，也可以为其它属性，不能用于更新主机属性中的云区域字段      |
| bk_isp_name  | string | 否    | 所属运营商，也可以为其它属性，不能用于更新主机属性中的云区域字段 |



### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "update":[
      {
        "properties":{
          "bk_host_name":"batch_update",
          "operator": "admin",
          "bk_comment": "test",
          "bk_isp_name": "1"
        },
        "bk_host_id":46
      }
    ]
}
```


### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": null
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

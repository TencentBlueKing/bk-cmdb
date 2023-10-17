### 功能描述

新增主机到资源池

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_supplier_account |  string     | 否     | 开发商账号 |
| host_info      |  dict    | 是     | 主机信息 |
| bk_biz_id      |  int     | 否     | 业务ID，如果指定业务ID则将主机添加到该业务空闲机池，否则默认添加到资源池 |
| directory | int | 否 | 资源池目录ID |

#### host_info

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_host_innerip |  string   | 是     | 主机内网ip |
| import_from     |  string   | 是     | 主机导入来源,以api方式导入为3 |
| bk_cloud_id     |  int      | 是     | 云区域ID |
| bk_host_name | string | 否 | 主机名，也可以为其它属性 |
| operator | string | 否 | 主要维护人，也可以为其它属性 |
| bk_comment | string | 否 | 备注，也可以为其它属性 |

### 请求参数示例
```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 3,
    "host_info": {
        "0":{
            "bk_host_innerip": "127.0.0.200",
            "bk_host_name": "host11",
            "bk_cloud_id": 0,
            "import_from": "3",
            "operator": "admin",
            "bk_comment": "comment"
        },
         "1":{
            "bk_host_innerip": "127.0.0.201",
            "bk_host_name": "host12",
            "bk_cloud_id": 1,
            "import_from": "3",
            "operator": "admin",
            "bk_comment": "comment"
        }
    },
    "directory": 1
}
```
示例中host_info的"0"表示行数，可按顺序递增
### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "permission": null,
    "data": {
        "success": [
            "0", "1"
        ]
    }
}
```
### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| data    | object | 请求返回的数据                           |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |

#### data

| 字段    | 类型  | 描述           |
| ------- | ----- | -------------- |
| success | array | 添加成功的主机 |

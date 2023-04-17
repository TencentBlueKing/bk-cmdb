### 功能描述

新增云主机到业务的空闲机模块 (云主机管理专用接口, 版本: v3.10.19+, 权限: 业务主机编辑)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型           | 必选  | 描述                                  |
|-----------|--------------|-----|-------------------------------------|
| bk_biz_id | int          | 是   | 业务ID                                |
| host_info | array | 是   | 新增的云主机信息，数组长度最多为200，一批主机仅可同时成功或同时失败 |

#### host_info

主机信息，其中云区域ID和内网IP字段为必填字段，其它字段为主机模型中定义的属性字段。在此仅展示部分字段示例，其它字段请按需填写

| 字段              | 类型     | 必选  | 描述                        |
|-----------------|--------|-----|---------------------------|
| bk_cloud_id     | int    | 是   | 云区域ID                     |
| bk_host_innerip | string | 是   | IPv4格式的主机内网IP，多个IP之间用逗号分隔 |
| bk_host_name    | string | 否   | 主机名，也可以为其它属性              |
| operator        | string | 否   | 主要维护人，也可以为其它属性            |
| bk_comment      | string | 否   | 备注，也可以为其它属性               |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 123,
    "host_info": [
        {
            "bk_cloud_id": 0,
            "bk_host_innerip": "127.0.0.1",
            "bk_host_name": "host1",
            "operator": "admin",
            "bk_comment": "comment"
        },
        {
            "bk_cloud_id": 0,
            "bk_host_innerip": "127.0.0.2",
            "bk_host_name": "host2",
            "operator": "admin",
            "bk_comment": "comment"
        }
    ]
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
    "data": {
        "ids": [
            1,
            2
        ]
    }
}
```

### 返回结果参数说明

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |

#### data
| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
| ids | array | 创建成功的主机的ID数组 |
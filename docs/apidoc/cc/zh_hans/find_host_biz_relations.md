### 功能描述

根据主机ID查询业务相关信息

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段         | 类型    | 必选 | 描述                 |
|------------|-------|----|--------------------|
| bk_host_id | array | 是  | 主机ID数组，ID个数不能超过500 |
| bk_biz_id  | int   | 否  | 业务ID               |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "bk_host_id": [
        3,
        4
    ]
}
```

### 返回结果示例

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": [
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 59,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 61,
      "bk_set_id": 12,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 4,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    }
  ]
}
```

### 返回结果参数说明

#### response

| 字段         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

data 字段说明：

| 字段                  | 类型     | 描述    |
|---------------------|--------|-------|
| bk_biz_id           | int    | 业务ID  |
| bk_host_id          | int    | 主机ID  |
| bk_module_id        | int    | 模块ID  |
| bk_set_id           | int    | 集群ID  |
| bk_supplier_account | string | 开发商账户 |


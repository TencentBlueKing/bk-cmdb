### 功能描述

查询服务分类列表，根据业务ID查询，共用服务分类也会返回

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型  | 必选 | 描述   |
|-----------|-----|----|------|
| bk_biz_id | int | 是  | 业务ID |

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1
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
  "data": {
    "count": 20,
    "info": [
      {
        "bk_biz_id": 0,
        "id": 16,
        "name": "Apache",
        "bk_root_id": 14,
        "bk_parent_id": 14,
        "bk_supplier_account": "0",
        "is_built_in": true
      },
      {
        "bk_biz_id": 0,
        "id": 19,
        "name": "Ceph",
        "bk_root_id": 18,
        "bk_parent_id": 18,
        "bk_supplier_account": "0",
        "is_built_in": true
      },
      {
        "bk_biz_id": 1,
        "id": 1,
        "name": "Default",
        "bk_root_id": 1,
        "bk_supplier_account": "0",
        "is_built_in": true
      }
    ]
  }
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

#### data 字段说明

| 字段    | 类型    | 描述   |
|-------|-------|------|
| count | int   | 总数   |
| info  | array | 返回结果 |

#### info 字段说明

| 字段                  | 类型     | 描述      |
|---------------------|--------|---------|
| id                  | int    | 服务分类ID  |
| name                | string | 服务分类名称  |
| bk_root_id          | int    | 根服务分类ID |
| bk_parent_id        | int    | 父服务分类ID |
| is_built_in         | bool   | 是否内置    |
| bk_supplier_account | string | 开发商帐户名称 |

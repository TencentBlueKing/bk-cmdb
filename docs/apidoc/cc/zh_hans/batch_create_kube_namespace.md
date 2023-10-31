### 功能描述

批量创建namespace (版本：v3.12.1+，权限：容器命名空间新建)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型    | 必选  | 描述                     |
|-----------|-------|-----|------------------------|
| bk_biz_id | int   | 是   | 业务id                   |
| data      | array | 是   | namespace数组, 一次限制创建200 |

#### data[x]

| 字段              | 类型     | 必选  | 描述                         |
|-----------------|--------|-----|----------------------------|
| bk_cluster_id   | int    | 是   | cmdb里标识cluster的唯一id        |
| name            | string | 是   | 命名空间名称                     |
| labels          | map    | 否   | 标签                         |
| resource_quotas | array  | 否   | 命名空间CPU与内存的requests与limits |

#### resource_quotas[x]

| 字段             | 类型     | 必选  | 描述                                                                                                                 |
|----------------|--------|-----|--------------------------------------------------------------------------------------------------------------------|
| hard           | object | 否   | 每个命名资源所需的硬限制                                                                                                       |
| scopes         | array  | 否   | 配额作用域,可选值为："Terminating"、"NotTerminating"、"BestEffort"、"NotBestEffort"、"PriorityClass"、"CrossNamespacePodAffinity" |
| scope_selector | object | 否   | 作用域选择器                                                                                                             |

#### scope_selector

| 字段                | 类型  | 必选    | 描述    |
|-------------------|-----|-------|-------|
| match_expressions | 否   | array | 匹配表达式 |

#### match_expressions[x]

| 字段         | 类型     | 必选  | 描述                                                                                                                 |
|------------|--------|-----|--------------------------------------------------------------------------------------------------------------------|
| scope_name | array  | 是   | 配额作用域,可选值为："Terminating"、"NotTerminating"、"BestEffort"、"NotBestEffort"、"PriorityClass"、"CrossNamespacePodAffinity" |
| operator   | string | 是   | 选择器操作符，可选值为："In"、"NotIn"、"Exists"、"DoesNotExist"                                                                   |
| values     | array  | 否   | 字符串数组，如果操作符为"In"或"NotIn",不能为空，如果为"Exists"或"DoesNotExist"，必须为空                                                      |

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "data": [
    {
      "bk_cluster_id": 1,
      "name": "test",
      "labels": {
        "test": "test",
        "test2": "test2"
      },
      "resource_quotas": [
        {
          "hard": {
            "memory": "20000Gi",
            "pods": "100",
            "cpu": "10k"
          },
          "scope_selector": {
            "match_expressions": [
              {
                "values": [
                  "high"
                ],
                "operator": "In",
                "scope_name": "PriorityClass"
              }
            ]
          }
        }
      ]
    }
  ]
}
```

### 返回结果示例

```json

{
  "result": true,
  "code": 0,
  "data": {
    "ids": [
      1
    ]
  },
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

**注意：**

- 返回的data中的namespaceID数组顺序与参数中的数组数据顺序保持一致。

### 返回结果参数说明

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

#### data

| 字段  | 类型    | 描述                   |
|-----|-------|----------------------|
| ids | array | namespace在cc中的唯一标识数组 |

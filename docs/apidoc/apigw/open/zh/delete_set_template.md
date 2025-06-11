### 描述

根据业务ID和集群模板ID列表删除指定业务下的集群模板(权限：集群模板删除权限)

### 输入参数

| 参数名称             | 参数类型  | 必选 | 描述       |
|------------------|-------|----|----------|
| bk_biz_id        | int   | 是  | 业务ID     |
| set_template_ids | array | 是  | 集群模板ID列表 |

### 调用示例

```json
{
    "bk_biz_id": 20,
    "set_template_ids": [59]
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

### 描述

更新对象实例(权限：模型实例编辑权限)

- 该接口只适用于自定义层级模型和通用模型实例上，不适用于业务、集群、模块、主机等模型实例

### 输入参数

| 参数名称         | 参数类型   | 必选 | 描述                         |
|--------------|--------|----|----------------------------|
| bk_obj_id    | string | 是  | 模型ID                       |
| bk_inst_id   | int    | 是  | 实例ID                       |
| bk_inst_name | string | 否  | 实例名，也可以为其它自定义字段            |
| bk_biz_id    | int    | 否  | 业务ID， 当删除的是自定义主线层级模型实例时则必传 |

注意：当操作的是自定义主线层级模型实例时，而又有使用权限中心的，对于cmdb小于3.9的版本，还需要传包含实例所在业务id的metadata参数，否则会导致权限中心鉴权失败，格式为
"metadata": {
"label": {
"bk_biz_id": "64"
}
}

### 调用示例

```json
{
    "bk_obj_id": "1",
    "bk_inst_id": 0,
    "bk_inst_name": "test"
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
| data       | object | 无数据返回                      |

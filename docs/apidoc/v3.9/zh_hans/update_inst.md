### 功能描述

更新对象实例

- 该接口只适用于自定义层级模型和通用模型实例上，不适用于业务、集群、模块、主机等模型实例

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型      | 必选   |  描述                            |
|---------------------|------------|--------|----------------------------------|
| bk_obj_id           | string     | 是     | 模型ID       |
| bk_inst_id          | int        | 是     | 实例ID |
| bk_inst_name        | string     | 否     | 实例名，也可以为其它自定义字段   |
| bk_biz_id                  | int        | 否     | 业务ID， 当删除的是自定义主线层级模型实例时则必传|

 注意：当操作的是自定义主线层级模型实例时，而又有使用权限中心的，对于cmdb小于3.9的版本，还需要传包含实例所在业务id的metadata参数，否则会导致权限中心鉴权失败，格式为
"metadata": {
    "label": {
        "bk_biz_id": "64"
    }
}

### 请求参数示例(通用实例示例)

```json
{
    "bk_supplier_account": "0",
    "bk_obj_id": "1",
    "bk_inst_id": 0,
    "bk_inst_name": "test"
 }
```

### json

```json

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}
```

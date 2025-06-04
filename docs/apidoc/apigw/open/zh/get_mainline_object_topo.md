### 描述

获取主线模型的业务拓扑

### 输入参数

### 调用示例

```json
{
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [
        {
            "bk_obj_id": "biz",
            "bk_obj_name": "Business",
            "bk_supplier_account": "0",
            "bk_next_obj": "set",
            "bk_next_name": "Set",
            "bk_pre_obj_id": "",
            "bk_pre_obj_name": ""
        },
        {
            "bk_obj_id": "set",
            "bk_obj_name": "Set",
            "bk_supplier_account": "0",
            "bk_next_obj": "module",
            "bk_next_name": "Module",
            "bk_pre_obj_id": "biz",
            "bk_pre_obj_name": "Business"
        },
        {
            "bk_obj_id": "module",
            "bk_obj_name": "Module",
            "bk_supplier_account": "0",
            "bk_next_obj": "host",
            "bk_next_name": "Host",
            "bk_pre_obj_id": "set",
            "bk_pre_obj_name": "Set"
        },
        {
            "bk_obj_id": "host",
            "bk_obj_name": "Host",
            "bk_supplier_account": "0",
            "bk_next_obj": "",
            "bk_next_name": "",
            "bk_pre_obj_id": "module",
            "bk_pre_obj_name": "Module"
        }
    ]
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

#### data

| 参数名称                | 参数类型   | 描述              |
|---------------------|--------|-----------------|
| bk_obj_id           | string | 模型的唯一ID         |
| bk_obj_name         | string | 模型名称            |
| bk_supplier_account | string | 开发商帐户名称         |
| bk_next_obj         | string | 当前模型的下一个模型唯一ID  |
| bk_next_name        | string | 当前模型的下一个模型名称    |
| bk_pre_obj_id       | string | 当前模型的前一个模型的唯一ID |
| bk_pre_obj_name     | string | 当前模型的前一个模型的名称   |

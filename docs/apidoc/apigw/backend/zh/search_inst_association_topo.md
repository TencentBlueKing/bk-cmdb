### 描述

查询实例关联拓扑

### 输入参数

| 参数名称       | 参数类型   | 必选 | 描述   |
|------------|--------|----|------|
| bk_obj_id  | string | 是  | 模型id |
| bk_inst_id | int    | 是  | 实例id |

### 调用示例

```json
{
    "bk_obj_id":"test",
    "bk_inst_id":1
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "data": [
        {
            "id": "",
            "bk_obj_id": "biz",
            "bk_obj_icon": "icon-cc-business",
            "bk_inst_id": 0,
            "bk_obj_name": "business",
            "bk_inst_name": "",
            "asso_id": 0,
            "count": 1,
            "children": [
                {
                    "id": "6",
                    "bk_obj_id": "biz",
                    "bk_obj_icon": "icon-cc-business",
                    "bk_inst_id": 6,
                    "bk_obj_name": "business",
                    "bk_inst_name": "",
                    "asso_id": 558
                }
            ]
        }
    ],
    "message": "success",
    "permission": null,
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

| 参数名称         | 参数类型         | 描述               |
|--------------|--------------|------------------|
| bk_inst_id   | int          | 实例ID             |
| bk_inst_name | string       | 实例用于展示的名字        |
| bk_obj_icon  | string       | 模型图标的名字          |
| bk_obj_id    | string       | 模型ID             |
| bk_obj_name  | string       | 模型用于展示的名字        |
| children     | object array | 本模型下所有被关联的实例的集合  |
| count        | int          | children 包含节点的数量 |

#### children

| 参数名称         | 参数类型   | 描述        |
|--------------|--------|-----------|
| bk_inst_id   | int    | 实例ID      |
| bk_inst_name | string | 实例用于展示的名字 |
| bk_obj_icon  | string | 模型图标的名字   |
| bk_obj_id    | string | 模型ID      |
| bk_obj_name  | string | 模型用于展示的名字 |
| asso_id      | string | 关联id      |

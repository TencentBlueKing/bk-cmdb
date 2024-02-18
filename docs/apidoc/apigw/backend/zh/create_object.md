### 描述

创建模型(权限：模型新建权限)

### 输入参数

| 参数名称                 | 参数类型   | 必选 | 描述                           |
|----------------------|--------|----|------------------------------|
| creator              | string | 否  | 本条数据创建者                      |
| bk_classification_id | string | 是  | 对象模型的分类ID，只能用英文字母序列命名        |
| bk_obj_id            | string | 是  | 对象模型的ID，只能用英文字母序列命名          |
| bk_obj_name          | string | 是  | 对象模型的名字，用于展示，可以使用人类可以阅读的任何语言 |                                             |
| bk_obj_icon          | string | 否  | 对象模型的ICON信息，用于前端显示           |
| obj_sort_number      | int    | 否  | 对象模型在所属模型分组下的排序序号            |

### 调用示例

```json
{
    "creator": "admin",
    "bk_classification_id": "test",
    "bk_obj_name": "test",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "test",
    "obj_sort_number": 1
}
```

### 响应示例

```json
{
    "code": 0,
    "permission": null,
    "result": true,
    "message": "success",
    "data": {
        "description": "",
        "bk_ishidden": false,
        "bk_classification_id": "test",
        "creator": "admin",
        "bk_obj_name": "test",
        "bk_ispaused": false,
        "last_time": null,
        "bk_obj_id": "test",
        "create_time": null,
        "bk_supplier_account": "0",
        "position": "",
        "bk_obj_icon": "icon-cc-business",
        "modifier": "",
        "id": 2000002118,
        "ispre": false,
        "obj_sort_number": 1
    }
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

| 参数名称                 | 参数类型               | 描述                   |
|----------------------|--------------------|----------------------|
| id                   | int                | 新增的数据记录的ID           |
| bk_classification_id | int                | 对象模型的分类ID            |
| creator              | string             | 创建者                  |
| modifier             | string             | 最后修改人员               |
| create_time          | string             | 创建时间                 |
| last_time            | string             | 更新时间                 |
| bk_supplier_account  | string             | 开发商账号                |
| bk_obj_id            | string             | 模型类型                 |
| bk_obj_name          | string             | 模型名称                 |
| bk_obj_icon          | string             | 对象模型的ICON信息，用于前端显示   |
| position             | json object string | 用于前端展示的坐标            |
| ispre                | bool               | 是否预定义, true or false |
| obj_sort_number      | int                | 对象模型在所属模型分组下的排序序号    |

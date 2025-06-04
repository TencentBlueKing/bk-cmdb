### 描述

根据可选条件查询模型(权限：模型查看权限)

### 输入参数

| 参数名称                 | 参数类型   | 必选 | 描述                           |
|----------------------|--------|----|------------------------------|
| creator              | string | 否  | 本条数据创建者                      |
| modifier             | string | 否  | 本条数据的最后修改人员                  |
| bk_classification_id | string | 否  | 对象模型的分类ID，只能用英文字母序列命名        |
| bk_obj_id            | string | 否  | 对象模型的ID，只能用英文字母序列命名          |
| bk_obj_name          | string | 否  | 对象模型的名字，用于展示，可以使用人类可以阅读的任何语言 |
| obj_sort_number      | int    | 否  | 对象模型在所属模型分组下的排序序号            |

### 调用示例

```json
{
    "creator": "user",
    "modifier": "user",
    "bk_classification_id": "test",
    "bk_obj_id": "biz",
    "bk_obj_name": "aaa",
    "obj_sort_number": 1
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": [
        {
            "bk_classification_id": "bk_organization",
            "create_time": "2018-03-08T11:30:28.005+08:00",
            "creator": "cc_system",
            "description": "",
            "id": 4,
            "bk_ispaused": false,
            "ispre": true,
            "last_time": null,
            "modifier": "",
            "bk_obj_icon": "icon-XXX",
            "bk_obj_id": "XX",
            "bk_obj_name": "XXX",
            "position": "{\"test_obj\":{\"x\":-253,\"y\":137}}",
            "bk_supplier_account": "0",
            "obj_sort_number": 1
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

| 参数名称                 | 参数类型               | 描述                                                                                     |
|----------------------|--------------------|----------------------------------------------------------------------------------------|
| id                   | int                | 数据记录的ID                                                                                |
| creator              | string             | 本条数据创建者                                                                                |
| modifier             | string             | 本条数据的最后修改人员                                                                            |
| bk_classification_id | string             | 对象模型的分类ID，只能用英文字母序列命名                                                                  |
| bk_obj_id            | string             | 对象模型的ID，只能用英文字母序列命名                                                                    |
| bk_obj_name          | string             | 对象模型的名字，用于展示                                                                           |
| bk_supplier_account  | string             | 开发商账号                                                                                  |
| bk_ispaused          | bool               | 是否停用, true or false                                                                    |
| ispre                | bool               | 是否预定义, true or false                                                                   |
| bk_obj_icon          | string             | 对象模型的ICON信息，用于前端显示，取值可参考[(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json) |
| position             | json object string | 用于前端展示的坐标                                                                              |
| description          | string             | 数据的描述信息                                                                                |
| obj_sort_number      | int                | 对象模型在所属模型分组下的排序序号                                                                      |

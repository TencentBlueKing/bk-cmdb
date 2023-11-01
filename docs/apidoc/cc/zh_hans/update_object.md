### 功能描述

更新模型定义

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型              | 必选   |  描述                                   |
|---------------------|--------------------|--------|-----------------------------------------|
| id                  | int                | 否     | 对象模型的ID，作为更新操作的条件    |
| modifier            | string             | 否     | 本条数据的最后修改人员    |
| bk_classification_id| string             | 是     | 对象模型的分类ID，只能用英文字母序列命名|
| bk_obj_name         | string             | 否     | 对象模型的名字                          |
| bk_obj_icon         | string             | 否     | 对象模型的ICON信息，用于前端显示，取值可参考[(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json)|
| position            | json object string | 否     | 用于前端展示的坐标                      |
| obj_sort_number     | int    | 否     | 对象模型在所属模型分组下的排序序号；更新该值时当设置的值超过分组模型中该值的最大值，则更新的值为最大值加一，如当设置的值为999，而当前分组模型中该值的最大值为6，则更新值设置为7 |


### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id": 1,
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc2_test_inst",
    "bk_obj_icon": "icon-cc-business",
    "position":"{\"ff\":{\"x\":-863,\"y\":1}}",
    "obj_sort_number": 1
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | object | 无数据返回 |

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



### 请求参数示例

```python
{
    "id": 1,
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc2_test_inst",
    "bk_supplier_account": "0",
    "bk_obj_icon": "icon-cc-business",
    "position":"{\"ff\":{\"x\":-863,\"y\":1}}"
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}
```

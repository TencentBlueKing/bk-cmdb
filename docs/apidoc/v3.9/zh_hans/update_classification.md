### 功能描述

更新模型分类

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                   |  类型    | 必选   |  描述                                      |
|------------------------|----------|--------|--------------------------------------------|
| id                     | int      | 否     | 目标数据的记录ID，作为更新操作的条件       |
| bk_classification_name | string   | 否     | 分类名 |
| bk_classification_icon | string   | 否     | 模型分类的图标,取值可参考，取值可参考[(classIcon.json)](resource_define/classIcon.json) |




### 请求参数示例

```python
{
    "id": 1,
    "bk_classification_name": "cc_test_new",
    "bk_classification_icon": "icon-cc-business"
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

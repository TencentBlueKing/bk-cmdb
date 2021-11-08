### 功能描述

创建模型

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选   |  描述                                                    |
|----------------------|------------|--------|----------------------------------------------------------|
| creator              |string      | 否     | 本条数据创建者                                           |
| bk_classification_id | string     | 是     | 对象模型的分类ID，只能用英文字母序列命名                 |
| bk_obj_id            | string     | 是     | 对象模型的ID，只能用英文字母序列命名                     |
| bk_obj_name          | string     | 是     | 对象模型的名字，用于展示，可以使用人类可以阅读的任何语言 |                                             |
| bk_obj_icon          | string     | 否     | 对象模型的ICON信息，用于前端显示，取值可参考[(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json)|


### 请求参数示例

```python
{
    "creator": "admin",
    "bk_classification_id": "test",
    "bk_obj_name": "test",
    "bk_supplier_account": "0",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "test"
}
```


### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": 1038
    }
}
```

### 返回结果参数说明

#### data

| 字段      | 类型      | 描述               |
|-----------|-----------|--------------------|
| id        | int       | 新增的数据记录的ID |

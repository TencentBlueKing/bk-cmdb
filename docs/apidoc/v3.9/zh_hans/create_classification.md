### 功能描述

添加模型分类

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_classification_id       | string     | 是     | 分类ID，英文描述用于系统内部使用           |
| bk_classification_name     | string     | 是     | 分类名     |
| bk_classification_icon     | string     | 否     | 模型分类的图标,取值可参考，取值可参考[(classIcon.json)](resource_define/classIcon.json)|



### 请求参数示例

```python
{
    "bk_classification_id": "cs_test",
    "bk_classification_name": "test_name",
    "bk_classification_icon": "icon-cc-business"
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": 18
    }
}
```

### 返回结果参数说明

#### data

| 字段       | 类型      | 描述                |
|----------- |-----------|--------------------|
| id         | int       | 新增数据记录的ID   |

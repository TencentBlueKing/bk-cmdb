### 功能描述

查询模型分类

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|

### 请求参数示例

``` python
{
    "bk_supplier_account": "0"
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
     "data": [
         {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "主机管理",
            "bk_classification_type": "inner",
            "id": 1
         }
     ]
}
```

### 返回结果参数说明

#### data

| 字段                   | 类型     | 描述                                                                                          |
|------------------------|----------|-----------------------------------------------------------------------------------------------|
| bk_classification_id   | string   | 分类ID，英文描述用于系统内部使用                                                              |
| bk_classification_name | string   | 分类名                                                                                        |
| bk_classification_type | string   | 用于对分类进行分类（如：inner代码为内置分类，空字符串为自定义分类）                           |
| bk_classification_icon | string   | 模型分类的图标,取值可参考，取值可参考[(classIcon.json)](resource_define/classIcon.json) |
| id                     | int      | 数据记录ID                                                                                    |

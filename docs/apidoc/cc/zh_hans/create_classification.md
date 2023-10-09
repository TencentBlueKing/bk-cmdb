### 功能描述

添加模型分类

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_classification_id       | string     | 是     | 分类ID，英文描述用于系统内部使用           |
| bk_classification_name     | string     | 是     | 分类名     |
| bk_classification_icon     | string     | 否     | 模型分类的图标|



### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "data": {
        "id": 11,
        "bk_classification_id": "cs_test",
        "bk_classification_name": "test_name",
        "bk_classification_type": "",
        "bk_classification_icon": "icon-cc-business",
        "bk_supplier_account": ""
    },
    "message": "success",
    "permission": null,
    "request_id": "76e9134a953b4055bb55853bb248dcb7"
    }
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data

| 字段       | 类型      | 描述                |
|----------- |-----------|--------------------|
| id         | int       | 新增数据记录的ID   |
| bk_classification_id       | string          | 分类ID，英文描述用于系统内部使用           |
| bk_classification_name     | string        | 分类名     |
| bk_classification_icon     | string         | 模型分类的图标|
| bk_classification_type | string   | 用于对分类进行分类（如：inner代码为内置分类，空字符串为自定义分类）                           |
| bk_supplier_account| string| 开发商账号|
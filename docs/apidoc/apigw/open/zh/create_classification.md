### 描述

添加模型分类(权限：模型分组新建权限)

### 输入参数

| 参数名称                   | 参数类型   | 必选 | 描述                |
|------------------------|--------|----|-------------------|
| bk_classification_id   | string | 是  | 分类ID，英文描述用于系统内部使用 |
| bk_classification_name | string | 是  | 分类名               |
| bk_classification_icon | string | 否  | 模型分类的图标           |

### 调用示例

```json
{
    "bk_classification_id": "cs_test",
    "bk_classification_name": "test_name",
    "bk_classification_icon": "icon-cc-business"
}
```

### 响应示例

```json
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

| 参数名称                   | 参数类型   | 描述                                   |
|------------------------|--------|--------------------------------------|
| id                     | int    | 新增数据记录的ID                            |
| bk_classification_id   | string | 分类ID，英文描述用于系统内部使用                    |
| bk_classification_name | string | 分类名                                  |
| bk_classification_icon | string | 模型分类的图标                              |
| bk_classification_type | string | 用于对分类进行分类（如：inner代码为内置分类，空字符串为自定义分类） |
| bk_supplier_account    | string | 开发商账号                                |

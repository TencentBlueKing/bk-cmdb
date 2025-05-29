### 描述

查询模型分类

### 输入参数

### 调用示例

```json
{
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
     "data": [
         {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "主机管理",
            "bk_classification_type": "inner",
            "bk_supplier_account": "0",
            "id": 1
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

| 参数名称                   | 参数类型   | 描述                                                                    |
|------------------------|--------|-----------------------------------------------------------------------|
| bk_classification_id   | string | 分类ID，英文描述用于系统内部使用                                                     |
| bk_classification_name | string | 分类名                                                                   |
| bk_classification_type | string | 用于对分类进行分类（如：inner代码为内置分类，空字符串为自定义分类）                                  |
| bk_classification_icon | string | 模型分类的图标,取值可参考，取值可参考[(classIcon.json)](resource_define/classIcon.json) |
| id                     | int    | 数据记录ID                                                                |
| bk_supplier_account    | string | 开发商账户                                                                 |

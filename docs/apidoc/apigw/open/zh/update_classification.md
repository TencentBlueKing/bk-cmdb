### 描述

更新模型分类(权限：模型分组编辑权限)

### 输入参数

| 参数名称                   | 参数类型   | 必选 | 描述                                                                    |
|------------------------|--------|----|-----------------------------------------------------------------------|
| id                     | int    | 否  | 目标数据的记录ID，作为更新操作的条件                                                   |
| bk_classification_name | string | 否  | 分类名                                                                   |
| bk_classification_icon | string | 否  | 模型分类的图标,取值可参考，取值可参考[(classIcon.json)](resource_define/classIcon.json) |

### 调用示例

```json
{
    "id": 1,
    "bk_classification_name": "cc_test_new",
    "bk_classification_icon": "icon-cc-business"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": "success"
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

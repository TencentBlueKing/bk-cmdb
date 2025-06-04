### 描述

通过对象模型的分组ID查询普通模型拓扑(权限：模型拓扑视图编辑权限)

### 输入参数

| 参数名称                 | 参数类型   | 必选 | 描述        |
|----------------------|--------|----|-----------|
| bk_classification_id | string | 是  | 对象模型的分组ID |

### 调用示例

```json
{
    "bk_classification_id": "test"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": [
        {
           "arrows": "to",
           "from": {
               "bk_classification_id": "bk_host_manage",
               "bk_obj_id": "host",
               "bk_obj_name": "主机",
               "position": "{\"bk_host_manage\":{\"x\":-357,\"y\":-344},\"lhmtest\":{\"x\":163,\"y\":75}}",
               "bk_supplier_account": "0"
           },
           "label": "switch_to_host",
           "label_name": "",
           "label_type": "",
           "to": {
               "bk_classification_id": "bk_network",
               "bk_obj_id": "bk_switch",
               "bk_obj_name": "交换机",
               "position": "{\"bk_network\":{\"x\":-172,\"y\":-160}}",
               "bk_supplier_account": "0"
           }
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

| 参数名称       | 参数类型   | 描述                      |
|------------|--------|-------------------------|
| arrows     | string | 取值 to（单向） 或 to,from（双向） |
| label_name | string | 关联关系的名字                 |
| label      | string | 表明From通过哪个字段关联到To的      |
| from       | string | 对象模型的英文id，拓扑关系的发起方      |
| to         | string | 对象模型的英文ID，拓扑关系的终止方      |

#### from、to

| 参数名称                 | 参数类型               | 描述        |
|----------------------|--------------------|-----------|
| bk_classification_id | string             | 分类ID      |
| bk_obj_id            | string             | 模型id      |
| bk_obj_name          | string             | 模型名称      |
| bk_supplier_account  | string             | 开发商账号     |
| position             | json object string | 用于前端展示的坐标 |

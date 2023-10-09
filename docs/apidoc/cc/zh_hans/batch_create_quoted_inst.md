### 功能描述

批量创建被引用的模型的实例(v3.10.30+，权限：源模型实例的编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段             | 类型           | 必选  | 描述                  |
|----------------|--------------|-----|---------------------|
| bk_obj_id      | string       | 是   | 源模型ID               |
| bk_property_id | string       | 是   | 源模型引用该模型的属性ID       |
| data           | object array | 是   | 需要创建的实例信息，最多不能超过50个 |

#### data[n]

| 参数          | 类型     | 必选              | 描述                                    |
|-------------|--------|-----------------|---------------------------------------|
| bk_inst_id  | int64  | 否               | 源模型实例ID，未填写的情况下需要通过创建源模型实例的接口和源模型实例关联 |
| name        | string | 取决于属性中的"是否必选"配置 | 名称，此处仅为示例，实际字段由模型属性决定                 |
| operator    | string | 取决于属性中的"是否必选"配置 | 维护人，此处仅为示例，实际字段由模型属性决定                | 
| description | string | 取决于属性中的"是否必选"配置 | 描述，此处仅为示例，实际字段由模型属性决定                 |

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "data": [
    {
      "bk_inst_id": 123,
      "name": "test",
      "operator": "user",
      "description": "test instance"
    }
  ]
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "ids": [
      1,
      2
    ]
  },
  "request_id": "dsda1122adasadadada2222"
}
```

**注意：**

- 返回的data中的ID数组顺序与参数中的数组数据顺序保持一致。

### 返回结果参数说明

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

#### data

| 字段  | 类型          | 描述               |
|-----|-------------|------------------|
| ids | int64 array | 创建的实例在cc中的唯一标识数组 |

### 描述

根据服务实例id和设置的标签为服务实例添加标签.(权限：服务实例编辑权限)

### 输入参数

| 参数名称         | 参数类型   | 必选 | 描述                    |
|--------------|--------|----|-----------------------|
| instance_ids | array  | 是  | 服务实例ID,一次最多支持输入100个ID |
| labels       | object | 是  | 添加的Label              |
| bk_biz_id    | int    | 是  | 业务ID                  |

#### labels 字段说明

- key 校验规则: `^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- value 校验规则: `^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### 调用示例

```json
{
  "bk_biz_id": 1,
  "instance_ids": [59, 62],
  "labels": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null

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
